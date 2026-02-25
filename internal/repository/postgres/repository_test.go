package postgres

import (
	"context"
	"testing"

	"github.com/pressly/goose"
	"github.com/stretchr/testify/suite"
	"github.com/tmozzze/org_struct_api/internal/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// RepoTestSuite - struct for testing
type RepoTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *Repo
}

// TearDownTest - cleanup after each test
func (s *RepoTestSuite) TearDownTest() {
	err := s.db.Exec("TRUNCATE TABLE employees, departments RESTART IDENTITY CASCADE").Error
	s.NoError(err, "failed to cleanup database after test")
}

// SetupSuite - initializing for testing
func (s *RepoTestSuite) SetupSuite() {

	// DSN
	dsn := "host=localhost user=user password=password dbname=pgdb port=5432 sslmode=disable"
	// Migrations directory
	migrationsDir := "../../../database/migrations"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		s.T().Fatalf("failed to connect to test database: %v", err)
	}
	s.db = db

	sqlDB, _ := s.db.DB()

	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		s.T().Fatalf("failed to run migrations: %v", err)
	}

	s.repo = NewRepository(db)
}

// TestUniqueNameConstraint - test Department unique name constraint
func (s *RepoTestSuite) TestUniqueNameConstraint() {
	ctx := context.Background()

	// In Root
	deptA := &models.Department{Name: "UniqueDept"}
	s.NoError(s.repo.Department().Create(ctx, deptA))

	deptB := &models.Department{Name: "UniqueDept"}
	err := s.repo.Department().Create(ctx, deptB)
	s.Error(err, "should not allow duplicate department names in Root")

	// In same parent
	childA := &models.Department{Name: "ChildDept", ParentID: &deptA.ID}
	s.NoError(s.repo.Department().Create(ctx, childA))

	childB := &models.Department{Name: "ChildDept", ParentID: &deptA.ID}
	err = s.repo.Department().Create(ctx, childB)
	s.Error(err, "should not allow duplicate department names in same parent")

	// In different parents
	deptC := &models.Department{Name: "Dept C"}
	s.NoError(s.repo.Department().Create(ctx, deptC))

	childC := &models.Department{Name: "ChildDept", ParentID: &deptC.ID}
	s.NoError(s.repo.Department().Create(ctx, childC), "should allow duplicate department names in different parents")

}

// TestGetByID_RecursiveTree - test for DepartmentRepo GetById with recursive
func (s *RepoTestSuite) TestGetByID_RecursiveTree() {
	ctx := context.Background()

	// Root --> Child ->> Grandchild
	root := &models.Department{Name: "Root"}
	s.NoError(s.repo.Department().Create(ctx, root))

	child := &models.Department{Name: "Child", ParentID: &root.ID}
	s.NoError(s.repo.Department().Create(ctx, child))

	grandChild := &models.Department{Name: "Grandchild", ParentID: &child.ID}
	s.NoError(s.repo.Department().Create(ctx, grandChild))

	// Create employee
	emp := &models.Employee{
		FullName:     "Oleg Moroz",
		Position:     "Developer",
		DepartmentID: grandChild.ID,
	}
	s.NoError(s.repo.Employee().Create(ctx, emp))

	// Root with depth = 2
	res, err := s.repo.Department().GetByID(ctx, root.ID, 2, true)

	s.NoError(err)
	s.Equal("Root", res.Name)
	s.Len(res.Children, 1, "Root must have 1 child")
	s.Len(res.Children[0].Children, 1, "Child must have 1 child")

	grandChildRes := res.Children[0].Children[0]
	s.Len(grandChildRes.Employees, 1, "Grandchild must have 1 employee")
	s.Equal("Oleg Moroz", grandChildRes.Employees[0].FullName)
}

// TestDeleteWithReassign - test for DepartmentRepo DeleteWithReassign
func (s *RepoTestSuite) TestDeleteWithReassign() {
	ctx := context.Background()

	deptA := &models.Department{Name: "Dept A"}
	s.NoError(s.repo.Department().Create(ctx, deptA))

	deptB := &models.Department{Name: "Dept B"}
	s.NoError(s.repo.Department().Create(ctx, deptB))

	emp := &models.Employee{
		FullName:     "Oleg Moroz",
		Position:     "Developer",
		DepartmentID: deptA.ID,
	}
	s.NoError(s.repo.Employee().Create(ctx, emp))

	err := s.repo.Department().DeleteWithReassign(ctx, deptA.ID, deptB.ID)
	s.NoError(err)

	exists, _ := s.repo.Department().Exists(ctx, deptA.ID)
	s.False(exists, "Dept A should be deleted")

	var empRes models.Employee
	err = s.db.WithContext(ctx).First(&empRes, emp.ID).Error
	s.NoError(err)
	s.Equal(deptB.ID, empRes.DepartmentID, "Employee should be reassigned to Dept B")
}

func TestRepoSuite(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}

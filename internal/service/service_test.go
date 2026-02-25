package service

import (
	"context"
	"testing"

	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/tmozzze/org_struct_api/internal/domain"
	"github.com/tmozzze/org_struct_api/internal/domain/dto"
	"github.com/tmozzze/org_struct_api/internal/domain/models"
)

// MOCKS

type MockDepartmentRepo struct {
	mock.Mock
}

func (m *MockDepartmentRepo) Create(ctx context.Context, dept *models.Department) error {
	args := m.Called(ctx, dept)
	return args.Error(0)
}

func (m *MockDepartmentRepo) GetByID(ctx context.Context, id int, depth int, includeEmployees bool) (*models.Department, error) {
	args := m.Called(ctx, id, depth, includeEmployees)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *MockDepartmentRepo) GetByIDSimple(ctx context.Context, id int) (*models.Department, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *MockDepartmentRepo) GetByNameAndParent(ctx context.Context, name string, parentID *int) (*models.Department, error) {
	args := m.Called(ctx, name, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *MockDepartmentRepo) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockDepartmentRepo) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDepartmentRepo) DeleteWithReassign(ctx context.Context, id int, reassignToID int) error {
	args := m.Called(ctx, id, reassignToID)
	return args.Error(0)
}

func (m *MockDepartmentRepo) Exists(ctx context.Context, id int) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

type MockRepoWrapper struct {
	mock.Mock
	deptRepo *MockDepartmentRepo
}

func (m *MockRepoWrapper) Department() domain.DepartmentRepository {
	return m.deptRepo
}
func (m *MockRepoWrapper) Employee() domain.EmployeeRepository {
	return nil
}

// SUITE

type DepartmentServiceTestSuite struct {
	suite.Suite
	repo     *MockDepartmentRepo
	wrapper  *MockRepoWrapper
	service  domain.DepartmentService
	validate *validator.Validate
}

func (suite *DepartmentServiceTestSuite) SetupTest() {
	suite.repo = new(MockDepartmentRepo)
	suite.wrapper = &MockRepoWrapper{deptRepo: suite.repo}
	suite.validate = validator.New()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	suite.service = newDepartmentService(suite.wrapper, logger, suite.validate)
}

func TestDepartmentServiceSuite(t *testing.T) {
	suite.Run(t, new(DepartmentServiceTestSuite))
}

// TESTS

func (suite *DepartmentServiceTestSuite) TestCreate_Success() {
	req := &dto.CreateDepartmentRequest{
		Name:     "Backend",
		ParentID: nil,
	}

	suite.repo.On("GetByNameAndParent", mock.Anything, "Backend", mock.Anything).Return(nil, nil)
	suite.repo.On("Create", mock.Anything, mock.AnythingOfType("*models.Department")).Return(nil)

	resp, err := suite.service.Create(context.Background(), req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Backend", resp.Name)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *DepartmentServiceTestSuite) TestCreate_DuplicateName() {
	req := &dto.CreateDepartmentRequest{Name: "HR"}

	suite.repo.On("GetByNameAndParent", mock.Anything, "HR", mock.Anything).
		Return(&models.Department{ID: 1, Name: "HR"}, nil)

	resp, err := suite.service.Create(context.Background(), req)

	assert.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, domain.ErrDuplicateName)
	assert.Nil(suite.T(), resp)
}

func (suite *DepartmentServiceTestSuite) TestUpdate_CycleError() {
	// Try relocation 1 in his grandchild 3
	// Tree: 1 -> 2 -> 3
	idToMove := 1
	newParentID := 3

	req := &dto.UpdateDepartmentRequest{
		ParentID: &newParentID,
	}

	suite.repo.On("GetByIDSimple", mock.Anything, idToMove).Return(&models.Department{ID: 1, ParentID: nil}, nil)
	suite.repo.On("Exists", mock.Anything, newParentID).Return(true, nil)

	// checkCycle: 3 -> 2 -> 1
	suite.repo.On("GetByIDSimple", mock.Anything, 3).Return(&models.Department{ID: 3, ParentID: ptr(2)}, nil)
	suite.repo.On("GetByIDSimple", mock.Anything, 2).Return(&models.Department{ID: 2, ParentID: ptr(1)}, nil)

	resp, err := suite.service.Update(context.Background(), idToMove, req)

	assert.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, domain.ErrCycleConstraint)
	assert.Nil(suite.T(), resp)
}

func (suite *DepartmentServiceTestSuite) TestDelete_ReassignSuccess() {
	idToDelete := 10
	reassignID := 20
	mode := domain.ModeReassign

	req := &dto.DeleteDepartmentRequest{
		Mode:         mode,
		ReassignToID: &reassignID,
	}

	suite.repo.On("Exists", mock.Anything, idToDelete).Return(true, nil)
	suite.repo.On("Exists", mock.Anything, reassignID).Return(true, nil)
	suite.repo.On("DeleteWithReassign", mock.Anything, idToDelete, reassignID).Return(nil)

	err := suite.service.Delete(context.Background(), idToDelete, req)

	assert.NoError(suite.T(), err)
	suite.repo.AssertExpectations(suite.T())
}

func (suite *DepartmentServiceTestSuite) TestDelete_ReassignToSameDepartment() {
	idToDelete := 20
	reassignID := 20
	mode := domain.ModeReassign

	req := &dto.DeleteDepartmentRequest{
		Mode:         mode,
		ReassignToID: &reassignID,
	}

	suite.repo.On("Exists", mock.Anything, idToDelete).Return(true, nil)
	suite.repo.On("Exists", mock.Anything, reassignID).Return(true, nil)
	suite.repo.On("DeleteWithReassign", mock.Anything, idToDelete, reassignID).Return(nil)

	err := suite.service.Delete(context.Background(), idToDelete, req)

	assert.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, domain.ErrInvalidReassignToID)
}

func ptr(i int) *int {
	return &i
}

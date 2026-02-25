package postgres

import (
	"github.com/tmozzze/org_struct_api/internal/domain"
	"gorm.io/gorm"
)

// Repo - main repository struct
type Repo struct {
	department domain.DepartmentRepository
	employee   domain.EmployeeRepository
}

// NewRepository - constructor for Repo
func NewRepository(db *gorm.DB) *Repo {
	return &Repo{
		department: newDepartmentRepo(db),
		employee:   newEmployeeRepo(db),
	}
}

// Department - return DepartmentRepository
func (r *Repo) Department() domain.DepartmentRepository {
	return r.department
}

// Employee - return EmployeeRepository
func (r *Repo) Employee() domain.EmployeeRepository {
	return r.employee
}

package domain

import (
	"context"

	"github.com/tmozzze/org_struct_api/internal/domain/models"
)

// Repository - interface for data repositories
type Repository interface {
	Department() DepartmentRepository
	Employee() EmployeeRepository
}

// DepartmentRepository - interface for department data operations
type DepartmentRepository interface {
	Create(ctx context.Context, dept *models.Department) error
	GetByID(ctx context.Context, id int, depth int, includeEmployees bool) (*models.Department, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) error
	Delete(ctx context.Context, id int) error
	DeleteWithReassign(ctx context.Context, id int, reassignToID int) error
	Exists(ctx context.Context, id int) (bool, error)
}

// EmployeeRepository - interface for employee data operations
type EmployeeRepository interface {
	Create(ctx context.Context, emp *models.Employee) error
	UpdateDepartmentForEmployees(ctx context.Context, oldDeptID int, newDeptID int) error
}

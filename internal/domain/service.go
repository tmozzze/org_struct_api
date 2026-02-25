package domain

import (
	"context"

	"github.com/tmozzze/org_struct_api/internal/domain/dto"
)

const (
	// ModeCascade - delete department and all its sub-departments and employees
	ModeCascade = "cascade"
	// ModeReassign - delete department and reassign all its sub-departments and employees to another department
	ModeReassign = "reassign"
	// DateFormat - standard date format for the application
	DateFormat = "2006-01-02"
)

// Service -
type Service interface {
	Employee() EmployeeService
	Department() DepartmentService
}

// DepartmentService - interface for department business logic
type DepartmentService interface {
	Create(ctx context.Context, req *dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error)
	GetByID(ctx context.Context, id int, req *dto.GetByIDRequest) (*dto.DepartmentResponse, error)
	Update(ctx context.Context, id int, req *dto.UpdateDepartmentRequest) (*dto.DepartmentResponse, error)
	Delete(ctx context.Context, id int, req *dto.DeleteDepartmentRequest) error
}

// EmployeeService - interface for employee business logic
type EmployeeService interface {
	Create(ctx context.Context, deptID int, req *dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error)
}

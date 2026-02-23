package postgres

import (
	"context"

	"github.com/tmozzze/org_struct_api/internal/domain/models"
	"gorm.io/gorm"
)

// EmployeeRepo - PostgreSQL implementation of EmployeeRepository
type EmployeeRepo struct {
	db *gorm.DB
}

// NewEmployeeRepo - constructor for EmployeeRepo
func NewEmployeeRepo(db *gorm.DB) *EmployeeRepo {
	return &EmployeeRepo{
		db: db,
	}
}

// CreateEmployee - create a new employee in department
func (r *EmployeeRepo) Create(ctx context.Context, emp *models.Employee) error {
	return r.db.WithContext(ctx).Create(emp).Error
}

// UpdateDepartmentForEmployees - update department for all employees in oldDeptID to newDeptID
func (r *EmployeeRepo) UpdateDepartmentForEmployees(ctx context.Context, oldDeptID int, newDeptID int) error {
	err := r.db.WithContext(ctx).Model(&models.Employee{}).
		Where("department_id = ?", oldDeptID).
		Update("department_id", newDeptID).Error
	return err
}

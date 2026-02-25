package postgres

import (
	"context"
	"fmt"

	"github.com/tmozzze/org_struct_api/internal/domain/models"
	"gorm.io/gorm"
)

type employeeRepo struct {
	db *gorm.DB
}

func newEmployeeRepo(db *gorm.DB) *employeeRepo {
	return &employeeRepo{
		db: db,
	}
}

// CreateEmployee - create a new employee in department
func (r *employeeRepo) Create(ctx context.Context, emp *models.Employee) error {
	const op = "postgres.employee.Create"

	result := r.db.WithContext(ctx).Create(emp)
	if result.Error != nil {
		return fmt.Errorf("%s: failed to create employee: %w", op, result.Error)
	}

	return nil
}

// UpdateDepartmentForEmployees - update department for all employees in oldDeptID to newDeptID
func (r *employeeRepo) UpdateDepartmentForEmployees(ctx context.Context, oldDeptID int, newDeptID int) error {
	const op = "postgres.employee.UpdateDepartmentForEmployees"

	result := r.db.WithContext(ctx).Model(&models.Employee{}).
		Where("department_id = ?", oldDeptID).
		Update("department_id", newDeptID)

	if result.Error != nil {
		return fmt.Errorf("%s: failed to update department for employees from department id: %d to department id: %d: %w", op, oldDeptID, newDeptID, result.Error)
	}

	return nil
}

package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/tmozzze/org_struct_api/internal/domain"
	"github.com/tmozzze/org_struct_api/internal/domain/models"
	"gorm.io/gorm"
)

type departmentRepo struct {
	db *gorm.DB
}

func newDepartmentRepo(db *gorm.DB) *departmentRepo {
	return &departmentRepo{
		db: db,
	}
}

// Create - create a new department
func (r *departmentRepo) Create(ctx context.Context, dept *models.Department) error {
	const op = "postgres.department.Create"

	result := r.db.WithContext(ctx).Create(dept)
	if result.Error != nil {
		return fmt.Errorf("%s: failed to create department: %w", op, result.Error)
	}

	return nil
}

// GetByID - get department by ID with optional depth and employees
func (r *departmentRepo) GetByID(ctx context.Context, id int, depth int, includeEmployees bool) (*models.Department, error) {
	const op = "postgres.department.GetByID"

	var dept models.Department
	query := r.db.WithContext(ctx)

	if includeEmployees {
		// Sorting employees by full name
		query = query.Preload("Employees", func(db *gorm.DB) *gorm.DB {
			return db.Order("full_name ASC")
		})
	}

	// Children
	currentPath := ""
	for i := 0; i < depth; i++ {
		if currentPath == "" {
			currentPath = "Children"
		} else {
			currentPath += ".Children"
		}

		query = query.Preload(currentPath)

		// Employees for all children
		if includeEmployees {
			query = query.Preload(currentPath+".Employees", func(db *gorm.DB) *gorm.DB {
				return db.Order("full_name ASC")
			})
		}
	}

	err := query.First(&dept, id).Error
	if err != nil {
		// NOT FOUND
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: failed to get department by id: %d: %w", op, id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: failed to get department by id: %d: %w", op, id, err)
	}

	return &dept, nil
}

// Update - update department
func (r *departmentRepo) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	const op = "postgres.department.Update"

	result := r.db.WithContext(ctx).
		Model(&models.Department{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("%s: failed to update department id: %d: %w", op, id, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%s: failed to update department id: %d: %w", op, id, domain.ErrNotFound)
	}

	return nil
}

// Delete - delete department
func (r *departmentRepo) Delete(ctx context.Context, id int) error {
	const op = "postgres.department.Delete"

	result := r.db.WithContext(ctx).Delete(&models.Department{}, id)
	if result.Error != nil {
		return fmt.Errorf("%s: failed to delete department id: %d: %w", op, id, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%s: failed to delete department id: %d: %w", op, id, domain.ErrNotFound)
	}

	return nil
}

// DeleteWithReassign - delete department and reassign employees to another department
func (r *departmentRepo) DeleteWithReassign(ctx context.Context, id int, reassignToID int) error {
	const op = "postgres.department.DeleteWithReassign"

	// Start transaction
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update employees to new department
		if err := tx.Model(&models.Employee{}).
			Where("department_id = ?", id).
			Update("department_id", reassignToID).Error; err != nil {
			return fmt.Errorf("%s: failed to reassign id: %d: %w", op, id, err)
		}

		// Delete department
		result := tx.Delete(&models.Department{}, id)
		if result.Error != nil {
			return fmt.Errorf("%s: failed to delete with reassign department id: %d: %w", op, id, result.Error)
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("%s: failed to delete with reassign department id: %d: %w", op, id, domain.ErrNotFound)
		}

		return nil
	})

}

// GetByNameAndParent - get department by name and parent
func (r *departmentRepo) GetByNameAndParent(ctx context.Context, name string, parentID *int) (*models.Department, error) {
	const op = "postgres.department.GetByNameAndParent"

	var dept models.Department
	query := r.db.WithContext(ctx).Where("name = ?", name)

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	err := query.First(&dept).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // not err - jsut found nothing
		}
		return nil, fmt.Errorf("%s: failed to get department by name and parent: %w", op, err)
	}

	return &dept, nil
}

// GetByIDSimple - get department by ID without children and employees
func (r *departmentRepo) GetByIDSimple(ctx context.Context, id int) (*models.Department, error) {
	const op = "postgres.department.GetByIDSimple"

	var dept models.Department
	if err := r.db.WithContext(ctx).First(&dept, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: failed to get department: %w", op, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: get department by id simple failed: %d: %w", op, id, err)
	}
	return &dept, nil
}

// Exists - check if department exists by ID
func (r *departmentRepo) Exists(ctx context.Context, id int) (bool, error) {
	const op = "postgres.department.Exists"

	var count int64
	err := r.db.WithContext(ctx).Model(&models.Department{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("%s: failed to check existence id: %d: %w", op, id, err)
	}
	return count > 0, nil
}

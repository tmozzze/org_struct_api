package postgres

import (
	"context"
	"errors"

	"github.com/tmozzze/org_struct_api/internal/domain"
	"github.com/tmozzze/org_struct_api/internal/domain/models"
	"gorm.io/gorm"
)

// DepartmentRepo - PostgreSQL implementation of DepartmentRepository
type DepartmentRepo struct {
	db *gorm.DB
}

// NewDepartmentRepo - constructor for DepartmentRepo
func NewDepartmentRepo(db *gorm.DB) *DepartmentRepo {
	return &DepartmentRepo{
		db: db,
	}
}

// Create - create a new department
func (r *DepartmentRepo) Create(ctx context.Context, dept *models.Department) error {
	return r.db.WithContext(ctx).Create(dept).Error
}

// GetByID - get department by ID with optional depth and employees
func (r *DepartmentRepo) GetByID(ctx context.Context, id int, depth int, includeEmployees bool) (*models.Department, error) {
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
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &dept, nil
}

// Update - update department
func (r *DepartmentRepo) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&models.Department{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// Delete - delete department
func (r *DepartmentRepo) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.Department{}, id).Error
}

// DeleteWithReassign - delete department and reassign employees to another department
func (r *DepartmentRepo) DeleteWithReassign(ctx context.Context, id int, reassignToID int) error {
	// Start transaction
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update employees to new department
		if err := tx.Model(&models.Employee{}).
			Where("department_id = ?", id).
			Update("department_id", reassignToID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrNotFound
			}
			return err
		}

		// Delete department
		if err := tx.Delete(&models.Department{}, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrNotFound
			}
			return err
		}

		return nil
	})

}

// Exists - check if department exists by ID
func (r *DepartmentRepo) Exists(ctx context.Context, id int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Department{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

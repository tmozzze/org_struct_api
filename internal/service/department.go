package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/tmozzze/org_struct_api/internal/domain"
	"github.com/tmozzze/org_struct_api/internal/domain/dto"
	"github.com/tmozzze/org_struct_api/internal/domain/models"
)

type departmentService struct {
	repo     domain.Repository
	log      *slog.Logger
	validate *validator.Validate
}

func newDepartmentService(
	repo domain.Repository,
	log *slog.Logger,
	validate *validator.Validate,
) domain.DepartmentService {
	return &departmentService{repo: repo, log: log, validate: validate}
}

// Create - Create a new department
func (s *departmentService) Create(ctx context.Context, req *dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error) {
	const op = "service.department.Create"

	// Trimming space
	req.Name = strings.TrimSpace(req.Name)

	// Validation DTO
	if err := s.validate.Struct(req); err != nil {
		return nil, fmt.Errorf("%s: validation failed: %w", op, err)
	}
	// Check parent department exists
	if req.ParentID != nil {
		exists, err := s.repo.Department().Exists(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to check parent department existence: %w", op, err)
		}
		if !exists {
			return nil, fmt.Errorf("%s: parent department not found: %w", op, domain.ErrParentNotFound)
		}

	}

	// Check name unique
	existing, err := s.repo.Department().GetByNameAndParent(ctx, req.Name, req.ParentID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to check department name uniqueness: %w", op, err)
	}
	if existing != nil {
		return nil, fmt.Errorf("%s: department with name '%s' already exists: %w", op, req.Name, domain.ErrDuplicateName)
	}

	// Mapping DTO to Model
	dept := &models.Department{
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	// Go to repo
	if err := s.repo.Department().Create(ctx, dept); err != nil {
		return nil, fmt.Errorf("%s: failed to create department: %w", op, err)
	}

	// Mapping model to DTO
	resp := dto.NewDepartmentResponse(*dept)
	return &resp, nil

}

// GetByID - Get department by id with depth and include_employees options
func (s *departmentService) GetByID(ctx context.Context, id int, req *dto.GetByIDRequest) (*dto.DepartmentResponse, error) {
	const op = "service.department.GetByID"

	// Set max depth
	if req.Depth > 5 {
		req.Depth = 5
	}

	// Validation DTO
	if err := s.validate.Struct(req); err != nil {
		return nil, fmt.Errorf("%s: validation failed: %w", op, err)
	}

	// Go to repo
	dept, err := s.repo.Department().GetByID(ctx, id, req.Depth, req.IncludeEmployees)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%s: failed to get department by id: %w", op, domain.ErrDepartmentNotFound)
		}
		return nil, fmt.Errorf("%s: failed to get department by id: %w", op, err)
	}

	// Mapping model to DTO
	resp := dto.NewDepartmentResponse(*dept)
	return &resp, nil
}

// Update - Update department by id
func (s *departmentService) Update(ctx context.Context, id int, req *dto.UpdateDepartmentRequest) (*dto.DepartmentResponse, error) {
	const op = "service.department.Update"

	// Validation DTO
	if err := s.validate.Struct(req); err != nil {
		return nil, fmt.Errorf("%s: validation failed: %w", op, err)
	}

	// Get current department for validation
	current, err := s.repo.Department().GetByIDSimple(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%s: failed to get current department: %w", op, domain.ErrDepartmentNotFound)
		}
		return nil, fmt.Errorf("%s: failed to get current department: %w", op, err)
	}

	updates := make(map[string]interface{})

	// Trimming space
	if req.Name != nil {
		trimmedName := strings.TrimSpace(*req.Name)
		if trimmedName == "" {
			return nil, domain.ErrEmptyConstraint
		}
		updates["name"] = trimmedName
	}

	// Check parent department
	if req.ParentID != nil {
		exists, err := s.repo.Department().Exists(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to check parent department existence: %w", op, err)
		}
		if !exists {
			return nil, fmt.Errorf("%s: parent department with id '%d' not found: %w", op, *req.ParentID, domain.ErrParentNotFound)
		}

		newParentID := *req.ParentID

		if newParentID == id {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrCycleConstraint)
		}

		if err := s.checkCycle(ctx, id, newParentID); err != nil {
			return nil, fmt.Errorf("%s: failed to check cycle constraint: %w", op, err)
		}

		updates["parent_id"] = newParentID
	}

	// If no fields to update
	if len(updates) == 0 {
		resp := dto.NewDepartmentResponse(*current)
		return &resp, nil
	}

	// Check name unique
	nameToCheck := current.Name
	if req.Name != nil {
		nameToCheck = strings.TrimSpace(*req.Name)
	}

	parentIDToCheck := current.ParentID
	if req.ParentID != nil {
		parentIDToCheck = req.ParentID
	}

	existing, err := s.repo.Department().GetByNameAndParent(ctx, nameToCheck, parentIDToCheck)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to check department name unique constraint: %w", op, err)
	}
	if existing != nil && existing.ID != id {
		return nil, fmt.Errorf("%s: department with name '%s' and parent_id '%d' already exists: %w", op, nameToCheck, parentIDToCheck, domain.ErrDuplicateName)
	}

	// Go to repo to update
	if err := s.repo.Department().Update(ctx, id, updates); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%s: failed to update department: %w", op, domain.ErrDepartmentNotFound)
		}
		return nil, fmt.Errorf("%s: failed to update department: %w", op, err)
	}

	// Get updated department
	updatedDept, err := s.repo.Department().GetByID(ctx, id, 1, false)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get updated department: %w", op, err)
	}
	// Mapping model to DTO
	resp := dto.NewDepartmentResponse(*updatedDept)
	return &resp, nil
}

// Delete - Delete department by id with mode cascade or reassign
func (s *departmentService) Delete(ctx context.Context, id int, req *dto.DeleteDepartmentRequest) error {
	const op = "service.department.Delete"

	// Validation DTO
	if err := s.validate.Struct(req); err != nil {
		return fmt.Errorf("%s: validation failed: %w", op, err)
	}

	// Check department exists
	exists, err := s.repo.Department().Exists(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: failed to check department existence: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: department with id '%d' does not exist: %w", op, id, domain.ErrDepartmentNotFound)
	}

	// Reassign Mode
	if req.Mode == domain.ModeReassign {
		// Validate reassign_to_id
		if *req.ReassignToID == id {
			return fmt.Errorf("%s: reassign_to_id cannot be the same as department id '%d': %w", op, id, domain.ErrInvalidReassignToID)
		}

		// Check reassign_to department exists
		exists, err := s.repo.Department().Exists(ctx, *req.ReassignToID)
		if err != nil {
			return fmt.Errorf("%s: failed to check reassign_to_id department existence: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: reassign_to_id department with id '%d' does not exist: %w", op, *req.ReassignToID, domain.ErrDepartmentNotFound)
		}

		// Go to repo
		if err := s.repo.Department().DeleteWithReassign(ctx, id, *req.ReassignToID); err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return fmt.Errorf("%s: failed to delete department with reassign, department not found: %w", op, domain.ErrDepartmentNotFound)
			}
			return fmt.Errorf("%s: failed to delete department with reassign: %w", op, err)
		}
	} else {
		// Cascade Mode - Go to repo just Delete
		if err := s.repo.Department().Delete(ctx, id); err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return fmt.Errorf("%s: failed to delete department, department not found: %w", op, domain.ErrDepartmentNotFound)
			}
			return fmt.Errorf("%s: failed to delete department: %w", op, err)
		}
	}

	return nil
}

func (s *departmentService) checkCycle(ctx context.Context, movingID int, newParentID int) error {
	currParentID := &newParentID

	for currParentID != nil {
		if *currParentID == movingID {
			return domain.ErrCycleConstraint
		}

		parent, err := s.repo.Department().GetByIDSimple(ctx, *currParentID)
		if err != nil {
			return err
		}
		currParentID = parent.ParentID
	}
	return nil
}

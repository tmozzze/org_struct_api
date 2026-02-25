package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/tmozzze/org_struct_api/internal/domain"
	"github.com/tmozzze/org_struct_api/internal/domain/dto"
	"github.com/tmozzze/org_struct_api/internal/domain/models"
)

type employeeService struct {
	repo     domain.Repository
	log      *slog.Logger
	validate *validator.Validate
}

func newEmployeeService(
	repo domain.Repository,
	log *slog.Logger,
	validate *validator.Validate,
) domain.EmployeeService {
	return &employeeService{repo: repo, log: log, validate: validate}
}

// Create - Create a new employee in a department
func (s *employeeService) Create(ctx context.Context, deptID int, req *dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error) {
	const op = "service.employee.Create"

	// Validation
	if err := s.validate.Struct(req); err != nil {
		return nil, fmt.Errorf("%s: validation failed: %w", op, err)
	}

	// Check department exists
	exists, err := s.repo.Department().Exists(ctx, deptID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to check department existence: %w", op, err)
	}
	if !exists {
		return nil, fmt.Errorf("%s: parent department not found: %w", op, domain.ErrDepartmentNotFound)
	}

	// Parsing date
	var hiredAt *time.Time
	if req.HiredAt != nil {
		t, err := time.Parse(domain.DateFormat, *req.HiredAt)
		if err != nil {
			return nil, fmt.Errorf("%s: invalid date format: %w", op, fmt.Errorf("invalid date format for hired_at, expected YYYY-MM-DD: %w", err))
		}
		hiredAt = &t
	}

	// Mapping DTO to model
	emp := &models.Employee{
		DepartmentID: deptID,
		FullName:     req.FullName,
		Position:     req.Position,
		HiredAt:      hiredAt,
	}

	// Go to repo
	if err := s.repo.Employee().Create(ctx, emp); err != nil {
		return nil, fmt.Errorf("%s: failed to create employee: %w", op, err)
	}

	// Mapping model to DTO
	resp := dto.NewEmployeeResponse(*emp)
	return &resp, nil
}

package service

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/tmozzze/org_struct_api/internal/domain"
)

// Service - main service struct
type Service struct {
	department domain.DepartmentService
	employee   domain.EmployeeService
	log        *slog.Logger
	validate   *validator.Validate
}

// NewService - constructor for Service
func NewService(
	repo domain.Repository,
	log *slog.Logger,
	validate *validator.Validate,
) domain.Service {
	return &Service{
		department: newDepartmentService(repo, log, validate),
		employee:   newEmployeeService(repo, log, validate),
		log:        log,
		validate:   validate,
	}
}

// Department - return DepartmentService
func (s *Service) Department() domain.DepartmentService {
	return s.department
}

// Employee - return EmployeeService
func (s *Service) Employee() domain.EmployeeService {
	return s.employee
}

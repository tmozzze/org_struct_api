package dto

import (
	"time"

	"github.com/tmozzze/org_struct_api/internal/domain/models"
)

// CreateEmployeeRequest - request payload for creating an employee
type CreateEmployeeRequest struct {
	FullName string  `json:"full_name" validate:"required,min=1,max=200"`
	Position string  `json:"position" validate:"required,min=1,max=200"`
	HiredAt  *string `json:"hired_at" validate:"omitempty,datetime=2006-01-02"`
}

// EmployeeResponse - response payload for employee data
type EmployeeResponse struct {
	ID           int       `json:"id"`
	DepartmentID int       `json:"department_id"`
	FullName     string    `json:"full_name"`
	Position     string    `json:"position"`
	HiredAt      *string   `json:"hired_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// NewEmployeeResponse - convert Employee model to EmployeeResponse DTO
func NewEmployeeResponse(m models.Employee) EmployeeResponse {
	var hiredAtStr *string
	if m.HiredAt != nil {
		str := m.HiredAt.Format("2006-01-02")
		hiredAtStr = &str
	}
	return EmployeeResponse{
		ID:           m.ID,
		DepartmentID: m.DepartmentID,
		FullName:     m.FullName,
		Position:     m.Position,
		HiredAt:      hiredAtStr,
		CreatedAt:    m.CreatedAt,
	}
}

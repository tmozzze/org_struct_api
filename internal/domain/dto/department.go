package dto

import (
	"time"

	"github.com/tmozzze/org_struct_api/internal/domain/models"
)

// CreateDepartmentRequest - request payload for creating a department
type CreateDepartmentRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=200"`
	ParentID *int   `json:"parent_id" validate:"omitempty,gt=0"`
}

// GetByIDRequest - request payload for getting by id
type GetByIDRequest struct {
	Depth            int  `json:"depth" validate:"min=1,max=5"`
	IncludeEmployees bool `json:"include_employees"`
}

// UpdateDepartmentRequest - request payload for updating a department
type UpdateDepartmentRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=1,max=200"`
	ParentID *int    `json:"parent_id" validate:"omitempty,gt=0"`
}

// DeleteDepartmentRequest - request payload for deleting
type DeleteDepartmentRequest struct {
	Mode         string `json:"mode" validate:"required,oneof=cascade reassign"`
	ReassignToID *int   `json:"reassign_to_id" validate:"required_if=Mode reassign,omitempty,gt=0"`
}

// DepartmentResponse - response payload for department data
type DepartmentResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	ParentID  *int      `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`

	Employees []EmployeeResponse   `json:"employees,omitempty"`
	Children  []DepartmentResponse `json:"children,omitempty"`
}

// NewDepartmentResponse - convert Department model to DepartmentResponse DTO
func NewDepartmentResponse(m models.Department) DepartmentResponse {
	resp := DepartmentResponse{
		ID:        m.ID,
		Name:      m.Name,
		ParentID:  m.ParentID,
		CreatedAt: m.CreatedAt,
	}

	if len(m.Employees) > 0 {
		resp.Employees = make([]EmployeeResponse, len(m.Employees))
		for i, emp := range m.Employees {
			resp.Employees[i] = NewEmployeeResponse(emp)
		}
	}

	if len(m.Children) > 0 {
		resp.Children = make([]DepartmentResponse, len(m.Children))
		for i, child := range m.Children {
			resp.Children[i] = NewDepartmentResponse(child)
		}
	}

	return resp
}

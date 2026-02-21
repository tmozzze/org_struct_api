package models

import "time"

// Department - represent a department in organisation
type Department struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(200);not null;index:idx_parent_name,unique" validate:"required,min=1,max=200"`
	ParentID  *int      `json:"parent_id" gorm:"index:idx_parent_name,unique" validate:"omitempty,gt=0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	Employees []Employee   `json:"employees,omitempty" gorm:"foreignKey:DepartmentID;constraint:OnDelete:CASCADE"`
	Children  []Department `json:"children,omitempty" gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE"`
}

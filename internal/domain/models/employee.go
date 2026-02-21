package models

import "time"

// Employee - represent an employee in organisation
type Employee struct {
	ID           int        `json:"id" gorm:"primaryKey"`
	DepartmentID int        `json:"department_id" gorm:"not null;index" validate:"required,gt=0"`
	FullName     string     `json:"full_name" gorm:"type:varchar(200);not null" validate:"required,min=1,max=200"`
	Position     string     `json:"position" gorm:"type:varchar(200);not null" validate:"required,min=1,max=200"`
	HiredAt      *time.Time `json:"hired_at" gorm:"type:date"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

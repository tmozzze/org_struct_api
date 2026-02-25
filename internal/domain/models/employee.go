package models

import "time"

// Employee - represent an employee in organisation
type Employee struct {
	ID           int        `json:"id" gorm:"primaryKey"`
	DepartmentID int        `json:"department_id" gorm:"not null;index"`
	FullName     string     `json:"full_name" gorm:"type:varchar(200);not null"`
	Position     string     `json:"position" gorm:"type:varchar(200);not null"`
	HiredAt      *time.Time `json:"hired_at" gorm:"type:date"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

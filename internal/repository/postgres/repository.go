package postgres

import "gorm.io/gorm"

// Repo - main repository struct
type Repo struct {
	DepartmentRepo *DepartmentRepo
	EmployeeRepo   *EmployeeRepo
}

// NewRepository - constructor for Repo
func NewRepository(db *gorm.DB) *Repo {
	return &Repo{
		DepartmentRepo: NewDepartmentRepo(db),
		EmployeeRepo:   NewEmployeeRepo(db),
	}
}

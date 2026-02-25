package domain

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrDepartmentNotFound = errors.New("department not found")
	ErrParentNotFound     = errors.New("parent not found")

	ErrDuplicateName = errors.New("duplicate name")
	ErrAlreadyExist  = errors.New("entity already exists")

	ErrCycleConstraint  = errors.New("cycle constraint")
	ErrLengthConstraint = errors.New("length constraint")
	ErrEmptyConstraint  = errors.New("empty constraint")

	ErrInvalidReassignToID = errors.New("invalid reassign_to_id")
)

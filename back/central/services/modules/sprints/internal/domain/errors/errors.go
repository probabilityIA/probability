package errors

import "errors"

var (
	ErrSprintNotFound   = errors.New("sprint not found")
	ErrNameRequired     = errors.New("name is required")
	ErrInvalidStatus    = errors.New("invalid status")
	ErrInvalidStartDate = errors.New("invalid start_date, expected YYYY-MM-DD or RFC3339")
	ErrInvalidEndDate   = errors.New("invalid end_date, expected YYYY-MM-DD or RFC3339")
	ErrInvalidDateRange = errors.New("end_date must be after start_date")
	ErrCreatorRequired  = errors.New("created_by_id is required")
)

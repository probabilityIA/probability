package errors

import "errors"

var (
	ErrAnnouncementNotFound = errors.New("announcement not found")
	ErrInvalidDateRange     = errors.New("starts_at must be before ends_at")
	ErrInvalidDisplayType   = errors.New("invalid display type, must be: modal_image, modal_text, or ticker")
	ErrInvalidFrequencyType = errors.New("invalid frequency type")
	ErrInvalidStatus        = errors.New("invalid status")
	ErrCategoryNotFound     = errors.New("category not found")
	ErrTargetsRequired      = errors.New("target business IDs are required when announcement is not global")
)

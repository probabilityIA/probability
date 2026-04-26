package errors

import "errors"

var (
	ErrTicketNotFound      = errors.New("ticket not found")
	ErrCommentNotFound     = errors.New("comment not found")
	ErrAttachmentNotFound  = errors.New("attachment not found")
	ErrForbidden           = errors.New("forbidden")
	ErrInvalidStatus       = errors.New("invalid status")
	ErrInvalidPriority     = errors.New("invalid priority")
	ErrInvalidType         = errors.New("invalid type")
	ErrInvalidSeverity     = errors.New("invalid severity")
	ErrTitleRequired       = errors.New("title is required")
	ErrDescriptionRequired = errors.New("description is required")
	ErrAssigneeNotFound    = errors.New("assigned user not found")
)

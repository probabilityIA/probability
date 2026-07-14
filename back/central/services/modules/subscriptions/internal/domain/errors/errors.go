package errors

import "errors"

var (
	ErrSubscriptionTypeNotFound = errors.New("subscription type not found")
	ErrSubscriptionTypeInactive = errors.New("subscription type is not active")
	ErrInvalidSubscriptionType  = errors.New("subscription type name, code and price are required")
	ErrInsufficientBalance      = errors.New("insufficient wallet balance")
	ErrInvalidMonths            = errors.New("months must be greater than zero")
	ErrSubscriptionNotFound     = errors.New("subscription not found")
	ErrInvalidModuleCode        = errors.New("invalid module code")
	ErrOverrideNotFound         = errors.New("override not found")
)

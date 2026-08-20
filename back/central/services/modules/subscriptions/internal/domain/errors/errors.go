package errors

import "errors"

var (
	ErrSubscriptionTypeNotFound    = errors.New("subscription type not found")
	ErrSubscriptionTypeInactive    = errors.New("subscription type is not active")
	ErrInvalidSubscriptionType     = errors.New("subscription type name, code and price are required")
	ErrInsufficientBalance         = errors.New("insufficient wallet balance")
	ErrInvalidMonths               = errors.New("months must be greater than zero")
	ErrSubscriptionNotFound        = errors.New("subscription not found")
	ErrInvalidModuleCode           = errors.New("invalid module code")
	ErrOverrideNotFound            = errors.New("override not found")
	ErrBusinessRequired            = errors.New("business_id is required for a custom plan")
	ErrInvalidDateRange            = errors.New("end_date must be after start_date")
	ErrSubscriptionAlreadyReverted = errors.New("subscription already reverted")
	ErrSubscriptionNotPaid         = errors.New("only paid subscriptions can be reverted")
	ErrNothingToReactivate         = errors.New("subscription has no prior state to reactivate")
	ErrInvalidDays                 = errors.New("days must be greater than zero")
	ErrNoOverageDue                = errors.New("no hay cargo de excedente pendiente de pago")
)

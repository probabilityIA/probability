package errors

import "errors"

var (
	// ErrOrderAlreadyExists indicates that an order with the same external ID already exists for the integration
	ErrOrderAlreadyExists = errors.New("order with this external_id already exists for this integration")
	// ErrOrderNotFound indicates that the requested order does not exist
	ErrOrderNotFound = errors.New("order not found")
	// ErrInvalidStatus indicates that the status code is not valid
	ErrInvalidStatus = errors.New("invalid status code")
	// ErrInvalidStatusTransition indicates that the transition between statuses is not allowed
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	// ErrOrderInTerminalState indicates that the order is in a terminal state and cannot be changed
	ErrOrderInTerminalState = errors.New("order is in a terminal state and cannot be changed")
)

package errors

import "errors"

var (
	ErrRouteNotFound       = errors.New("route not found")
	ErrStopNotFound        = errors.New("route stop not found")
	ErrInvalidTransition   = errors.New("invalid status transition")
	ErrRouteNotPlanned     = errors.New("route must be in planned status")
	ErrRouteNotInProgress  = errors.New("route must be in in_progress status")
	ErrDriverNotFound      = errors.New("driver not found")
	ErrVehicleNotFound     = errors.New("vehicle not found")
	ErrOrderNotFound       = errors.New("order not found")
	ErrStopIDsMismatch     = errors.New("stop IDs do not match route stops")
)

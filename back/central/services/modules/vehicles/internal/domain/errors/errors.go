package errors

import "errors"

var (
	ErrVehicleNotFound       = errors.New("vehicle not found")
	ErrDuplicateLicensePlate = errors.New("a vehicle with this license plate already exists in your business")
)

package errors

import "errors"

var (
	ErrDriverNotFound          = errors.New("driver not found")
	ErrDuplicateIdentification = errors.New("a driver with this identification already exists in your business")
)

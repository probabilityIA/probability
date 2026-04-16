package errors

import "errors"

var (
	ErrClientNotFound  = errors.New("client not found")
	ErrDuplicateEmail  = errors.New("a client with this email already exists in your business")
	ErrDuplicateDni    = errors.New("a client with this DNI already exists in your business")
	ErrClientHasOrders = errors.New("client has orders and cannot be deleted")
)

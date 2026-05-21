package errors

import "errors"

var (
	ErrGroupNotFound      = errors.New("client group not found")
	ErrGroupNameRequired  = errors.New("group name is required")
	ErrGroupNameDuplicate = errors.New("a client group with that name already exists")
	ErrClientNotFound     = errors.New("client not found")
	ErrInvalidPrice       = errors.New("price must be greater than or equal to zero")
	ErrTargetRequired     = errors.New("either client_group_id or client_id is required")
	ErrTargetAmbiguous    = errors.New("only one of client_group_id or client_id may be set")
	ErrProductNotFound    = errors.New("product not found")
	ErrNoClients          = errors.New("at least one client id is required")
)

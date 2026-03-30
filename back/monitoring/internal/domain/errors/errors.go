package errors

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrAccessDenied       = errors.New("access denied: only platform admins can access monitoring")
	ErrContainerNotFound  = errors.New("container not found")
	ErrInvalidAction      = errors.New("invalid action, must be: restart, stop, start")
	ErrActionFailed       = errors.New("container action failed")
)

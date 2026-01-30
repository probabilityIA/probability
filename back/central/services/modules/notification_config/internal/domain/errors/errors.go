package errors

import "errors"

var (
	// ErrNotificationConfigNotFound se lanza cuando no se encuentra una configuración
	ErrNotificationConfigNotFound = errors.New("notification config not found")

	// ErrInvalidConditions se lanza cuando las condiciones son inválidas
	ErrInvalidConditions = errors.New("invalid notification conditions")

	// ErrInvalidConfig se lanza cuando la configuración es inválida
	ErrInvalidConfig = errors.New("invalid notification config")

	// ErrDuplicateConfig se lanza cuando ya existe una configuración similar
	ErrDuplicateConfig = errors.New("duplicate notification config")
)

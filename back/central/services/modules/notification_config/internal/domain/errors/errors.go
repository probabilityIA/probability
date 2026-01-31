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

	// ErrNotificationTypeNotFound se lanza cuando no se encuentra un tipo de notificación
	ErrNotificationTypeNotFound = errors.New("notification type not found")

	// ErrNotificationEventTypeNotFound se lanza cuando no se encuentra un tipo de evento de notificación
	ErrNotificationEventTypeNotFound = errors.New("notification event type not found")

	// ErrNotFound se lanza cuando no se encuentra un recurso genérico
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput se lanza cuando los datos de entrada son inválidos
	ErrInvalidInput = errors.New("invalid input data")
)

package errors

import "errors"

var (
	ErrInvalidSignature = errors.New("webhook signature inválida")
	ErrEmptyAlerts      = errors.New("payload sin alertas")
)

package errors

import "errors"

var (
	ErrTokenRequired    = errors.New("el token del dispositivo es requerido")
	ErrPlatformInvalid  = errors.New("la plataforma debe ser android, ios o web")
	ErrBusinessRequired = errors.New("el negocio es requerido")
	ErrFCMNotConfigured = errors.New("el envio de push no esta configurado")
)

func IsNonRetryable(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrTokenRequired) ||
		errors.Is(err, ErrPlatformInvalid) ||
		errors.Is(err, ErrBusinessRequired) ||
		errors.Is(err, ErrFCMNotConfigured)
}

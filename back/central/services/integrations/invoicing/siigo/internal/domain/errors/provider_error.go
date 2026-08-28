package errors

import "errors"

type ProviderError struct {
	Codigo  string
	Mensaje string
}

func (e *ProviderError) Error() string {
	return e.Mensaje
}

func NewProviderError(codigo, mensaje string) *ProviderError {
	return &ProviderError{Codigo: codigo, Mensaje: mensaje}
}

func CodeFromError(err error) string {
	if err == nil {
		return ""
	}
	var pe *ProviderError
	if errors.As(err, &pe) && pe.Codigo != "" {
		return pe.Codigo
	}
	return "api_error"
}

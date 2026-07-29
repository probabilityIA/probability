package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("credenciales de Shipit invalidas")
	ErrCancelNotSupported = errors.New("Shipit no soporta cancelacion de envios via API; debe gestionarse desde el panel de Shipit")
	ErrCommuneNotFound    = errors.New("no se encontro la comuna de destino en Shipit")
	ErrEmptyResponse      = errors.New("respuesta vacia de Shipit")
)

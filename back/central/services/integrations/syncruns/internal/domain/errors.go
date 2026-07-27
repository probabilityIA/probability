package domain

import "errors"

var (
	ErrIntegrationNotOwned = errors.New("la integracion no pertenece al negocio")
	ErrInvalidKind         = errors.New("kind invalido: use inventory o products")
)

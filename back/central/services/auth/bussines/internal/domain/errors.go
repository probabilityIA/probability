package domain

import "errors"

var (
	ErrBusinessCodeAlreadyExists   = errors.New("el código del negocio ya está en uso")
	ErrBusinessDomainAlreadyExists = errors.New("el dominio personalizado ya está en uso")
)

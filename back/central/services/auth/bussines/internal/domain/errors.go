package domain

import "errors"

var (
	ErrBusinessNotFound            = errors.New("negocio no encontrado")
	ErrBusinessCodeAlreadyExists   = errors.New("el código del negocio ya está en uso")
	ErrBusinessDomainAlreadyExists = errors.New("el dominio personalizado ya está en uso")
	ErrBusinessTypeIDRequired      = errors.New("el tipo de negocio es obligatorio")
	ErrBusinessTypeIDInvalid       = errors.New("el tipo de negocio especificado no existe o no es válido")

	ErrBusinessTypeNotFound          = errors.New("tipo de negocio no encontrado")
	ErrBusinessTypeNameAlreadyExists = errors.New("el nombre del tipo de negocio ya está en uso")
)

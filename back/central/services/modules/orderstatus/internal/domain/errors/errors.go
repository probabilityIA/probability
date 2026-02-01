package errors

import "errors"

var (
	// ErrMappingNotFound indica que no se encontró el mapeo de estado
	ErrMappingNotFound = errors.New("order status mapping not found")

	// ErrMappingAlreadyExists indica que ya existe un mapeo para la combinación IntegrationType + OriginalStatus
	ErrMappingAlreadyExists = errors.New("mapping already exists for this integration type and original status")

	// ErrInvalidID indica que el ID proporcionado es inválido
	ErrInvalidID = errors.New("invalid ID")

	// ErrInvalidFilters indica que los filtros proporcionados son inválidos
	ErrInvalidFilters = errors.New("invalid filters")
)

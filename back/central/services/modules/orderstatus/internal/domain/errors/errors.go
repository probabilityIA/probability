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

	// ErrOrderStatusNotFound indica que no se encontró el estado de orden
	ErrOrderStatusNotFound = errors.New("order status not found")

	// ErrOrderStatusHasMappings indica que el estado de orden tiene mapeos activos y no puede eliminarse
	ErrOrderStatusHasMappings = errors.New("order status has mappings and cannot be deleted")

	// ErrChannelStatusNotFound indica que no se encontró el estado de canal
	ErrChannelStatusNotFound = errors.New("channel status not found")

	// ErrChannelStatusHasMappings indica que el estado de canal tiene mapeos activos y no puede eliminarse
	ErrChannelStatusHasMappings = errors.New("channel status has order status mappings and cannot be deleted")
)

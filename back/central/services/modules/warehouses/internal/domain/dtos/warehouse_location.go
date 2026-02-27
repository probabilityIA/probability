package dtos

// CreateLocationDTO datos para crear una ubicación dentro de una bodega
type CreateLocationDTO struct {
	WarehouseID   uint
	BusinessID    uint
	Name          string
	Code          string
	Type          string
	IsActive      bool
	IsFulfillment bool
	Capacity      *int
}

// UpdateLocationDTO datos para actualizar una ubicación
type UpdateLocationDTO struct {
	ID            uint
	WarehouseID   uint
	BusinessID    uint
	Name          string
	Code          string
	Type          string
	IsActive      bool
	IsFulfillment bool
	Capacity      *int
}

// ListLocationsParams parámetros para listar ubicaciones de una bodega
type ListLocationsParams struct {
	WarehouseID uint
	BusinessID  uint
}

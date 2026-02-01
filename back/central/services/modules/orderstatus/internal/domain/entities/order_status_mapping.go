package entities

import "time"

// OrderStatusMapping representa un mapeo de estado de orden en el dominio
// PURO - Sin tags, sin dependencias de frameworks
type OrderStatusMapping struct {
	ID                uint
	IntegrationTypeID uint
	OriginalStatus    string
	OrderStatusID     uint
	IsActive          bool
	Priority          int
	Description       string
	CreatedAt         time.Time
	UpdatedAt         time.Time

	// Relaciones (opcionales, para cuando se carga con join)
	IntegrationType *IntegrationTypeInfo
	OrderStatus     *OrderStatusInfo
}

package domain

import "time"

// OrderStatusMapping representa un mapeo de estado de orden en el dominio
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
	IntegrationType *IntegrationTypeInfo `json:"integration_type,omitempty"`
	OrderStatus     *OrderStatusInfo     `json:"order_status,omitempty"`
}

// IntegrationTypeInfo contiene información básica del tipo de integración
type IntegrationTypeInfo struct {
	ID       uint   `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

// OrderStatusInfo contiene información básica del estado de orden
type OrderStatusInfo struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Color       string `json:"color"`
}

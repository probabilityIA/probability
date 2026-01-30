package request

import "context"

// OrderRepository define la interfaz para obtener datos de órdenes
type OrderRepository interface {
	GetByID(ctx context.Context, id string) (*OrderData, error)
}

// OrderData representa la información de la orden necesaria para validar condiciones
type OrderData struct {
	ID              string
	OrderNumber     string
	Status          string
	PaymentMethodID uint
	CustomerPhone   string
	TotalAmount     float64
	Currency        string
	BusinessID      *uint
}

// IntegrationRepository define la interfaz para obtener integraciones
type IntegrationRepository interface {
	GetWhatsAppByBusinessID(ctx context.Context, businessID uint) (*IntegrationData, error)
}

// IntegrationData representa la información de la integración
type IntegrationData struct {
	ID         uint
	BusinessID uint
	IsActive   bool
}

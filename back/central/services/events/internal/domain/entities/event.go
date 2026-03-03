package entities

import (
	"time"
)

// Event es la estructura base para todos los eventos del sistema.
// Representa un evento genérico que puede ser de órdenes, facturación,
// envíos o integraciones.
type Event struct {
	ID            string
	Type          string
	Category      string // "order", "invoice", "shipment", "integration"
	BusinessID    uint
	IntegrationID uint
	Timestamp     time.Time
	Data          map[string]interface{}
	Metadata      map[string]interface{}
}

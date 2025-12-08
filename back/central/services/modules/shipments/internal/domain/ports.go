package domain

import (
	"context"
)

// ───────────────────────────────────────────
//
//	REPOSITORY INTERFACE
//
// ───────────────────────────────────────────

// IRepository define todos los métodos de repositorio del módulo shipments
type IRepository interface {
	// CRUD Operations
	CreateShipment(ctx context.Context, shipment *Shipment) error
	GetShipmentByID(ctx context.Context, id uint) (*Shipment, error)
	GetShipmentByTrackingNumber(ctx context.Context, trackingNumber string) (*Shipment, error)
	GetShipmentsByOrderID(ctx context.Context, orderID string) ([]Shipment, error)
	ListShipments(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]Shipment, int64, error)
	UpdateShipment(ctx context.Context, shipment *Shipment) error
	DeleteShipment(ctx context.Context, id uint) error

	// Validation
	ShipmentExists(ctx context.Context, orderID string, trackingNumber string) (bool, error)
}


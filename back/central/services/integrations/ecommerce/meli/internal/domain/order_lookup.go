package domain

import "context"

type MeliOrderRef struct {
	IntegrationID uint
	ShipmentID    int64
}

type IOrderLookupRepository interface {
	GetMeliShipmentByOrderID(ctx context.Context, orderID string) (*MeliOrderRef, error)
}

package domain

import (
	"context"
	"time"
)

type ManifestFilter struct {
	BusinessID      uint
	IncludeChildren bool
	Carrier         string
}

type ManifestShipmentRow struct {
	ShipmentID         uint
	OrderID            *string
	OrderNumber        string
	TrackingNumber     string
	Carrier            string
	CarrierCode        string
	CustomerName       string
	CustomerDocument   string
	DestinationAddress string
	DestinationCity    string
	DestinationState   string
	Weight             float64
	DeclaredValue      float64
	CodTotal           float64
	BusinessID         uint
	BusinessName       string
	WarehouseName      string
	CreatedAt          time.Time
	ShipmentCreatedAt  *time.Time
	OrderCreatedAt     *time.Time
	ShipmentStatus     string
	OrderStatus        string
}

type ManifestPDFInput struct {
	BusinessID   uint
	BusinessName string
	BusinessNIT  string
	BusinessCode string
	OriginCity   string
	GeneratedAt  time.Time
	GeneratedBy  string
	ManifestNo   string
	Carrier      string
	Rows         []ManifestShipmentRow
}

type IManifestRepository interface {
	ListPendingForManifest(ctx context.Context, filter ManifestFilter) ([]ManifestShipmentRow, error)
	GetBusinessForManifest(ctx context.Context, businessID uint) (*ManifestBusinessInfo, error)
	GetChildBusinessIDs(ctx context.Context, parentID uint) ([]uint, error)
}

type ManifestBusinessInfo struct {
	ID       uint
	Name     string
	NIT      string
	Code     string
	Address  string
	City     string
	ParentID *uint
}

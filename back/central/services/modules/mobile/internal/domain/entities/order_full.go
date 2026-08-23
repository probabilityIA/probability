package entities

import "time"

type OrderSummary struct {
	ID              string
	OrderNumber     string
	InternalNumber  string
	IntegrationID   uint
	IntegrationType string
	Platform        string
	Status          string
	StatusID        *uint
	IsPaid          bool
	IsCod           bool
	CodTotal        *float64
	Subtotal        float64
	Tax             float64
	Discount        float64
	ShippingCost    float64
	TotalAmount     float64
	Currency        string
	CustomerName    string
	CustomerEmail   string
	CustomerPhone   string
	ShippingStreet  string
	ShippingCity    string
	WarehouseName   string
	UserName        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderLine struct {
	ID         uint
	ProductID  *string
	SKU        string
	Name       string
	Quantity   int
	UnitPrice  float64
	TotalPrice float64
}

type ShipmentSummary struct {
	ID                   uint
	TrackingNumber       *string
	TrackingURL          *string
	Carrier              *string
	GuideURL             *string
	Status               string
	CarrierStatus        *string
	DestinationCity      string
	InsuranceCost        *float64
	TotalCost            *float64
	CarrierCost          *float64
	AppliedMargin        *float64
	CodCarrierFee        *float64
	CodProbabilityMargin *float64
	CreatedAt            time.Time
}

type InvoiceSummary struct {
	ID            uint
	InvoiceNumber string
	Status        string
	TotalAmount   float64
	Currency      string
	InvoiceURL    *string
	CUFE          *string
	IssuedAt      *time.Time
	CreatedAt     time.Time
}

type OrderFull struct {
	Order    OrderSummary
	Items    []OrderLine
	Shipment *ShipmentSummary
	Invoice  *InvoiceSummary
}

package entities

import "time"

type OrderItemInfo struct {
	Name     string
	Variant  string
	Quantity int
}

type TrackingEvent struct {
	Date        string
	Status      string
	RawStatus   string
	Description string
}

type ShipmentInfo struct {
	TrackingNumber      string
	Carrier             string
	Status              string
	CarrierStatus       string
	CarrierStatusDetail string
	CodCollectAmount    *float64
	TotalCost           *float64
	DestinationCity     string
	OrderNumber         string
	IsTest              bool
	CreatedAt           time.Time
	ShippedAt           *time.Time
	DeliveredAt         *time.Time
	EstimatedDelivery   *time.Time
	Events              []TrackingEvent
}

type OrderInfo struct {
	Number        string
	Platform      string
	CreatedAt     time.Time
	StatusName    string
	PaymentStatus string
	IsPaid        bool
	Total         float64
	Currency      string
	IsCod         bool
	CodTotal      *float64
	CustomerName  string
	City          string
	IsConfirmed   bool
	Novelty       string
	HasInvoice    bool
	IsTest        bool
	Items         []OrderItemInfo
	Shipment      *ShipmentInfo
}

type OrderSummary struct {
	Number         string
	CreatedAt      time.Time
	CustomerName   string
	Total          float64
	Currency       string
	StatusName     string
	IsCod          bool
	TrackingNumber string
	ShipmentStatus string
	IsTest         bool
}

type StatusCount struct {
	Name  string
	Count int64
}

type OrdersOverview struct {
	From           time.Time
	To             time.Time
	Orders         int64
	Amount         float64
	CashOnDelivery int64
	WithoutGuide   int64
	WithNovelty    int64
	ByStatus       []StatusCount
	ShipmentStatus []StatusCount
}

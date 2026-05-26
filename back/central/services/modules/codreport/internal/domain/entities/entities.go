package entities

import "time"

type CarrierConfig struct {
	ID                 uint
	BusinessID         uint
	CarrierName        string
	DiscountPercentage float64
	IsActive           bool
}

type CodOrder struct {
	OrderID       string
	OrderNumber   string
	ShipmentID    uint
	CustomerName  string
	Carrier       string
	CodTotal      float64
	CodCarrierFee float64
	ShippingCost  float64
	DiscountPct   float64
	Discount      float64
	Net           float64
	Currency      string
	Status        string
	Collected     bool
	CreatedAt     time.Time
	DeliveredAt   *time.Time
	CutStatus     string
}

type CarrierAggregate struct {
	Carrier        string
	OrdersCount    int
	TotalCollected float64
	DiscountPct    float64
	TotalDiscount  float64
	TotalNet       float64
}

type MonthlyPoint struct {
	Month     string
	Label     string
	Orders    int
	Collected float64
	Discount  float64
	Net       float64
}

type CodSummary struct {
	TotalCollected  float64
	TotalPending    float64
	TotalDiscount   float64
	TotalNet        float64
	OrdersCollected int
	OrdersPending   int
	ByCarrier       []CarrierAggregate
	Monthly         []MonthlyPoint
}

type PaymentCut struct {
	ID              uint
	BusinessID      uint
	PeriodStart     time.Time
	PeriodEnd       time.Time
	Status          string
	OrdersCount     int
	TotalCollected  float64
	TotalDiscount   float64
	TotalNet        float64
	ByCarrier       []CarrierAggregate
	ConfirmedBy     uint
	ConfirmedByName string
	ConfirmedAt     *time.Time
}

type WeekAggregate struct {
	WeekStart time.Time
	Carrier   string
	Orders    int
	Collected float64
}

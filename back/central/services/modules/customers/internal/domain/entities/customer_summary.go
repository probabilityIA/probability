package entities

import "time"

type CustomerSummary struct {
	ID                uint
	CustomerID        uint
	BusinessID        uint
	TotalOrders       int
	DeliveredOrders   int
	CancelledOrders   int
	InProgressOrders  int
	TotalSpent        float64
	AvgTicket         float64
	TotalPaidOrders   int
	AvgDeliveryScore  float64
	FirstOrderAt      *time.Time
	LastOrderAt       *time.Time
	PreferredPlatform string
	LastUpdatedAt     time.Time
}

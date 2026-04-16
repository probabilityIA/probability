package entities

import "time"

type CustomerHistory struct {
	TotalOrders       int
	TotalSpent        float64
	AvgOrderValue     float64
	FirstOrderDate    *time.Time
	LastOrderDate     *time.Time
	NoveltyCount      int
	CODOrderCount     int
	DistinctAddresses int
	FailedPayments    int
}

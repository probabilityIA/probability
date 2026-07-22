package domain

import "time"

type IntegrationStats struct {
	IntegrationID    uint
	OrdersCount      int64
	OrdersInProgress int64
	OrdersDelivered  int64
	OrdersCancelled  int64
	OrdersReturned   int64
	ProductsCount    int64
	LastOrderAt      *time.Time
}

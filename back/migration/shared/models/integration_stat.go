package models

import "time"

type IntegrationStat struct {
	IntegrationID    uint  `gorm:"primaryKey" json:"integration_id"`
	BusinessID       uint  `gorm:"not null;index" json:"business_id"`
	OrdersTotal      int64 `gorm:"not null;default:0" json:"orders_total"`
	OrdersInProgress int64 `gorm:"not null;default:0" json:"orders_in_progress"`
	OrdersDelivered  int64 `gorm:"not null;default:0" json:"orders_delivered"`
	OrdersCancelled  int64 `gorm:"not null;default:0" json:"orders_cancelled"`
	OrdersReturned   int64 `gorm:"not null;default:0" json:"orders_returned"`
	ProductsCount    int64 `gorm:"not null;default:0" json:"products_count"`
	LastOrderAt      *time.Time
	UpdatedAt        time.Time
}

func (IntegrationStat) TableName() string {
	return "integration_stats"
}

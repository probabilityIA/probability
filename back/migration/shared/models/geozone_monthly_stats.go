package models

import "time"

type GeozoneMonthlyStats struct {
	ID              uint64    `gorm:"primaryKey"`
	BusinessID      uint      `gorm:"not null;index;uniqueIndex:uq_geozone_monthly_stats,priority:1"`
	Period          time.Time `gorm:"type:date;not null;uniqueIndex:uq_geozone_monthly_stats,priority:2"`
	GeozoneID       uint      `gorm:"not null;uniqueIndex:uq_geozone_monthly_stats,priority:3"`
	GeozoneType     string    `gorm:"size:32;not null"`
	Carrier         string    `gorm:"size:128;not null;default:'';uniqueIndex:uq_geozone_monthly_stats,priority:4"`
	TotalShipments  int       `gorm:"not null;default:0"`
	Delivered       int       `gorm:"not null;default:0"`
	Cancelled       int       `gorm:"not null;default:0"`
	Returned        int       `gorm:"not null;default:0"`
	InTransit       int       `gorm:"not null;default:0"`
	Failed          int       `gorm:"not null;default:0"`
	TotalAttempts   int       `gorm:"not null;default:0"`
	AvgDeliveryDays *float64  `gorm:"type:decimal(6,2)"`
	ComputedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (GeozoneMonthlyStats) TableName() string { return "geozone_monthly_stats" }

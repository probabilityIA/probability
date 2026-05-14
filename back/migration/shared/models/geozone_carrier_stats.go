package models

import "time"

type GeozoneCarrierStats struct {
	GeozoneLevel     string    `gorm:"size:20;primaryKey"`
	GeozoneID        uint64    `gorm:"primaryKey"`
	CarrierKey       string    `gorm:"size:64;primaryKey"`
	CarrierDisplay   string    `gorm:"size:128;not null;default:''"`
	Total            int64     `gorm:"not null;default:0"`
	Delivered        int64     `gorm:"not null;default:0"`
	Cancelled        int64     `gorm:"not null;default:0"`
	Returned         int64     `gorm:"not null;default:0"`
	Failed           int64     `gorm:"not null;default:0"`
	InTransit        int64     `gorm:"not null;default:0"`
	SampleSufficient bool      `gorm:"not null;default:false"`
	LastRefreshedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (GeozoneCarrierStats) TableName() string { return "geozone_carrier_stats" }

package models

import (
	"time"

	"gorm.io/gorm"
)

type CodPaymentCut struct {
	gorm.Model
	BusinessID       uint      `gorm:"not null;index;uniqueIndex:idx_cod_cut_biz_period,priority:1"`
	PeriodStart      time.Time `gorm:"type:date;not null;uniqueIndex:idx_cod_cut_biz_period,priority:2"`
	PeriodEnd        time.Time `gorm:"type:date;not null;uniqueIndex:idx_cod_cut_biz_period,priority:3"`
	Status           string    `gorm:"size:32;not null;default:'confirmed';index"`
	OrdersCount      int       `gorm:"not null;default:0"`
	TotalCollected   float64   `gorm:"type:decimal(15,2);not null;default:0"`
	TotalDiscount    float64   `gorm:"type:decimal(15,2);not null;default:0"`
	TotalNet         float64   `gorm:"type:decimal(15,2);not null;default:0"`
	CarrierBreakdown string    `gorm:"type:jsonb"`
	ConfirmedBy      uint      `gorm:"not null;default:0"`
	ConfirmedByName  string    `gorm:"size:160"`
	ConfirmedAt      *time.Time

	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

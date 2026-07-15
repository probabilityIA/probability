package models

import (
	"time"

	"gorm.io/gorm"
)

type CodPaymentCutOrder struct {
	gorm.Model
	CodPaymentCutID uint      `gorm:"not null;index"`
	BusinessID      uint      `gorm:"not null;index"`
	OrderID         string    `gorm:"type:varchar(36);not null;uniqueIndex:idx_cod_cut_order_order"`
	Carrier         string    `gorm:"size:120"`
	CodAmount       float64   `gorm:"type:decimal(15,2);not null;default:0"`
	PaidAt          time.Time `gorm:"not null"`

	CodPaymentCut CodPaymentCut `gorm:"foreignKey:CodPaymentCutID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

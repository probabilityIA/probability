package models

import (
	"time"

	"gorm.io/gorm"
)

// ───────────────────────────────────────────
//
//	BUSINESS SUBSCRIPTIONS - Historial de pagos y estado de suscripción
//
// ───────────────────────────────────────────
type BusinessSubscription struct {
	gorm.Model
	BusinessID       uint    `gorm:"not null;index"`
	Amount           float64 `gorm:"type:numeric(12,2);not null"`
	StartDate        time.Time
	EndDate          time.Time
	Status           string  `gorm:"size:20;default:'pending'"` // 'paid', 'pending', 'rejected'
	PaymentReference *string `gorm:"size:255"`                  // Ref de pago o comprobante
	Notes            *string `gorm:"type:text"`

	// Relaciones
	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName especifica el nombre de la tabla
func (BusinessSubscription) TableName() string {
	return "business_subscriptions"
}

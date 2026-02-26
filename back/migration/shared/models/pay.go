package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ───────────────────────────────────────────
//
//	PAYMENT TRANSACTIONS - Transacciones de pago via pasarelas externas
//
// ───────────────────────────────────────────
type PaymentTransaction struct {
	gorm.Model
	BusinessID      uint           `gorm:"not null;index"`
	Amount          float64        `gorm:"not null"`
	Currency        string         `gorm:"size:10;not null;default:'COP'"`
	Status          string         `gorm:"size:30;not null;index"` // pending|processing|completed|failed|cancelled
	GatewayCode     string         `gorm:"size:50;not null;index"` // "nequi"
	ExternalID      *string        `gorm:"size:255"`               // ID de la transacción en el gateway
	Reference       string         `gorm:"size:100;not null;uniqueIndex"`
	PaymentMethod   string         `gorm:"size:50"` // "qr_code"|"payment_link"
	Description     string         `gorm:"size:500"`
	CallbackURL     *string        `gorm:"size:500"`
	Metadata        datatypes.JSON `gorm:"type:jsonb"`
	GatewayResponse datatypes.JSON `gorm:"type:jsonb"`

	// Relaciones
	Business Business         `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SyncLogs []PaymentSyncLog `gorm:"foreignKey:PaymentTransactionID"`
}

// PaymentSyncLog registra cada intento de cobro
type PaymentSyncLog struct {
	gorm.Model
	PaymentTransactionID uint           `gorm:"not null;index"`
	Status               string         `gorm:"size:30;not null"` // processing|completed|failed|cancelled
	RetryCount           int            `gorm:"default:0"`
	GatewayRequest       datatypes.JSON `gorm:"type:jsonb"`
	GatewayResponse      datatypes.JSON `gorm:"type:jsonb"`
	ErrorMessage         *string        `gorm:"type:text"`
	NextRetryAt          *time.Time

	// Relaciones
	PaymentTransaction PaymentTransaction `gorm:"foreignKey:PaymentTransactionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

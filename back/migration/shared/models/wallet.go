package models

import (
	"time"

	"github.com/google/uuid"
)

// ───────────────────────────────────────────
//
//	WALLET - Billetera por negocio
//
// ───────────────────────────────────────────

type Wallet struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BusinessID uint      `gorm:"not null;uniqueIndex"`
	Balance    float64   `gorm:"type:decimal(15,2);default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Relaciones
	Business Business `gorm:"foreignKey:BusinessID"`
}

func (Wallet) TableName() string { return "wallet" }

// WalletTransaction - Transacciones de la billetera (recargas, débitos)
type WalletTransaction struct {
	ID                   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	WalletID             uuid.UUID `gorm:"type:uuid;not null;index"`
	Amount               float64   `gorm:"type:decimal(15,2);not null"`
	Type                 string    `gorm:"type:varchar(50);not null"`          // RECHARGE|USAGE
	Status               string    `gorm:"type:varchar(50);not null;default:'PENDING'"` // PENDING|COMPLETED|FAILED
	Reference            string    `gorm:"type:varchar(255)"`
	QrCode               string    `gorm:"type:text"`
	PaymentTransactionID *uint     `gorm:"index"` // FK opcional a payment_transactions
	CreatedAt            time.Time

	// Relaciones
	Wallet Wallet `gorm:"foreignKey:WalletID"`
}

func (WalletTransaction) TableName() string { return "transaction" }

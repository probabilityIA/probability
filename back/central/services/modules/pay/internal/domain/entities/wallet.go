package entities

import (
	"time"

	"github.com/google/uuid"
)

// Wallet tipos de transacción
const (
	WalletTxTypeRecharge = "RECHARGE"
	WalletTxTypeUsage    = "USAGE"
)

// Wallet estados de transacción
const (
	WalletTxStatusPending   = "PENDING"
	WalletTxStatusCompleted = "COMPLETED"
	WalletTxStatusFailed    = "FAILED"
)

// Wallet es la billetera de un negocio
type Wallet struct {
	ID         uuid.UUID
	BusinessID uint
	Balance    float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// WalletTransaction es un movimiento de la billetera
type WalletTransaction struct {
	ID                   uuid.UUID
	WalletID             uuid.UUID
	Amount               float64
	Type                 string // RECHARGE|USAGE
	Status               string // PENDING|COMPLETED|FAILED
	Reference            string
	QrCode               string
	PaymentTransactionID *uint
	CreatedAt            time.Time
}

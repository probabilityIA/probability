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

// Wallet conceptos (categoría) de transacción
const (
	WalletTxConceptGuide        = "GUIDE"
	WalletTxConceptSubscription = "SUBSCRIPTION"
	WalletTxConceptExtraUsage   = "EXTRA_USAGE"
	WalletTxConceptRecharge     = "RECHARGE"
	WalletTxConceptRefund       = "REFUND"
	WalletTxConceptAdjustment   = "ADJUSTMENT"
	WalletTxConceptOther        = "OTHER"
)

// ValidWalletTxConcept indica si el concepto es uno de los permitidos
func ValidWalletTxConcept(c string) bool {
	switch c {
	case WalletTxConceptGuide, WalletTxConceptSubscription, WalletTxConceptExtraUsage,
		WalletTxConceptRecharge, WalletTxConceptRefund, WalletTxConceptAdjustment, WalletTxConceptOther:
		return true
	default:
		return false
	}
}

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
	Concept              string // GUIDE|SUBSCRIPTION|EXTRA_USAGE|RECHARGE|REFUND|ADJUSTMENT|OTHER
	Reference            string
	QrCode               string
	PaymentTransactionID *uint
	UserID               *uint
	IntegrationTypeID    *uint
	IntegrationID        *uint
	GatewayRequest       []byte
	GatewayResponse      []byte
	IntegrationImageURL  string
	CreatedAt            time.Time
	BusinessID           uint
}

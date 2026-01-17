package domain

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BusinessID uint      `gorm:"not null;uniqueIndex"`
	Balance    float64   `gorm:"type:decimal(15,2);default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type TransactionType string

const (
	TransactionTypeRecharge TransactionType = "RECHARGE"
	TransactionTypeUsage    TransactionType = "USAGE"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
)

type Transaction struct {
	ID        uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	WalletID  uuid.UUID         `gorm:"type:uuid;not null;index"`
	Amount    float64           `gorm:"type:decimal(15,2);not null"`
	Type      TransactionType   `gorm:"type:varchar(50);not null"`
	Status    TransactionStatus `gorm:"type:varchar(50);not null;default:'PENDING'"`
	Reference string            `gorm:"type:varchar(255)"` // Nequi Transaction ID
	QrCode    string            `gorm:"type:text"`         // The QR string
	CreatedAt time.Time
}

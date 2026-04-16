package entities

import "time"

// PaymentSyncLog registra cada intento de procesamiento de pago
type PaymentSyncLog struct {
	ID                   uint
	PaymentTransactionID uint
	Status               string
	RetryCount           int
	ErrorMessage         *string
	NextRetryAt          *time.Time
	CreatedAt            time.Time
}

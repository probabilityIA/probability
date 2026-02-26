package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
)

// IRepository define todas las operaciones de persistencia del módulo de pagos
type IRepository interface {
	// PaymentTransactions
	CreatePaymentTransaction(ctx context.Context, tx *entities.PaymentTransaction) error
	GetPaymentTransactionByID(ctx context.Context, id uint) (*entities.PaymentTransaction, error)
	GetPaymentTransactionByReference(ctx context.Context, ref string) (*entities.PaymentTransaction, error)
	UpdatePaymentTransaction(ctx context.Context, tx *entities.PaymentTransaction) error
	ListPaymentTransactions(ctx context.Context, businessID uint, page, pageSize int) ([]*entities.PaymentTransaction, int64, error)

	// PaymentSyncLogs
	CreateSyncLog(ctx context.Context, log *entities.PaymentSyncLog) error
	UpdateSyncLog(ctx context.Context, log *entities.PaymentSyncLog) error
	GetPendingSyncLogRetries(ctx context.Context, limit int) ([]*entities.PaymentSyncLog, error)
	CancelPendingSyncLogs(ctx context.Context, transactionID uint) error
	GetSyncLogsByTransactionID(ctx context.Context, transactionID uint) ([]*entities.PaymentSyncLog, error)
}

// IRequestPublisher publica solicitudes de pago a la cola pay.requests
type IRequestPublisher interface {
	PublishPaymentRequest(ctx context.Context, msg *dtos.PaymentRequestMessage) error
}

// ISSEPublisher publica actualizaciones de pago a Redis Pub/Sub
type ISSEPublisher interface {
	PublishPaymentCompleted(ctx context.Context, tx *entities.PaymentTransaction) error
	PublishPaymentFailed(ctx context.Context, tx *entities.PaymentTransaction, errMsg string) error
	PublishPaymentProcessing(ctx context.Context, tx *entities.PaymentTransaction) error
}

// IUseCase define todos los casos de uso del módulo de pagos
type IUseCase interface {
	CreatePayment(ctx context.Context, dto *dtos.CreatePaymentDTO) (*entities.PaymentTransaction, error)
	ProcessPaymentResponse(ctx context.Context, msg *dtos.PaymentResponseMessage) error
	RetryPayment(ctx context.Context, transactionID uint) error
	GetPayment(ctx context.Context, id uint) (*entities.PaymentTransaction, error)
	ListPayments(ctx context.Context, businessID uint, page, pageSize int) ([]*entities.PaymentTransaction, int64, error)
}

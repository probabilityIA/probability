package ports

import (
	"context"

	"github.com/google/uuid"
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

	// Wallet
	GetWalletByBusinessID(ctx context.Context, businessID uint) (*entities.Wallet, error)
	GetWalletByID(ctx context.Context, walletID uuid.UUID) (*entities.Wallet, error)
	CreateWallet(ctx context.Context, wallet *entities.Wallet) error
	UpdateWallet(ctx context.Context, wallet *entities.Wallet) error
	GetAllWallets(ctx context.Context) ([]*entities.Wallet, error)

	// WalletTransactions
	CreateWalletTransaction(ctx context.Context, tx *entities.WalletTransaction) error
	GetWalletTransactionByID(ctx context.Context, id uuid.UUID) (*entities.WalletTransaction, error)
	GetWalletTransactionByReference(ctx context.Context, reference string) (*entities.WalletTransaction, error)
	SaveWalletTransactionGatewayResponse(ctx context.Context, id uuid.UUID, response []byte) error
	UpdateWalletTransaction(ctx context.Context, tx *entities.WalletTransaction) error
	GetTransactionsByWalletID(ctx context.Context, walletID uuid.UUID) ([]*entities.WalletTransaction, error)
	GetPendingRechargeTransactions(ctx context.Context) ([]*entities.WalletTransaction, error)
	GetProcessedTransactions(ctx context.Context) ([]*entities.WalletTransaction, error)
	DeleteTransactionsByWalletIDAndType(ctx context.Context, walletID uuid.UUID, txType string) error
	DeleteAllTransactionsByWalletID(ctx context.Context, walletID uuid.UUID) error

	GetFinancialStats(ctx context.Context, dto *dtos.FinancialStatsDTO) (*dtos.FinancialStatsResponse, error)

	GetBoldCredentials(ctx context.Context) (*dtos.BoldCredentials, error)
	GetBoldCredentialsForBusiness(ctx context.Context, businessID uint) (*dtos.BoldCredentials, error)
	GetBoldLinkCredentialsForBusiness(ctx context.Context, businessID uint) (*dtos.BoldCredentials, error)
	GetBoldIntegrationForBusiness(ctx context.Context, businessID uint) (*dtos.BoldBusinessIntegration, error)

	RecordBoldWebhookEvent(ctx context.Context, event *dtos.BoldWebhookEvent) (created bool, err error)
	MarkBoldWebhookProcessed(ctx context.Context, id uuid.UUID, paymentTransactionID *uint, processErr error) error
	LinkBoldWebhookToWalletTransaction(ctx context.Context, eventID, walletTransactionID uuid.UUID) error
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
	ProcessBoldWebhookMessage(ctx context.Context, msg *dtos.BoldWebhookMessage) error
	ApplyWalletRechargeStatus(ctx context.Context, in *dtos.WalletRechargeStatusInput) error
}

// IWalletUseCase define los casos de uso de la billetera
type IWalletUseCase interface {
	GetWallet(ctx context.Context, businessID uint) (*entities.Wallet, error)
	RechargeWallet(ctx context.Context, dto *dtos.RechargeWalletDTO) (*entities.WalletTransaction, error)
	ApproveTransaction(ctx context.Context, transactionID string) error
	RejectTransaction(ctx context.Context, transactionID string) error
	ManualDebit(ctx context.Context, dto *dtos.ManualDebitDTO) error
	DebitForGuide(ctx context.Context, dto *dtos.DebitForGuideDTO) error
	GetAllWallets(ctx context.Context) ([]*entities.Wallet, error)
	GetPendingTransactions(ctx context.Context) ([]*entities.WalletTransaction, error)
	GetProcessedTransactions(ctx context.Context) ([]*entities.WalletTransaction, error)
	GetTransactionsByBusinessID(ctx context.Context, businessID uint) ([]*entities.WalletTransaction, error)
	ClearRechargeHistory(ctx context.Context, businessID uint) error
	AdminAdjustBalance(ctx context.Context, dto *dtos.AdminAdjustBalanceDTO) error
	GetFinancialStats(ctx context.Context, dto *dtos.FinancialStatsDTO) (*dtos.FinancialStatsResponse, error)

	BoldGenerateSignature(ctx context.Context, businessID uint, amount float64, currency string) (*dtos.BoldSignatureResponse, error)
	GetBoldStatus(ctx context.Context, boldOrderID string) (*dtos.BoldStatusResponse, error)
	SyncBoldRecharge(ctx context.Context, businessID uint, orderID string) (*dtos.BoldStatusResponse, error)
	BoldSimulatePayment(ctx context.Context, dto *dtos.BoldSimulateDTO) (*dtos.BoldSimulateResponse, error)
}

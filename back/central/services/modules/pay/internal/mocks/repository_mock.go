package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) CreatePaymentTransaction(ctx context.Context, tx *entities.PaymentTransaction) error {
	return m.Called(ctx, tx).Error(0)
}

func (m *RepositoryMock) GetPaymentTransactionByID(ctx context.Context, id uint) (*entities.PaymentTransaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.PaymentTransaction), args.Error(1)
}

func (m *RepositoryMock) GetPaymentTransactionByReference(ctx context.Context, ref string) (*entities.PaymentTransaction, error) {
	args := m.Called(ctx, ref)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.PaymentTransaction), args.Error(1)
}

func (m *RepositoryMock) UpdatePaymentTransaction(ctx context.Context, tx *entities.PaymentTransaction) error {
	return m.Called(ctx, tx).Error(0)
}

func (m *RepositoryMock) GetStorefrontCheckoutByReference(ctx context.Context, reference string) (*entities.StorefrontCheckoutSnapshot, error) {
	args := m.Called(ctx, reference)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.StorefrontCheckoutSnapshot), args.Error(1)
}

func (m *RepositoryMock) MarkStorefrontCheckoutStatus(ctx context.Context, reference, status string) error {
	return m.Called(ctx, reference, status).Error(0)
}

func (m *RepositoryMock) ListPaymentTransactions(ctx context.Context, businessID uint, page, pageSize int) ([]*entities.PaymentTransaction, int64, error) {
	args := m.Called(ctx, businessID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entities.PaymentTransaction), args.Get(1).(int64), args.Error(2)
}

func (m *RepositoryMock) CreateSyncLog(ctx context.Context, l *entities.PaymentSyncLog) error {
	return m.Called(ctx, l).Error(0)
}

func (m *RepositoryMock) UpdateSyncLog(ctx context.Context, l *entities.PaymentSyncLog) error {
	return m.Called(ctx, l).Error(0)
}

func (m *RepositoryMock) GetPendingSyncLogRetries(ctx context.Context, limit int) ([]*entities.PaymentSyncLog, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.PaymentSyncLog), args.Error(1)
}

func (m *RepositoryMock) CancelPendingSyncLogs(ctx context.Context, transactionID uint) error {
	return m.Called(ctx, transactionID).Error(0)
}

func (m *RepositoryMock) GetSyncLogsByTransactionID(ctx context.Context, transactionID uint) ([]*entities.PaymentSyncLog, error) {
	args := m.Called(ctx, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.PaymentSyncLog), args.Error(1)
}

func (m *RepositoryMock) GetWalletByBusinessID(ctx context.Context, businessID uint) (*entities.Wallet, error) {
	args := m.Called(ctx, businessID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Wallet), args.Error(1)
}

func (m *RepositoryMock) GetWalletByID(ctx context.Context, walletID uuid.UUID) (*entities.Wallet, error) {
	args := m.Called(ctx, walletID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Wallet), args.Error(1)
}

func (m *RepositoryMock) CreateWallet(ctx context.Context, wallet *entities.Wallet) error {
	return m.Called(ctx, wallet).Error(0)
}

func (m *RepositoryMock) UpdateWallet(ctx context.Context, wallet *entities.Wallet) error {
	return m.Called(ctx, wallet).Error(0)
}

func (m *RepositoryMock) GetAllWallets(ctx context.Context) ([]*entities.Wallet, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Wallet), args.Error(1)
}

func (m *RepositoryMock) CreateWalletTransaction(ctx context.Context, tx *entities.WalletTransaction) error {
	return m.Called(ctx, tx).Error(0)
}

func (m *RepositoryMock) GetWalletTransactionByID(ctx context.Context, id uuid.UUID) (*entities.WalletTransaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.WalletTransaction), args.Error(1)
}

func (m *RepositoryMock) GetWalletTransactionByReference(ctx context.Context, reference string) (*entities.WalletTransaction, error) {
	args := m.Called(ctx, reference)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.WalletTransaction), args.Error(1)
}

func (m *RepositoryMock) SaveWalletTransactionGatewayResponse(ctx context.Context, id uuid.UUID, response []byte) error {
	return m.Called(ctx, id, response).Error(0)
}

func (m *RepositoryMock) UpdateWalletTransaction(ctx context.Context, tx *entities.WalletTransaction) error {
	return m.Called(ctx, tx).Error(0)
}

func (m *RepositoryMock) GetTransactionsByWalletID(ctx context.Context, walletID uuid.UUID) ([]*entities.WalletTransaction, error) {
	args := m.Called(ctx, walletID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.WalletTransaction), args.Error(1)
}

func (m *RepositoryMock) GetPendingRechargeTransactions(ctx context.Context) ([]*entities.WalletTransaction, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.WalletTransaction), args.Error(1)
}

func (m *RepositoryMock) GetProcessedTransactions(ctx context.Context) ([]*entities.WalletTransaction, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.WalletTransaction), args.Error(1)
}

func (m *RepositoryMock) DeleteTransactionsByWalletIDAndType(ctx context.Context, walletID uuid.UUID, txType string) error {
	return m.Called(ctx, walletID, txType).Error(0)
}

func (m *RepositoryMock) DeleteAllTransactionsByWalletID(ctx context.Context, walletID uuid.UUID) error {
	return m.Called(ctx, walletID).Error(0)
}

func (m *RepositoryMock) GetFinancialStats(ctx context.Context, dto *dtos.FinancialStatsDTO) (*dtos.FinancialStatsResponse, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.FinancialStatsResponse), args.Error(1)
}

func (m *RepositoryMock) GetBoldCredentials(ctx context.Context) (*dtos.BoldCredentials, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.BoldCredentials), args.Error(1)
}

func (m *RepositoryMock) GetBoldCredentialsForBusiness(ctx context.Context, businessID uint) (*dtos.BoldCredentials, error) {
	args := m.Called(ctx, businessID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.BoldCredentials), args.Error(1)
}

func (m *RepositoryMock) GetBoldLinkCredentialsForBusiness(ctx context.Context, businessID uint) (*dtos.BoldCredentials, error) {
	args := m.Called(ctx, businessID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.BoldCredentials), args.Error(1)
}

func (m *RepositoryMock) GetBoldIntegrationForBusiness(ctx context.Context, businessID uint) (*dtos.BoldBusinessIntegration, error) {
	args := m.Called(ctx, businessID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtos.BoldBusinessIntegration), args.Error(1)
}

func (m *RepositoryMock) RecordPaymentWebhookEvent(ctx context.Context, event *dtos.PaymentWebhookEvent) (bool, error) {
	args := m.Called(ctx, event)
	return args.Bool(0), args.Error(1)
}

func (m *RepositoryMock) MarkPaymentWebhookProcessed(ctx context.Context, id uuid.UUID, paymentTransactionID *uint, processErr error) error {
	return m.Called(ctx, id, paymentTransactionID, processErr).Error(0)
}

func (m *RepositoryMock) LinkPaymentWebhookToWalletTransaction(ctx context.Context, eventID, walletTransactionID uuid.UUID) error {
	return m.Called(ctx, eventID, walletTransactionID).Error(0)
}

func (m *RepositoryMock) ListLowBalanceAlertCandidates(ctx context.Context, threshold float64) ([]entities.LowBalanceAlertCandidate, error) {
	args := m.Called(ctx, threshold)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.LowBalanceAlertCandidate), args.Error(1)
}

func (m *RepositoryMock) GetLowBalanceAlertCandidate(ctx context.Context, businessID uint, threshold float64) (*entities.LowBalanceAlertCandidate, error) {
	args := m.Called(ctx, businessID, threshold)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.LowBalanceAlertCandidate), args.Error(1)
}

func (m *RepositoryMock) HasRecentSystemAlert(ctx context.Context, phoneNumber string) (bool, error) {
	args := m.Called(ctx, phoneNumber)
	return args.Bool(0), args.Error(1)
}

func (m *RepositoryMock) GetWalletKPISelection(ctx context.Context) (*entities.WalletKPISelection, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.WalletKPISelection), args.Error(1)
}

func (m *RepositoryMock) UpdateWalletKPISelection(ctx context.Context, selection *entities.WalletKPISelection) error {
	return m.Called(ctx, selection).Error(0)
}

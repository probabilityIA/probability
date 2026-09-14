package app

import (
	"context"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeQueue struct {
	published []string
}

func (f *fakeQueue) Publish(ctx context.Context, queueName string, message []byte) error {
	f.published = append(f.published, queueName)
	return nil
}
func (f *fakeQueue) PublishToExchange(ctx context.Context, exchangeName, routingKey string, message []byte) error {
	return nil
}
func (f *fakeQueue) Consume(ctx context.Context, queueName string, handler func([]byte) error) error {
	return nil
}
func (f *fakeQueue) ConsumeConcurrent(ctx context.Context, queueName string, handler func([]byte) error, workers int) error {
	return nil
}
func (f *fakeQueue) Close() error                                                    { return nil }
func (f *fakeQueue) DeclareQueue(queueName string, durable bool) error               { return nil }
func (f *fakeQueue) DeclareExchange(exchangeName, exchangeType string, durable bool) error {
	return nil
}
func (f *fakeQueue) BindQueue(queueName, exchangeName, routingKey string) error { return nil }
func (f *fakeQueue) Ping() error                                               { return nil }

func mysticStylePlan() *entities.SubscriptionType {
	return &entities.SubscriptionType{ID: 4, Code: "Plan-Membresia", Name: "Plan Membresia", Price: 99000, Active: true, Payable: true}
}

func TestAutoRenewIfEnabled_SaldoInsuficiente_QuedaRegistradoEnLaAuditoria(t *testing.T) {
	var auditLogs []*entities.SubscriptionAuditLog
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return mysticStylePlan(), nil
		},
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			return nil, nil
		},
		FindSuperAdminUserIDFn: func(ctx context.Context) (uint, error) { return 1, nil },
		CreateAuditLogFn: func(ctx context.Context, log *entities.SubscriptionAuditLog) error {
			auditLogs = append(auditLogs, log)
			return nil
		},
	}
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 762047.75 - 762047.75, nil },
	}

	uc := &UseCase{repo: repo, wallet: wallet, log: mocks.NewSilentLogger()}

	business := entities.ExpiringBusiness{
		BusinessID: 36, SubscriptionTypeID: 4, AutoPaymentEnabled: true,
		EndDate: time.Date(2026, 9, 11, 0, 0, 0, 0, time.UTC),
	}

	renewed := uc.autoRenewIfEnabled(context.Background(), business)

	assert.False(t, renewed, "sin saldo no debe renovar")
	require.Len(t, auditLogs, 1, "debe quedar rastro de por que no se renovo")
	assert.Equal(t, entities.AuditActionAutoRenewInsufficientFunds, auditLogs[0].Action)
	assert.Contains(t, auditLogs[0].Description, "99000")
}

func TestAutoRenewIfEnabled_SaldoSuficiente_NoRegistraFalloDeSaldo(t *testing.T) {
	var auditLogs []*entities.SubscriptionAuditLog
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return mysticStylePlan(), nil
		},
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			return nil, nil
		},
		FindSuperAdminUserIDFn: func(ctx context.Context) (uint, error) { return 1, nil },
		CreateAuditLogFn: func(ctx context.Context, log *entities.SubscriptionAuditLog) error {
			auditLogs = append(auditLogs, log)
			return nil
		},
		CreateSubscriptionAndActivateFn: func(ctx context.Context, sub *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error {
			return nil
		},
	}
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 762047.75, nil },
		DebitFn: func(ctx context.Context, businessID uint, amount float64, reference, concept string, userID uint) error {
			return nil
		},
	}

	uc := &UseCase{repo: repo, wallet: wallet, announcements: &mocks.AnnouncementsGatewayMock{}, log: mocks.NewSilentLogger()}

	business := entities.ExpiringBusiness{
		BusinessID: 36, SubscriptionTypeID: 4, AutoPaymentEnabled: true,
		EndDate: time.Date(2026, 9, 11, 0, 0, 0, 0, time.UTC),
	}

	renewed := uc.autoRenewIfEnabled(context.Background(), business)

	assert.True(t, renewed)
	for _, l := range auditLogs {
		assert.NotEqual(t, entities.AuditActionAutoRenewInsufficientFunds, l.Action,
			"con saldo suficiente no debe quedar un log de saldo insuficiente")
	}
}

func TestNotifyPaymentWindowIfNeeded_YaAvisadoHoy_NoReenviaElMismoDia(t *testing.T) {
	queue := &fakeQueue{}
	repo := &mocks.RepositoryMock{
		HasAuditLogSinceFn: func(ctx context.Context, businessID uint, action string, since time.Time) (bool, error) {
			return true, nil
		},
	}
	uc := &UseCase{repo: repo, rabbit: queue, log: mocks.NewSilentLogger()}

	business := entities.ExpiringBusiness{BusinessID: 36, EndDate: time.Date(2026, 9, 11, 0, 0, 0, 0, time.UTC)}
	now := time.Date(2026, 9, 14, 13, 0, 0, 0, time.UTC)

	uc.notifyPaymentWindowIfNeeded(context.Background(), business, now)

	assert.Empty(t, queue.published, "ya se aviso hoy, no debe reenviar")
}

func TestNotifyPaymentWindowIfNeeded_AvisadoAyer_ReenviaHoy(t *testing.T) {
	queue := &fakeQueue{}
	repo := &mocks.RepositoryMock{
		HasAuditLogSinceFn: func(ctx context.Context, businessID uint, action string, since time.Time) (bool, error) {
			return false, nil
		},
		GetWhatsAppContactFn: func(ctx context.Context, businessID uint) (*entities.WhatsAppContact, error) {
			return &entities.WhatsAppContact{Phone: "3001234567", BusinessName: "Mystic Rose"}, nil
		},
		FindSuperAdminUserIDFn: func(ctx context.Context) (uint, error) { return 1, nil },
		CreateAuditLogFn: func(ctx context.Context, log *entities.SubscriptionAuditLog) error {
			return nil
		},
	}
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 762047.75, nil },
	}
	uc := &UseCase{repo: repo, wallet: wallet, rabbit: queue, log: mocks.NewSilentLogger()}

	business := entities.ExpiringBusiness{BusinessID: 36, EndDate: time.Date(2026, 9, 11, 0, 0, 0, 0, time.UTC)}
	now := time.Date(2026, 9, 14, 13, 0, 0, 0, time.UTC)

	uc.notifyPaymentWindowIfNeeded(context.Background(), business, now)

	require.Len(t, queue.published, 1, "el aviso de ayer no debe bloquear el de hoy")
}

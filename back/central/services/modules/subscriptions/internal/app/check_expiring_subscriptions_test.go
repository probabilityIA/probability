package app

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type auditLogStore struct {
	mu   sync.Mutex
	logs []*entities.SubscriptionAuditLog
}

func (s *auditLogStore) create(ctx context.Context, log *entities.SubscriptionAuditLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.CreatedAt = time.Now()
	s.logs = append(s.logs, log)
	return nil
}

func (s *auditLogStore) hasSince(businessID uint, action string, since time.Time) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, l := range s.logs {
		if l.BusinessID == businessID && l.Action == action && !l.CreatedAt.Before(since) {
			return true
		}
	}
	return false
}

func (s *auditLogStore) count(action string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for _, l := range s.logs {
		if l.Action == action {
			n++
		}
	}
	return n
}

func TestCheckExpiringSubscriptions_SaldoSuficiente_RenuevaSinAvisarNiExpirar(t *testing.T) {
	store := &auditLogStore{}
	activated := false
	markedExpired := false

	business := entities.ExpiringBusiness{
		BusinessID: 36, SubscriptionTypeID: 4, AutoPaymentEnabled: true,
		EndDate: time.Now().Add(-3 * 24 * time.Hour),
	}

	repo := &mocks.RepositoryMock{
		ListBusinessesJustExpiredFn: func(ctx context.Context, before time.Time) ([]entities.ExpiringBusiness, error) {
			return []entities.ExpiringBusiness{business}, nil
		},
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return mysticStylePlan(), nil
		},
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			return nil, nil
		},
		FindSuperAdminUserIDFn: func(ctx context.Context) (uint, error) { return 1, nil },
		CreateAuditLogFn:       store.create,
		CreateSubscriptionAndActivateFn: func(ctx context.Context, sub *entities.BusinessSubscription, subscriptionTypeID uint, endDate time.Time) error {
			activated = true
			return nil
		},
		MarkExpiredIfStillActiveFn: func(ctx context.Context, businessID uint, before time.Time) error {
			markedExpired = true
			return nil
		},
	}
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 762047.75, nil },
		DebitFn: func(ctx context.Context, businessID uint, amount float64, reference, concept string, userID uint) error {
			return nil
		},
	}
	queue := &fakeQueue{}

	uc := &UseCase{repo: repo, wallet: wallet, announcements: &mocks.AnnouncementsGatewayMock{}, rabbit: queue, log: mocks.NewSilentLogger()}

	err := uc.CheckExpiringSubscriptions(context.Background())

	require.NoError(t, err)
	assert.True(t, activated, "con saldo suficiente debe renovar")
	assert.Empty(t, queue.published, "si renovo, no debe avisar por whatsapp")
	assert.False(t, markedExpired, "no debe marcarse expirado si se renovo")
	assert.Equal(t, 0, store.count(entities.AuditActionAutoRenewInsufficientFunds))
}

func TestCheckExpiringSubscriptions_SaldoInsuficiente_AvisaYNoExpiraDentroDeLaGracia(t *testing.T) {
	store := &auditLogStore{}
	markedExpired := false

	business := entities.ExpiringBusiness{
		BusinessID: 36, SubscriptionTypeID: 4, AutoPaymentEnabled: true,
		EndDate: time.Now().Add(-3 * 24 * time.Hour),
	}

	repo := &mocks.RepositoryMock{
		ListBusinessesJustExpiredFn: func(ctx context.Context, before time.Time) ([]entities.ExpiringBusiness, error) {
			return []entities.ExpiringBusiness{business}, nil
		},
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return mysticStylePlan(), nil
		},
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			return nil, nil
		},
		FindSuperAdminUserIDFn: func(ctx context.Context) (uint, error) { return 1, nil },
		CreateAuditLogFn:       store.create,
		HasAuditLogSinceFn: func(ctx context.Context, businessID uint, action string, since time.Time) (bool, error) {
			return store.hasSince(businessID, action, since), nil
		},
		GetWhatsAppContactFn: func(ctx context.Context, businessID uint) (*entities.WhatsAppContact, error) {
			return &entities.WhatsAppContact{Phone: "3001234567", BusinessName: "Mystic Rose"}, nil
		},
		MarkExpiredIfStillActiveFn: func(ctx context.Context, businessID uint, before time.Time) error {
			markedExpired = true
			return nil
		},
	}
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 1000, nil },
	}
	queue := &fakeQueue{}

	uc := &UseCase{repo: repo, wallet: wallet, announcements: &mocks.AnnouncementsGatewayMock{}, rabbit: queue, log: mocks.NewSilentLogger()}

	err := uc.CheckExpiringSubscriptions(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 1, store.count(entities.AuditActionAutoRenewInsufficientFunds), "debe quedar rastro del intento fallido")
	require.Len(t, queue.published, 1, "debe avisar por whatsapp")
	assert.False(t, markedExpired, "3 dias de vencido con gracia de 6 no debe expirar todavia")
}

func TestCheckExpiringSubscriptions_SaldoInsuficiente_ExpiraAlPasarLaGracia(t *testing.T) {
	store := &auditLogStore{}
	markedExpired := false

	business := entities.ExpiringBusiness{
		BusinessID: 36, SubscriptionTypeID: 4, AutoPaymentEnabled: true,
		EndDate: time.Now().Add(-10 * 24 * time.Hour),
	}

	repo := &mocks.RepositoryMock{
		ListBusinessesJustExpiredFn: func(ctx context.Context, before time.Time) ([]entities.ExpiringBusiness, error) {
			return []entities.ExpiringBusiness{business}, nil
		},
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return mysticStylePlan(), nil
		},
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			return nil, nil
		},
		FindSuperAdminUserIDFn: func(ctx context.Context) (uint, error) { return 1, nil },
		CreateAuditLogFn:       store.create,
		HasAuditLogSinceFn: func(ctx context.Context, businessID uint, action string, since time.Time) (bool, error) {
			return store.hasSince(businessID, action, since), nil
		},
		GetWhatsAppContactFn: func(ctx context.Context, businessID uint) (*entities.WhatsAppContact, error) {
			return &entities.WhatsAppContact{Phone: "3001234567", BusinessName: "Mystic Rose"}, nil
		},
		MarkExpiredIfStillActiveFn: func(ctx context.Context, businessID uint, before time.Time) error {
			markedExpired = true
			return nil
		},
	}
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 1000, nil },
	}
	uc := &UseCase{repo: repo, wallet: wallet, announcements: &mocks.AnnouncementsGatewayMock{}, rabbit: &fakeQueue{}, log: mocks.NewSilentLogger()}

	err := uc.CheckExpiringSubscriptions(context.Background())

	require.NoError(t, err)
	assert.True(t, markedExpired, "pasados los 6 dias de gracia sin pagar, debe marcarse expirado")
}

func TestCheckExpiringSubscriptions_CorridaDosVecesElMismoDia_NoDuplicaElAvisoDeWhatsapp(t *testing.T) {
	store := &auditLogStore{}

	business := entities.ExpiringBusiness{
		BusinessID: 36, SubscriptionTypeID: 4, AutoPaymentEnabled: true,
		EndDate: time.Now().Add(-3 * 24 * time.Hour),
	}

	repo := &mocks.RepositoryMock{
		ListBusinessesJustExpiredFn: func(ctx context.Context, before time.Time) ([]entities.ExpiringBusiness, error) {
			return []entities.ExpiringBusiness{business}, nil
		},
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return mysticStylePlan(), nil
		},
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			return nil, nil
		},
		FindSuperAdminUserIDFn: func(ctx context.Context) (uint, error) { return 1, nil },
		CreateAuditLogFn:       store.create,
		HasAuditLogSinceFn: func(ctx context.Context, businessID uint, action string, since time.Time) (bool, error) {
			return store.hasSince(businessID, action, since), nil
		},
		GetWhatsAppContactFn: func(ctx context.Context, businessID uint) (*entities.WhatsAppContact, error) {
			return &entities.WhatsAppContact{Phone: "3001234567", BusinessName: "Mystic Rose"}, nil
		},
	}
	wallet := &mocks.WalletDebiterMock{
		GetBalanceFn: func(ctx context.Context, businessID uint) (float64, error) { return 1000, nil },
	}
	queue := &fakeQueue{}
	uc := &UseCase{repo: repo, wallet: wallet, announcements: &mocks.AnnouncementsGatewayMock{}, rabbit: queue, log: mocks.NewSilentLogger()}

	require.NoError(t, uc.CheckExpiringSubscriptions(context.Background()))
	require.NoError(t, uc.CheckExpiringSubscriptions(context.Background()))

	assert.Len(t, queue.published, 1, "la corrida de las 8pm no debe repetir el aviso que ya mando la de las 8am")
}

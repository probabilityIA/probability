package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	domainErrors "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/mocks"
)

func TestSyncByIntegration_Success_RecachesEnabledConfigs(t *testing.T) {
	ctx := context.Background()
	businessID := uint(26)
	integrationID := uint(35)
	cachedConfigIDs := []uint{}

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			return []entities.IntegrationNotificationConfig{
				{ID: 1, IntegrationID: integrationID, NotificationTypeID: 2, NotificationEventTypeID: 3, Enabled: true},
				{ID: 2, IntegrationID: integrationID, NotificationTypeID: 2, NotificationEventTypeID: 4, Enabled: false},
			}, nil
		},
		SyncConfigsFn: func(ctx context.Context, bID uint, iID uint, toCreate []*entities.IntegrationNotificationConfig, toUpdate []*entities.IntegrationNotificationConfig, toDeleteIDs []uint) error {
			return nil
		},
	}

	mockCacheManager := &mocks.CacheManagerMock{
		InvalidateConfigsByIntegrationFn: func(ctx context.Context, id uint) error {
			return nil
		},
		CacheConfigFn: func(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
			cachedConfigIDs = append(cachedConfigIDs, config.ID)
			return nil
		},
	}

	mockLogger := mocks.NewLoggerMock()
	useCase := New(mockRepo, &mocks.NotificationTypeRepositoryMock{}, &mocks.NotificationEventTypeRepositoryMock{}, mockCacheManager, &mocks.MessageAuditQuerierMock{}, &mocks.AIPauseCheckerMock{}, mockLogger)

	dto := dtos.SyncNotificationConfigsDTO{
		BusinessID:    businessID,
		IntegrationID: integrationID,
		Rules: []dtos.SyncRuleDTO{
			{NotificationTypeID: 2, NotificationEventTypeID: 3, Enabled: true},
			{NotificationTypeID: 2, NotificationEventTypeID: 4, Enabled: false},
		},
	}

	result, err := useCase.SyncByIntegration(ctx, dto)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}

	// Solo la config enabled (ID=1) debe ser cacheada
	if len(cachedConfigIDs) != 1 {
		t.Fatalf("expected 1 cached config, got %d", len(cachedConfigIDs))
	}
	if cachedConfigIDs[0] != 1 {
		t.Errorf("expected cached config ID 1, got %d", cachedConfigIDs[0])
	}
}

func TestSyncByIntegration_CacheErrorShouldNotFail(t *testing.T) {
	ctx := context.Background()
	businessID := uint(26)
	integrationID := uint(35)

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			return []entities.IntegrationNotificationConfig{
				{ID: 1, IntegrationID: integrationID, NotificationTypeID: 2, NotificationEventTypeID: 3, Enabled: true},
			}, nil
		},
		SyncConfigsFn: func(ctx context.Context, bID uint, iID uint, toCreate []*entities.IntegrationNotificationConfig, toUpdate []*entities.IntegrationNotificationConfig, toDeleteIDs []uint) error {
			return nil
		},
	}

	mockCacheManager := &mocks.CacheManagerMock{
		InvalidateConfigsByIntegrationFn: func(ctx context.Context, id uint) error {
			return nil
		},
		CacheConfigFn: func(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
			return errors.New("redis connection failed")
		},
	}

	mockLogger := mocks.NewLoggerMock()
	useCase := New(mockRepo, &mocks.NotificationTypeRepositoryMock{}, &mocks.NotificationEventTypeRepositoryMock{}, mockCacheManager, &mocks.MessageAuditQuerierMock{}, &mocks.AIPauseCheckerMock{}, mockLogger)

	dto := dtos.SyncNotificationConfigsDTO{
		BusinessID:    businessID,
		IntegrationID: integrationID,
		Rules: []dtos.SyncRuleDTO{
			{NotificationTypeID: 2, NotificationEventTypeID: 3, Enabled: true},
		},
	}

	result, err := useCase.SyncByIntegration(ctx, dto)

	if err != nil {
		t.Fatalf("cache error should not fail sync, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestSyncByIntegration_DuplicateRulesError(t *testing.T) {
	ctx := context.Background()

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			return []entities.IntegrationNotificationConfig{}, nil
		},
	}

	mockLogger := mocks.NewLoggerMock()
	useCase := New(mockRepo, &mocks.NotificationTypeRepositoryMock{}, &mocks.NotificationEventTypeRepositoryMock{}, &mocks.CacheManagerMock{}, &mocks.MessageAuditQuerierMock{}, &mocks.AIPauseCheckerMock{}, mockLogger)

	dto := dtos.SyncNotificationConfigsDTO{
		BusinessID:    26,
		IntegrationID: 35,
		Rules: []dtos.SyncRuleDTO{
			{NotificationTypeID: 2, NotificationEventTypeID: 3, Enabled: true},
			{NotificationTypeID: 2, NotificationEventTypeID: 3, Enabled: false}, // Duplicado
		},
	}

	_, err := useCase.SyncByIntegration(ctx, dto)

	if err == nil {
		t.Fatal("expected error for duplicate rules, got nil")
	}
	if !errors.Is(err, domainErrors.ErrDuplicateConfig) {
		t.Errorf("expected ErrDuplicateConfig, got %v", err)
	}
}

func TestSyncByIntegration_SyncConfigsError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("transaction failed")

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			return []entities.IntegrationNotificationConfig{}, nil
		},
		SyncConfigsFn: func(ctx context.Context, bID uint, iID uint, toCreate []*entities.IntegrationNotificationConfig, toUpdate []*entities.IntegrationNotificationConfig, toDeleteIDs []uint) error {
			return expectedErr
		},
	}

	mockLogger := mocks.NewLoggerMock()
	useCase := New(mockRepo, &mocks.NotificationTypeRepositoryMock{}, &mocks.NotificationEventTypeRepositoryMock{}, &mocks.CacheManagerMock{}, &mocks.MessageAuditQuerierMock{}, &mocks.AIPauseCheckerMock{}, mockLogger)

	dto := dtos.SyncNotificationConfigsDTO{
		BusinessID:    26,
		IntegrationID: 35,
		Rules: []dtos.SyncRuleDTO{
			{NotificationTypeID: 2, NotificationEventTypeID: 3, Enabled: true},
		},
	}

	_, err := useCase.SyncByIntegration(ctx, dto)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

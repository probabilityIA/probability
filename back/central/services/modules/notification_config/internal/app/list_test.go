package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/mocks"
)

func TestList_Success_NoFilters(t *testing.T) {
	// Arrange
	ctx := context.Background()

	expectedConfigs := []entities.IntegrationNotificationConfig{
		{
			ID:                      1,
			IntegrationID:           100,
			NotificationTypeID:      1,
			NotificationEventTypeID: 5,
			Enabled:                 true,
			Description:             "Config WhatsApp order.created",
			OrderStatusIDs:          []uint{1, 2},
			CreatedAt:               time.Now().Add(-48 * time.Hour),
			UpdatedAt:               time.Now().Add(-24 * time.Hour),
		},
		{
			ID:                      2,
			IntegrationID:           200,
			NotificationTypeID:      2,
			NotificationEventTypeID: 6,
			Enabled:                 false,
			Description:             "Config Email order.shipped",
			OrderStatusIDs:          []uint{3, 4},
			CreatedAt:               time.Now().Add(-36 * time.Hour),
			UpdatedAt:               time.Now().Add(-12 * time.Hour),
		},
	}

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			return expectedConfigs, nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	// Act
	result, err := useCase.List(ctx, dtos.FilterNotificationConfigDTO{})

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != len(expectedConfigs) {
		t.Fatalf("expected %d configs, got %d", len(expectedConfigs), len(result))
	}

	if result[0].ID != expectedConfigs[0].ID {
		t.Errorf("expected first config ID %d, got %d", expectedConfigs[0].ID, result[0].ID)
	}

	if result[0].IntegrationID != expectedConfigs[0].IntegrationID {
		t.Errorf("expected IntegrationID %d, got %d", expectedConfigs[0].IntegrationID, result[0].IntegrationID)
	}

	if result[0].Enabled != expectedConfigs[0].Enabled {
		t.Errorf("expected Enabled %v, got %v", expectedConfigs[0].Enabled, result[0].Enabled)
	}
}

func TestList_Success_WithFilters(t *testing.T) {
	// Arrange
	ctx := context.Background()
	integrationID := uint(100)

	filteredConfigs := []entities.IntegrationNotificationConfig{
		{
			ID:                      1,
			IntegrationID:           100,
			NotificationTypeID:      1,
			NotificationEventTypeID: 5,
			Enabled:                 true,
			Description:             "Config filtrada",
			OrderStatusIDs:          []uint{1},
		},
	}

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			// Verificar que el filtro se pase correctamente
			if filters.IntegrationID != nil && *filters.IntegrationID == integrationID {
				return filteredConfigs, nil
			}
			return []entities.IntegrationNotificationConfig{}, nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	filters := dtos.FilterNotificationConfigDTO{
		IntegrationID: &integrationID,
	}

	// Act
	result, err := useCase.List(ctx, filters)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 config, got %d", len(result))
	}

	if result[0].IntegrationID != integrationID {
		t.Errorf("expected IntegrationID %d, got %d", integrationID, result[0].IntegrationID)
	}
}

func TestList_Success_EmptyResult(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			return []entities.IntegrationNotificationConfig{}, nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	// Act
	result, err := useCase.List(ctx, dtos.FilterNotificationConfigDTO{})

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected 0 configs, got %d", len(result))
	}
}

func TestList_RepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expectedErr := errors.New("database query failed")

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			return nil, expectedErr
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	// Act
	result, err := useCase.List(ctx, dtos.FilterNotificationConfigDTO{})

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestList_FilterByNotificationTypeID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	notificationTypeID := uint(1) // WhatsApp

	whatsappConfigs := []entities.IntegrationNotificationConfig{
		{
			ID:                      1,
			IntegrationID:           100,
			NotificationTypeID:      1, // WhatsApp
			NotificationEventTypeID: 5,
			Enabled:                 true,
		},
		{
			ID:                      2,
			IntegrationID:           101,
			NotificationTypeID:      1, // WhatsApp
			NotificationEventTypeID: 6,
			Enabled:                 true,
		},
	}

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			if filters.NotificationTypeID != nil && *filters.NotificationTypeID == notificationTypeID {
				return whatsappConfigs, nil
			}
			return []entities.IntegrationNotificationConfig{}, nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	filters := dtos.FilterNotificationConfigDTO{
		NotificationTypeID: &notificationTypeID,
	}

	// Act
	result, err := useCase.List(ctx, filters)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 configs, got %d", len(result))
	}

	// Verificar que todas sean de WhatsApp
	for _, config := range result {
		if config.NotificationTypeID != notificationTypeID {
			t.Errorf("expected NotificationTypeID %d, got %d", notificationTypeID, config.NotificationTypeID)
		}
	}
}

func TestList_FilterByEnabled(t *testing.T) {
	// Arrange
	ctx := context.Background()
	enabledFilter := true

	activeConfigs := []entities.IntegrationNotificationConfig{
		{
			ID:                      1,
			IntegrationID:           100,
			NotificationTypeID:      1,
			NotificationEventTypeID: 5,
			Enabled:                 true,
		},
	}

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			if filters.Enabled != nil && *filters.Enabled == enabledFilter {
				return activeConfigs, nil
			}
			return []entities.IntegrationNotificationConfig{}, nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	filters := dtos.FilterNotificationConfigDTO{
		Enabled: &enabledFilter,
	}

	// Act
	result, err := useCase.List(ctx, filters)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 config, got %d", len(result))
	}

	if !result[0].Enabled {
		t.Error("expected config to be enabled")
	}
}

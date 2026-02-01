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
			ID:               1,
			IntegrationID:    100,
			NotificationType: "whatsapp",
			IsActive:         true,
			Conditions: entities.NotificationConditions{
				Trigger: "order.created",
			},
			Config: entities.NotificationConfig{
				TemplateName: "template1",
			},
			Priority:  1,
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:               2,
			IntegrationID:    200,
			NotificationType: "email",
			IsActive:         false,
			Conditions: entities.NotificationConditions{
				Trigger: "order.updated",
			},
			Config: entities.NotificationConfig{
				TemplateName: "template2",
			},
			Priority:  2,
			CreatedAt: time.Now().Add(-36 * time.Hour),
			UpdatedAt: time.Now().Add(-12 * time.Hour),
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

	filters := dtos.FilterNotificationConfigDTO{}

	// Act
	result, err := useCase.List(ctx, filters)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != len(expectedConfigs) {
		t.Errorf("expected %d configs, got %d", len(expectedConfigs), len(result))
	}

	for i, config := range result {
		if config.ID != expectedConfigs[i].ID {
			t.Errorf("config[%d]: expected ID %d, got %d", i, expectedConfigs[i].ID, config.ID)
		}
		if config.IntegrationID != expectedConfigs[i].IntegrationID {
			t.Errorf("config[%d]: expected IntegrationID %d, got %d", i, expectedConfigs[i].IntegrationID, config.IntegrationID)
		}
		if config.NotificationType != expectedConfigs[i].NotificationType {
			t.Errorf("config[%d]: expected NotificationType %s, got %s", i, expectedConfigs[i].NotificationType, config.NotificationType)
		}
	}
}

func TestList_Success_WithFilters(t *testing.T) {
	// Arrange
	ctx := context.Background()

	integrationID := uint(100)
	notificationType := "whatsapp"
	isActive := true
	trigger := "order.created"

	expectedConfigs := []entities.IntegrationNotificationConfig{
		{
			ID:               1,
			IntegrationID:    100,
			NotificationType: "whatsapp",
			IsActive:         true,
			Conditions: entities.NotificationConditions{
				Trigger: "order.created",
			},
			Config: entities.NotificationConfig{
				TemplateName: "template1",
			},
			Priority: 1,
		},
	}

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			// Validar que los filtros se pasaron correctamente
			if filters.IntegrationID != nil && *filters.IntegrationID != integrationID {
				t.Errorf("expected IntegrationID filter %d, got %d", integrationID, *filters.IntegrationID)
			}
			if filters.NotificationType != nil && *filters.NotificationType != notificationType {
				t.Errorf("expected NotificationType filter %s, got %s", notificationType, *filters.NotificationType)
			}
			if filters.IsActive != nil && *filters.IsActive != isActive {
				t.Errorf("expected IsActive filter %v, got %v", isActive, *filters.IsActive)
			}
			if filters.Trigger != nil && *filters.Trigger != trigger {
				t.Errorf("expected Trigger filter %s, got %s", trigger, *filters.Trigger)
			}

			return expectedConfigs, nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	filters := dtos.FilterNotificationConfigDTO{
		IntegrationID:    &integrationID,
		NotificationType: &notificationType,
		IsActive:         &isActive,
		Trigger:          &trigger,
	}

	// Act
	result, err := useCase.List(ctx, filters)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected 1 config, got %d", len(result))
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

	filters := dtos.FilterNotificationConfigDTO{}

	// Act
	result, err := useCase.List(ctx, filters)

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
	expectedErr := errors.New("database connection failed")

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

	filters := dtos.FilterNotificationConfigDTO{}

	// Act
	result, err := useCase.List(ctx, filters)

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

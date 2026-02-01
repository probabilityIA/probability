package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/mocks"
)

func TestGetByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()

	expectedConfig := &entities.IntegrationNotificationConfig{
		ID:               1,
		IntegrationID:    100,
		NotificationType: "whatsapp",
		IsActive:         true,
		Conditions: entities.NotificationConditions{
			Trigger:  "order.created",
			Statuses: []string{"pending", "processing"},
		},
		Config: entities.NotificationConfig{
			TemplateName:  "confirmacion_pedido",
			RecipientType: "customer",
			Language:      "es",
		},
		Description: "Notificación de confirmación",
		Priority:    1,
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now(),
	}

	mockRepo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.IntegrationNotificationConfig, error) {
			if id == 1 {
				return expectedConfig, nil
			}
			return nil, errors.New("not found")
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	// Act
	result, err := useCase.GetByID(ctx, 1)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if result.ID != expectedConfig.ID {
		t.Errorf("expected ID %d, got %d", expectedConfig.ID, result.ID)
	}

	if result.IntegrationID != expectedConfig.IntegrationID {
		t.Errorf("expected IntegrationID %d, got %d", expectedConfig.IntegrationID, result.IntegrationID)
	}

	if result.NotificationType != expectedConfig.NotificationType {
		t.Errorf("expected NotificationType %s, got %s", expectedConfig.NotificationType, result.NotificationType)
	}

	if result.IsActive != expectedConfig.IsActive {
		t.Errorf("expected IsActive %v, got %v", expectedConfig.IsActive, result.IsActive)
	}

	if result.Description != expectedConfig.Description {
		t.Errorf("expected Description %s, got %s", expectedConfig.Description, result.Description)
	}

	if result.Priority != expectedConfig.Priority {
		t.Errorf("expected Priority %d, got %d", expectedConfig.Priority, result.Priority)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expectedErr := errors.New("record not found")

	mockRepo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.IntegrationNotificationConfig, error) {
			return nil, expectedErr
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	// Act
	result, err := useCase.GetByID(ctx, 999)

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

func TestGetByID_RepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expectedErr := errors.New("database connection failed")

	mockRepo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.IntegrationNotificationConfig, error) {
			return nil, expectedErr
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	// Act
	result, err := useCase.GetByID(ctx, 1)

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

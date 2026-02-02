package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/mocks"
)

func TestGetByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	configID := uint(1)

	expectedConfig := &entities.IntegrationNotificationConfig{
		ID:                      configID,
		IntegrationID:           100,
		NotificationTypeID:      1,
		NotificationEventTypeID: 5,
		Enabled:                 true,
		Description:             "Notificación de confirmación",
		OrderStatusIDs:          []uint{1, 2},
	}

	mockRepo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.IntegrationNotificationConfig, error) {
			if id == configID {
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
	result, err := useCase.GetByID(ctx, configID)

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

	if result.NotificationTypeID != expectedConfig.NotificationTypeID {
		t.Errorf("expected NotificationTypeID %d, got %d", expectedConfig.NotificationTypeID, result.NotificationTypeID)
	}

	if result.NotificationEventTypeID != expectedConfig.NotificationEventTypeID {
		t.Errorf("expected NotificationEventTypeID %d, got %d", expectedConfig.NotificationEventTypeID, result.NotificationEventTypeID)
	}

	if result.Enabled != expectedConfig.Enabled {
		t.Errorf("expected Enabled %v, got %v", expectedConfig.Enabled, result.Enabled)
	}

	if result.Description != expectedConfig.Description {
		t.Errorf("expected Description %s, got %s", expectedConfig.Description, result.Description)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	configID := uint(999)
	expectedErr := errors.New("config not found")

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
	result, err := useCase.GetByID(ctx, configID)

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
	configID := uint(1)
	expectedErr := errors.New("database query failed")

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
	result, err := useCase.GetByID(ctx, configID)

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

func TestGetByID_ZeroID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	configID := uint(0)

	mockRepo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.IntegrationNotificationConfig, error) {
			return nil, errors.New("invalid ID")
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	// Act
	result, err := useCase.GetByID(ctx, configID)

	// Assert
	if err == nil {
		t.Fatal("expected error for zero ID, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestGetByID_WithRelations(t *testing.T) {
	// Arrange
	ctx := context.Background()
	configID := uint(1)

	mockNotificationType := &entities.NotificationType{
		ID:          1,
		Code:        "whatsapp",
		Name:        "WhatsApp",
		Description: "Notificaciones por WhatsApp",
	}

	mockEventType := &entities.NotificationEventType{
		ID:                 5,
		NotificationTypeID: 1,
		EventCode:          "order.created",
		EventName:          "Pedido Creado",
		Description:        "Se dispara cuando se crea un pedido",
	}

	expectedConfig := &entities.IntegrationNotificationConfig{
		ID:                      configID,
		IntegrationID:           100,
		NotificationTypeID:      1,
		NotificationEventTypeID: 5,
		Enabled:                 true,
		Description:             "Notificación de confirmación",
		OrderStatusIDs:          []uint{1, 2},
		NotificationType:        mockNotificationType,
		NotificationEventType:   mockEventType,
	}

	mockRepo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.IntegrationNotificationConfig, error) {
			if id == configID {
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
	result, err := useCase.GetByID(ctx, configID)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	// Verificar que las relaciones estén presentes en el DTO
	if result.NotificationTypeName == nil {
		t.Error("expected NotificationTypeName to be present")
	} else {
		if *result.NotificationTypeName != "WhatsApp" {
			t.Errorf("expected NotificationTypeName 'WhatsApp', got '%s'", *result.NotificationTypeName)
		}
	}

	if result.NotificationEventName == nil {
		t.Error("expected NotificationEventName to be present")
	} else {
		if *result.NotificationEventName != "Pedido Creado" {
			t.Errorf("expected NotificationEventName 'Pedido Creado', got '%s'", *result.NotificationEventName)
		}
	}
}

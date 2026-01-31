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

func TestUpdate_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()

	existingConfig := &entities.IntegrationNotificationConfig{
		ID:               1,
		IntegrationID:    100,
		NotificationType: "whatsapp",
		IsActive:         true,
		Conditions: entities.NotificationConditions{
			Trigger:  "order.created",
			Statuses: []string{"pending"},
		},
		Config: entities.NotificationConfig{
			TemplateName:  "confirmacion_pedido",
			RecipientType: "customer",
			Language:      "es",
		},
		Description: "Configuración original",
		Priority:    1,
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now().Add(-24 * time.Hour),
	}

	updatedConfig := &entities.IntegrationNotificationConfig{
		ID:               1,
		IntegrationID:    100,
		NotificationType: "email", // Cambiado
		IsActive:         false,   // Cambiado
		Conditions: entities.NotificationConditions{
			Trigger:  "order.updated", // Cambiado
			Statuses: []string{"processing"},
		},
		Config: entities.NotificationConfig{
			TemplateName:  "nuevo_template",
			RecipientType: "business",
			Language:      "en",
		},
		Description: "Configuración actualizada",
		Priority:    2,
		CreatedAt:   existingConfig.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	mockRepo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.IntegrationNotificationConfig, error) {
			if id == 1 {
				// Primera llamada: retorna config original
				return existingConfig, nil
			}
			// Segunda llamada (después del update): retorna config actualizada
			return updatedConfig, nil
		},
		UpdateFn: func(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
			return nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockLogger)

	newNotificationType := "email"
	newIsActive := false
	newDescription := "Configuración actualizada"
	newPriority := 2
	newConditions := entities.NotificationConditions{
		Trigger:  "order.updated",
		Statuses: []string{"processing"},
	}
	newConfig := entities.NotificationConfig{
		TemplateName:  "nuevo_template",
		RecipientType: "business",
		Language:      "en",
	}

	dto := dtos.UpdateNotificationConfigDTO{
		NotificationType: &newNotificationType,
		IsActive:         &newIsActive,
		Conditions:       &newConditions,
		Config:           &newConfig,
		Description:      &newDescription,
		Priority:         &newPriority,
	}

	// Act
	result, err := useCase.Update(ctx, 1, dto)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if result.NotificationType != newNotificationType {
		t.Errorf("expected NotificationType %s, got %s", newNotificationType, result.NotificationType)
	}

	if result.IsActive != newIsActive {
		t.Errorf("expected IsActive %v, got %v", newIsActive, result.IsActive)
	}

	if result.Description != newDescription {
		t.Errorf("expected Description %s, got %s", newDescription, result.Description)
	}

	if result.Priority != newPriority {
		t.Errorf("expected Priority %d, got %d", newPriority, result.Priority)
	}
}

func TestUpdate_PartialUpdate(t *testing.T) {
	// Arrange
	ctx := context.Background()

	existingConfig := &entities.IntegrationNotificationConfig{
		ID:               1,
		IntegrationID:    100,
		NotificationType: "whatsapp",
		IsActive:         true,
		Conditions: entities.NotificationConditions{
			Trigger:  "order.created",
			Statuses: []string{"pending"},
		},
		Config: entities.NotificationConfig{
			TemplateName:  "confirmacion_pedido",
			RecipientType: "customer",
			Language:      "es",
		},
		Description: "Configuración original",
		Priority:    1,
	}

	mockRepo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.IntegrationNotificationConfig, error) {
			return existingConfig, nil
		},
		UpdateFn: func(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
			return nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockLogger)

	// Solo actualizar IsActive
	newIsActive := false
	dto := dtos.UpdateNotificationConfigDTO{
		IsActive: &newIsActive,
	}

	// Act
	result, err := useCase.Update(ctx, 1, dto)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if result.IsActive != newIsActive {
		t.Errorf("expected IsActive %v, got %v", newIsActive, result.IsActive)
	}

	// Verificar que otros campos no cambiaron
	if result.NotificationType != existingConfig.NotificationType {
		t.Errorf("NotificationType should not change, expected %s, got %s", existingConfig.NotificationType, result.NotificationType)
	}
}

func TestUpdate_ConfigNotFound(t *testing.T) {
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
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockLogger)

	newIsActive := false
	dto := dtos.UpdateNotificationConfigDTO{
		IsActive: &newIsActive,
	}

	// Act
	result, err := useCase.Update(ctx, 999, dto)

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

func TestUpdate_RepositoryUpdateError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expectedErr := errors.New("database update failed")

	existingConfig := &entities.IntegrationNotificationConfig{
		ID:               1,
		IntegrationID:    100,
		NotificationType: "whatsapp",
		IsActive:         true,
		Priority:         1,
	}

	mockRepo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.IntegrationNotificationConfig, error) {
			return existingConfig, nil
		},
		UpdateFn: func(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
			return expectedErr
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockLogger)

	newIsActive := false
	dto := dtos.UpdateNotificationConfigDTO{
		IsActive: &newIsActive,
	}

	// Act
	result, err := useCase.Update(ctx, 1, dto)

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

func TestUpdate_GetUpdatedConfigError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expectedErr := errors.New("failed to fetch updated config")

	existingConfig := &entities.IntegrationNotificationConfig{
		ID:               1,
		IntegrationID:    100,
		NotificationType: "whatsapp",
		IsActive:         true,
		Priority:         1,
	}

	callCount := 0
	mockRepo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.IntegrationNotificationConfig, error) {
			callCount++
			if callCount == 1 {
				// Primera llamada: retorna config existente
				return existingConfig, nil
			}
			// Segunda llamada (después del update): retorna error
			return nil, expectedErr
		},
		UpdateFn: func(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
			return nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockLogger)

	newIsActive := false
	dto := dtos.UpdateNotificationConfigDTO{
		IsActive: &newIsActive,
	}

	// Act
	result, err := useCase.Update(ctx, 1, dto)

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

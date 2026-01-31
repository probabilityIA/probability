package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	domainErrors "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/mocks"
)

func TestCreate_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			// No hay configuraciones existentes
			return []entities.IntegrationNotificationConfig{}, nil
		},
		CreateFn: func(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
			// Simular que la BD asigna el ID
			config.ID = 1
			config.CreatedAt = time.Now()
			config.UpdatedAt = time.Now()
			return nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockLogger)

	dto := dtos.CreateNotificationConfigDTO{
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
		Description: "Notificación de confirmación de pedido",
		Priority:    1,
	}

	// Act
	result, err := useCase.Create(ctx, dto)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if result.ID == 0 {
		t.Error("expected ID to be set")
	}

	if result.IntegrationID != dto.IntegrationID {
		t.Errorf("expected IntegrationID %d, got %d", dto.IntegrationID, result.IntegrationID)
	}

	if result.NotificationType != dto.NotificationType {
		t.Errorf("expected NotificationType %s, got %s", dto.NotificationType, result.NotificationType)
	}

	if result.IsActive != dto.IsActive {
		t.Errorf("expected IsActive %v, got %v", dto.IsActive, result.IsActive)
	}

	if result.Description != dto.Description {
		t.Errorf("expected Description %s, got %s", dto.Description, result.Description)
	}

	if result.Priority != dto.Priority {
		t.Errorf("expected Priority %d, got %d", dto.Priority, result.Priority)
	}
}

func TestCreate_DuplicateConfig(t *testing.T) {
	// Arrange
	ctx := context.Background()

	existingConfig := entities.IntegrationNotificationConfig{
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
		Description: "Configuración existente",
		Priority:    1,
	}

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			// Retornar configuración existente con condiciones similares
			return []entities.IntegrationNotificationConfig{existingConfig}, nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockLogger)

	dto := dtos.CreateNotificationConfigDTO{
		IntegrationID:    100,
		NotificationType: "whatsapp",
		IsActive:         true,
		Conditions: entities.NotificationConditions{
			Trigger:  "order.created",
			Statuses: []string{"pending"}, // Mismas condiciones
		},
		Config: entities.NotificationConfig{
			TemplateName:  "otro_template",
			RecipientType: "customer",
			Language:      "es",
		},
		Description: "Intento de duplicado",
		Priority:    2,
	}

	// Act
	result, err := useCase.Create(ctx, dto)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, domainErrors.ErrDuplicateConfig) {
		t.Errorf("expected ErrDuplicateConfig, got %v", err)
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestCreate_RepositoryListError(t *testing.T) {
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
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockLogger)

	dto := dtos.CreateNotificationConfigDTO{
		IntegrationID:    100,
		NotificationType: "whatsapp",
		IsActive:         true,
		Conditions: entities.NotificationConditions{
			Trigger: "order.created",
		},
		Config: entities.NotificationConfig{
			TemplateName: "template",
		},
		Priority: 1,
	}

	// Act
	result, err := useCase.Create(ctx, dto)

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

func TestCreate_RepositoryCreateError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expectedErr := errors.New("failed to insert into database")

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			return []entities.IntegrationNotificationConfig{}, nil
		},
		CreateFn: func(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
			return expectedErr
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockLogger)

	dto := dtos.CreateNotificationConfigDTO{
		IntegrationID:    100,
		NotificationType: "email",
		IsActive:         true,
		Conditions: entities.NotificationConditions{
			Trigger: "order.updated",
		},
		Config: entities.NotificationConfig{
			TemplateName: "update_template",
		},
		Priority: 2,
	}

	// Act
	result, err := useCase.Create(ctx, dto)

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

func TestCreate_DifferentConditions_ShouldSucceed(t *testing.T) {
	// Arrange
	ctx := context.Background()

	existingConfig := entities.IntegrationNotificationConfig{
		ID:               1,
		IntegrationID:    100,
		NotificationType: "whatsapp",
		Conditions: entities.NotificationConditions{
			Trigger:  "order.created",
			Statuses: []string{"pending"},
		},
	}

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			return []entities.IntegrationNotificationConfig{existingConfig}, nil
		},
		CreateFn: func(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
			config.ID = 2
			return nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockLogger)

	dto := dtos.CreateNotificationConfigDTO{
		IntegrationID:    100,
		NotificationType: "whatsapp",
		IsActive:         true,
		Conditions: entities.NotificationConditions{
			Trigger:  "order.created",
			Statuses: []string{"processing"}, // Diferente status
		},
		Config: entities.NotificationConfig{
			TemplateName: "processing_template",
		},
		Priority: 1,
	}

	// Act
	result, err := useCase.Create(ctx, dto)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if result.ID == 0 {
		t.Error("expected ID to be set")
	}
}

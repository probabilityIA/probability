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
	businessID := uint(10)

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
	mockCacheManager := &mocks.CacheManagerMock{
		CacheConfigFn: func(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
			return nil
		},
	}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, &mocks.MessageAuditQuerierMock{}, mockLogger)

	dto := dtos.CreateNotificationConfigDTO{
		BusinessID:              &businessID,
		IntegrationID:           100,
		NotificationTypeID:      1, // WhatsApp
		NotificationEventTypeID: 5, // order.created
		Enabled:                 true,
		Description:             "Notificación de confirmación de pedido",
		OrderStatusIDs:          []uint{1, 2}, // pending, processing
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

	if result.NotificationTypeID != dto.NotificationTypeID {
		t.Errorf("expected NotificationTypeID %d, got %d", dto.NotificationTypeID, result.NotificationTypeID)
	}

	if result.NotificationEventTypeID != dto.NotificationEventTypeID {
		t.Errorf("expected NotificationEventTypeID %d, got %d", dto.NotificationEventTypeID, result.NotificationEventTypeID)
	}

	if result.Enabled != dto.Enabled {
		t.Errorf("expected Enabled %v, got %v", dto.Enabled, result.Enabled)
	}

	if result.Description != dto.Description {
		t.Errorf("expected Description %s, got %s", dto.Description, result.Description)
	}
}

func TestCreate_DuplicateConfig(t *testing.T) {
	// Arrange
	ctx := context.Background()

	existingConfig := entities.IntegrationNotificationConfig{
		ID:                      1,
		IntegrationID:           100,
		NotificationTypeID:      1, // WhatsApp
		NotificationEventTypeID: 5, // order.created
		Enabled:                 true,
		Description:             "Configuración existente",
		OrderStatusIDs:          []uint{1},
	}

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			// Retornar configuración existente con la misma combinación
			return []entities.IntegrationNotificationConfig{existingConfig}, nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, &mocks.MessageAuditQuerierMock{}, mockLogger)

	dto := dtos.CreateNotificationConfigDTO{
		IntegrationID:           100,
		NotificationTypeID:      1, // Mismo tipo
		NotificationEventTypeID: 5, // Mismo evento
		Enabled:                 true,
		Description:             "Intento de duplicado",
		OrderStatusIDs:          []uint{2}, // Diferente estado, pero aún es duplicado
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
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, &mocks.MessageAuditQuerierMock{}, mockLogger)

	dto := dtos.CreateNotificationConfigDTO{
		IntegrationID:           100,
		NotificationTypeID:      1,
		NotificationEventTypeID: 5,
		Enabled:                 true,
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
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, &mocks.MessageAuditQuerierMock{}, mockLogger)

	dto := dtos.CreateNotificationConfigDTO{
		IntegrationID:           100,
		NotificationTypeID:      2, // Email
		NotificationEventTypeID: 6, // order.updated
		Enabled:                 true,
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

func TestCreate_CacheErrorShouldNotFail(t *testing.T) {
	// Arrange
	ctx := context.Background()
	cacheErr := errors.New("redis connection failed")

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			return []entities.IntegrationNotificationConfig{}, nil
		},
		CreateFn: func(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
			config.ID = 1
			return nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{
		CacheConfigFn: func(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
			return cacheErr // Error en cache
		},
	}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, &mocks.MessageAuditQuerierMock{}, mockLogger)

	dto := dtos.CreateNotificationConfigDTO{
		IntegrationID:           100,
		NotificationTypeID:      1,
		NotificationEventTypeID: 5,
		Enabled:                 true,
	}

	// Act
	result, err := useCase.Create(ctx, dto)

	// Assert - No debe fallar si el cache falla (cache es secundario)
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

func TestCreate_DifferentIntegration_ShouldSucceed(t *testing.T) {
	// Arrange
	ctx := context.Background()

	existingConfig := entities.IntegrationNotificationConfig{
		ID:                      1,
		IntegrationID:           100,
		NotificationTypeID:      1,
		NotificationEventTypeID: 5,
		Enabled:                 true,
	}

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			// Filtro por integración diferente, no retorna nada
			if filters.IntegrationID != nil && *filters.IntegrationID == 200 {
				return []entities.IntegrationNotificationConfig{}, nil
			}
			return []entities.IntegrationNotificationConfig{existingConfig}, nil
		},
		CreateFn: func(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
			config.ID = 2
			return nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, &mocks.MessageAuditQuerierMock{}, mockLogger)

	dto := dtos.CreateNotificationConfigDTO{
		IntegrationID:           200, // Diferente integración
		NotificationTypeID:      1,
		NotificationEventTypeID: 5,
		Enabled:                 true,
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

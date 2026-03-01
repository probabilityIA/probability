package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/mocks"
)

func TestUpdate_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	configID := uint(1)
	newDescription := "Nueva descripción actualizada"
	newEnabled := false

	existingConfig := &entities.IntegrationNotificationConfig{
		ID:                      configID,
		IntegrationID:           100,
		NotificationTypeID:      1,
		NotificationEventTypeID: 5,
		Enabled:                 true,
		Description:             "Descripción original",
		OrderStatusIDs:          []uint{1, 2},
	}

	updatedConfig := &entities.IntegrationNotificationConfig{
		ID:                      configID,
		IntegrationID:           100,
		NotificationTypeID:      1,
		NotificationEventTypeID: 5,
		Enabled:                 newEnabled,
		Description:             newDescription,
		OrderStatusIDs:          []uint{1, 2},
	}

	mockRepo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.IntegrationNotificationConfig, error) {
			if id == configID {
				// Primera llamada retorna config original, segunda llamada retorna actualizada
				return existingConfig, nil
			}
			return nil, errors.New("not found")
		},
		UpdateFn: func(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
			// Simular actualización exitosa
			existingConfig = updatedConfig
			return nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{
		UpdateConfigInCacheFn: func(ctx context.Context, oldConfig, newConfig *entities.IntegrationNotificationConfig) error {
			return nil
		},
	}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, &mocks.MessageAuditQuerierMock{}, mockLogger)

	dto := dtos.UpdateNotificationConfigDTO{
		Description: &newDescription,
		IsActive:    &newEnabled,
	}

	// Act
	result, err := useCase.Update(ctx, configID, dto)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if result.ID != configID {
		t.Errorf("expected ID %d, got %d", configID, result.ID)
	}
}

func TestUpdate_OnlyDescription(t *testing.T) {
	// Arrange
	ctx := context.Background()
	configID := uint(1)
	newDescription := "Solo actualizar descripción"

	existingConfig := &entities.IntegrationNotificationConfig{
		ID:                      configID,
		IntegrationID:           100,
		NotificationTypeID:      1,
		NotificationEventTypeID: 5,
		Enabled:                 true,
		Description:             "Descripción original",
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
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, &mocks.MessageAuditQuerierMock{}, mockLogger)

	dto := dtos.UpdateNotificationConfigDTO{
		Description: &newDescription,
		IsActive:    nil, // No actualizar enabled
	}

	// Act
	result, err := useCase.Update(ctx, configID, dto)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestUpdate_OnlyEnabled(t *testing.T) {
	// Arrange
	ctx := context.Background()
	configID := uint(1)
	newEnabled := false

	existingConfig := &entities.IntegrationNotificationConfig{
		ID:                      configID,
		IntegrationID:           100,
		NotificationTypeID:      1,
		NotificationEventTypeID: 5,
		Enabled:                 true,
		Description:             "Descripción",
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
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, &mocks.MessageAuditQuerierMock{}, mockLogger)

	dto := dtos.UpdateNotificationConfigDTO{
		Description: nil, // No actualizar descripción
		IsActive:    &newEnabled,
	}

	// Act
	result, err := useCase.Update(ctx, configID, dto)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestUpdate_NotFound(t *testing.T) {
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

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, &mocks.MessageAuditQuerierMock{}, mockLogger)

	newDescription := "Nueva descripción"
	dto := dtos.UpdateNotificationConfigDTO{
		Description: &newDescription,
	}

	// Act
	result, err := useCase.Update(ctx, configID, dto)

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
	configID := uint(1)
	expectedErr := errors.New("database update failed")

	existingConfig := &entities.IntegrationNotificationConfig{
		ID:                      configID,
		IntegrationID:           100,
		NotificationTypeID:      1,
		NotificationEventTypeID: 5,
		Enabled:                 true,
		Description:             "Original",
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
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, &mocks.MessageAuditQuerierMock{}, mockLogger)

	newDescription := "Nueva descripción"
	dto := dtos.UpdateNotificationConfigDTO{
		Description: &newDescription,
	}

	// Act
	result, err := useCase.Update(ctx, configID, dto)

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

func TestUpdate_CacheErrorShouldNotFail(t *testing.T) {
	// Arrange
	ctx := context.Background()
	configID := uint(1)
	cacheErr := errors.New("redis connection failed")

	existingConfig := &entities.IntegrationNotificationConfig{
		ID:                      configID,
		IntegrationID:           100,
		NotificationTypeID:      1,
		NotificationEventTypeID: 5,
		Enabled:                 true,
		Description:             "Original",
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
	mockCacheManager := &mocks.CacheManagerMock{
		UpdateConfigInCacheFn: func(ctx context.Context, oldConfig, newConfig *entities.IntegrationNotificationConfig) error {
			return cacheErr // Error en cache
		},
	}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, &mocks.MessageAuditQuerierMock{}, mockLogger)

	newDescription := "Nueva descripción"
	dto := dtos.UpdateNotificationConfigDTO{
		Description: &newDescription,
	}

	// Act
	result, err := useCase.Update(ctx, configID, dto)

	// Assert - No debe fallar si el cache falla (cache es secundario)
	if err != nil {
		t.Fatalf("expected no error (cache error should be logged only), got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestUpdate_EmptyDTO(t *testing.T) {
	// Arrange
	ctx := context.Background()
	configID := uint(1)

	existingConfig := &entities.IntegrationNotificationConfig{
		ID:                      configID,
		IntegrationID:           100,
		NotificationTypeID:      1,
		NotificationEventTypeID: 5,
		Enabled:                 true,
		Description:             "Original",
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
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, &mocks.MessageAuditQuerierMock{}, mockLogger)

	dto := dtos.UpdateNotificationConfigDTO{
		Description: nil, // No actualizar nada
		IsActive:    nil,
	}

	// Act
	result, err := useCase.Update(ctx, configID, dto)

	// Assert - Debe permitir updates vacíos (no hace cambios)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}
}

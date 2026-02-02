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

func TestDeleteNotificationEventType_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	eventTypeID := uint(10)

	mockEventType := &entities.NotificationEventType{
		ID:                 eventTypeID,
		NotificationTypeID: 1,
		EventCode:          "order.created",
		EventName:          "Pedido Creado",
		Description:        "Se dispara cuando se crea un nuevo pedido",
	}

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			// No hay configuraciones usando este evento
			return []entities.IntegrationNotificationConfig{}, nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.NotificationEventType, error) {
			return mockEventType, nil
		},
		DeleteFn: func(ctx context.Context, id uint) error {
			return nil
		},
	}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	// Act
	err := useCase.DeleteNotificationEventType(ctx, eventTypeID)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteNotificationEventType_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	eventTypeID := uint(999)

	mockRepo := &mocks.RepositoryMock{}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.NotificationEventType, error) {
			return nil, errors.New("not found")
		},
	}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	// Act
	err := useCase.DeleteNotificationEventType(ctx, eventTypeID)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, domainErrors.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteNotificationEventType_HasActiveConfigs(t *testing.T) {
	// Arrange
	ctx := context.Background()
	eventTypeID := uint(10)

	mockEventType := &entities.NotificationEventType{
		ID:                 eventTypeID,
		NotificationTypeID: 1,
		EventCode:          "order.created",
		EventName:          "Pedido Creado",
		Description:        "Se dispara cuando se crea un nuevo pedido",
	}

	mockNotificationType := &entities.NotificationType{
		ID:   1,
		Code: "whatsapp",
		Name: "WhatsApp",
	}

	// Configuraciones activas usando este evento
	activeConfigs := []entities.IntegrationNotificationConfig{
		{
			ID:                      1,
			IntegrationID:           100,
			NotificationTypeID:      1,
			NotificationEventTypeID: eventTypeID,
			Enabled:                 true,
			Description:             "Config WhatsApp pedido creado",
			NotificationType:        mockNotificationType,
		},
		{
			ID:                      2,
			IntegrationID:           101,
			NotificationTypeID:      1,
			NotificationEventTypeID: eventTypeID,
			Enabled:                 true,
			Description:             "Otra config activa",
			NotificationType:        mockNotificationType,
		},
	}

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			// Retornar configuraciones activas
			return activeConfigs, nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.NotificationEventType, error) {
			return mockEventType, nil
		},
	}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	// Act
	err := useCase.DeleteNotificationEventType(ctx, eventTypeID)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Verificar que el mensaje contenga información sobre las configuraciones activas
	errMsg := err.Error()
	if errMsg == "" {
		t.Error("expected error message, got empty string")
	}

	// Verificar que mencione las 2 configuraciones activas
	expectedSubstrings := []string{
		"2 configuración(es) activa(s)",
		"no se puede eliminar",
	}

	for _, substr := range expectedSubstrings {
		if !contains(errMsg, substr) {
			t.Errorf("error message should contain '%s', got: %s", substr, errMsg)
		}
	}
}

func TestDeleteNotificationEventType_HasInactiveConfigs_ShouldSucceed(t *testing.T) {
	// Arrange
	ctx := context.Background()
	eventTypeID := uint(10)

	mockEventType := &entities.NotificationEventType{
		ID:                 eventTypeID,
		NotificationTypeID: 1,
		EventCode:          "order.created",
		EventName:          "Pedido Creado",
	}

	// Configuraciones INACTIVAS usando este evento
	inactiveConfigs := []entities.IntegrationNotificationConfig{
		{
			ID:                      1,
			IntegrationID:           100,
			NotificationTypeID:      1,
			NotificationEventTypeID: eventTypeID,
			Enabled:                 false, // Inactiva
			Description:             "Config inactiva",
		},
		{
			ID:                      2,
			IntegrationID:           101,
			NotificationTypeID:      1,
			NotificationEventTypeID: eventTypeID,
			Enabled:                 false, // Inactiva
		},
	}

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			return inactiveConfigs, nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.NotificationEventType, error) {
			return mockEventType, nil
		},
		DeleteFn: func(ctx context.Context, id uint) error {
			return nil
		},
	}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	// Act
	err := useCase.DeleteNotificationEventType(ctx, eventTypeID)

	// Assert - Debe permitir eliminar si solo hay configs inactivas
	if err != nil {
		t.Fatalf("expected no error when only inactive configs exist, got %v", err)
	}
}

func TestDeleteNotificationEventType_MixedActiveInactive_ShouldFail(t *testing.T) {
	// Arrange
	ctx := context.Background()
	eventTypeID := uint(10)

	mockEventType := &entities.NotificationEventType{
		ID:                 eventTypeID,
		NotificationTypeID: 1,
		EventCode:          "order.created",
		EventName:          "Pedido Creado",
	}

	mockNotificationType := &entities.NotificationType{
		ID:   1,
		Code: "whatsapp",
		Name: "WhatsApp",
	}

	// Mezcla de configuraciones activas e inactivas
	mixedConfigs := []entities.IntegrationNotificationConfig{
		{
			ID:                      1,
			IntegrationID:           100,
			NotificationTypeID:      1,
			NotificationEventTypeID: eventTypeID,
			Enabled:                 true, // Activa
			Description:             "Config activa",
			NotificationType:        mockNotificationType,
		},
		{
			ID:                      2,
			IntegrationID:           101,
			NotificationTypeID:      1,
			NotificationEventTypeID: eventTypeID,
			Enabled:                 false, // Inactiva
		},
	}

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			return mixedConfigs, nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.NotificationEventType, error) {
			return mockEventType, nil
		},
	}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	// Act
	err := useCase.DeleteNotificationEventType(ctx, eventTypeID)

	// Assert - Debe fallar si al menos UNA config está activa
	if err == nil {
		t.Fatal("expected error when at least one active config exists, got nil")
	}

	errMsg := err.Error()
	if !contains(errMsg, "1 configuración(es) activa(s)") {
		t.Errorf("error message should mention 1 active config, got: %s", errMsg)
	}
}

func TestDeleteNotificationEventType_RepositoryListError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	eventTypeID := uint(10)
	expectedErr := errors.New("database query failed")

	mockEventType := &entities.NotificationEventType{
		ID:                 eventTypeID,
		NotificationTypeID: 1,
		EventCode:          "order.created",
	}

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			return nil, expectedErr
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.NotificationEventType, error) {
			return mockEventType, nil
		},
		DeleteFn: func(ctx context.Context, id uint) error {
			return nil
		},
	}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	// Act - Si hay error en List, aún así debe permitir eliminar
	err := useCase.DeleteNotificationEventType(ctx, eventTypeID)

	// Assert - El código actual permite eliminar si hay error en List
	if err != nil {
		t.Fatalf("expected no error (ignores List error), got %v", err)
	}
}

func TestDeleteNotificationEventType_DeleteError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	eventTypeID := uint(10)
	expectedErr := errors.New("database delete failed")

	mockEventType := &entities.NotificationEventType{
		ID:                 eventTypeID,
		NotificationTypeID: 1,
		EventCode:          "order.created",
	}

	mockRepo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error) {
			return []entities.IntegrationNotificationConfig{}, nil
		},
	}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.NotificationEventType, error) {
			return mockEventType, nil
		},
		DeleteFn: func(ctx context.Context, id uint) error {
			return expectedErr
		},
	}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	// Act
	err := useCase.DeleteNotificationEventType(ctx, eventTypeID)

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

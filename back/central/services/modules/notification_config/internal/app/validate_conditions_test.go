package app

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/mocks"
)

func TestValidateConditions_Success(t *testing.T) {
	tests := []struct {
		name            string
		config          *entities.IntegrationNotificationConfig
		orderStatusID   uint
		paymentMethodID uint
		expectedResult  bool
	}{
		{
			name: "Sin filtros de estados - debe aceptar cualquier estado",
			config: &entities.IntegrationNotificationConfig{
				ID:                      1,
				IntegrationID:           100,
				NotificationTypeID:      1,
				NotificationEventTypeID: 5,
				Enabled:                 true,
				OrderStatusIDs:          []uint{}, // Sin filtros
			},
			orderStatusID:   1,
			paymentMethodID: 1,
			expectedResult:  true,
		},
		{
			name: "Estado permitido en la lista - debe aceptar",
			config: &entities.IntegrationNotificationConfig{
				ID:                      1,
				IntegrationID:           100,
				NotificationTypeID:      1,
				NotificationEventTypeID: 5,
				Enabled:                 true,
				OrderStatusIDs:          []uint{1, 2, 3}, // pending, processing, completed
			},
			orderStatusID:   2, // processing
			paymentMethodID: 1,
			expectedResult:  true,
		},
		{
			name: "Estado NO permitido - debe rechazar",
			config: &entities.IntegrationNotificationConfig{
				ID:                      1,
				IntegrationID:           100,
				NotificationTypeID:      1,
				NotificationEventTypeID: 5,
				Enabled:                 true,
				OrderStatusIDs:          []uint{1, 2}, // pending, processing
			},
			orderStatusID:   5, // cancelled
			paymentMethodID: 1,
			expectedResult:  false,
		},
		{
			name: "Primer estado en la lista - debe aceptar",
			config: &entities.IntegrationNotificationConfig{
				ID:                      1,
				IntegrationID:           100,
				NotificationTypeID:      1,
				NotificationEventTypeID: 5,
				Enabled:                 true,
				OrderStatusIDs:          []uint{1, 2, 3},
			},
			orderStatusID:   1, // Primer elemento
			paymentMethodID: 1,
			expectedResult:  true,
		},
		{
			name: "Último estado en la lista - debe aceptar",
			config: &entities.IntegrationNotificationConfig{
				ID:                      1,
				IntegrationID:           100,
				NotificationTypeID:      1,
				NotificationEventTypeID: 5,
				Enabled:                 true,
				OrderStatusIDs:          []uint{1, 2, 3},
			},
			orderStatusID:   3, // Último elemento
			paymentMethodID: 1,
			expectedResult:  true,
		},
		{
			name: "Un solo estado permitido - debe aceptar si coincide",
			config: &entities.IntegrationNotificationConfig{
				ID:                      1,
				IntegrationID:           100,
				NotificationTypeID:      1,
				NotificationEventTypeID: 5,
				Enabled:                 true,
				OrderStatusIDs:          []uint{2}, // Solo processing
			},
			orderStatusID:   2,
			paymentMethodID: 1,
			expectedResult:  true,
		},
		{
			name: "Un solo estado permitido - debe rechazar si no coincide",
			config: &entities.IntegrationNotificationConfig{
				ID:                      1,
				IntegrationID:           100,
				NotificationTypeID:      1,
				NotificationEventTypeID: 5,
				Enabled:                 true,
				OrderStatusIDs:          []uint{2}, // Solo processing
			},
			orderStatusID:   1, // pending
			paymentMethodID: 1,
			expectedResult:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &mocks.RepositoryMock{}
			mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
			mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
			mockCacheManager := &mocks.CacheManagerMock{}
			mockLogger := mocks.NewLoggerMock()

			useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

			// Act
			result := useCase.ValidateConditions(tt.config, tt.orderStatusID, tt.paymentMethodID)

			// Assert
			if result != tt.expectedResult {
				t.Errorf("expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestValidateConditions_PaymentMethod(t *testing.T) {
	// Arrange
	mockRepo := &mocks.RepositoryMock{}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	config := &entities.IntegrationNotificationConfig{
		ID:                      1,
		IntegrationID:           100,
		NotificationTypeID:      1,
		NotificationEventTypeID: 5,
		Enabled:                 true,
		OrderStatusIDs:          []uint{1}, // pending
	}

	// Act - PaymentMethod actualmente no se valida (TODO en el código)
	result := useCase.ValidateConditions(config, 1, 999)

	// Assert - Debe aceptar cualquier paymentMethodID por ahora
	if !result {
		t.Error("expected true for any paymentMethodID (not yet implemented)")
	}
}

func TestValidateConditions_EmptyOrderStatusIDs(t *testing.T) {
	// Arrange
	mockRepo := &mocks.RepositoryMock{}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	config := &entities.IntegrationNotificationConfig{
		ID:                      1,
		IntegrationID:           100,
		NotificationTypeID:      1,
		NotificationEventTypeID: 5,
		Enabled:                 true,
		OrderStatusIDs:          []uint{}, // Sin filtros
	}

	// Act - Probar con diferentes order status IDs
	tests := []uint{0, 1, 10, 100, 999}
	for _, statusID := range tests {
		result := useCase.ValidateConditions(config, statusID, 1)

		// Assert - Debe aceptar TODOS los estados
		if !result {
			t.Errorf("expected true for orderStatusID %d when OrderStatusIDs is empty", statusID)
		}
	}
}

func TestValidateConditions_NilConfig(t *testing.T) {
	// Arrange
	mockRepo := &mocks.RepositoryMock{}
	mockNotificationTypeRepo := &mocks.NotificationTypeRepositoryMock{}
	mockEventTypeRepo := &mocks.NotificationEventTypeRepositoryMock{}
	mockCacheManager := &mocks.CacheManagerMock{}
	mockLogger := mocks.NewLoggerMock()

	useCase := New(mockRepo, mockNotificationTypeRepo, mockEventTypeRepo, mockCacheManager, mockLogger)

	// Act - Con config nil debería hacer panic (comportamiento actual)
	// Este test documenta el comportamiento actual: no maneja nil
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic with nil config, but didn't panic")
		}
	}()

	// Este test verifica que SÍ hay panic con nil (comportamiento actual)
	// En producción, el caller debe asegurar que config no sea nil
	_ = useCase.ValidateConditions(nil, 1, 1)
}

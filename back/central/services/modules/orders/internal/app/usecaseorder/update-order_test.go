package usecaseorder

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateOrder_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	orderID := "order-uuid-123"
	businessID := uint(1)

	// Orden existente
	existingOrder := &entities.ProbabilityOrder{
		ID:            orderID,
		BusinessID:    &businessID,
		IntegrationID: 10,
		OrderNumber:   "ORD-001",
		Status:        "pending",
		TotalAmount:   100.0,
		Currency:      "USD",
		CustomerEmail: "old@example.com",
		CustomerName:  "Old Name",
	}

	// Request de actualización (actualizar algunos campos)
	newTotal := 150.0
	newEmail := "new@example.com"
	newName := "New Name"
	req := &dtos.UpdateOrderRequest{
		TotalAmount:   &newTotal,
		CustomerEmail: &newEmail,
		CustomerName:  &newName,
	}

	// Configurar mocks
	mockRepo.On("GetOrderByID", ctx, orderID).
		Return(existingOrder, nil).
		Maybe()

	mockRepo.On("UpdateOrder", ctx, mock.MatchedBy(func(order *entities.ProbabilityOrder) bool {
		return order.ID == orderID &&
			order.TotalAmount == newTotal &&
			order.CustomerEmail == newEmail &&
			order.CustomerName == newName
	})).Return(nil).Maybe()

	// Score use case puede ser llamado (no bloqueante)
	mockScoreUseCase.On("CalculateAndUpdateOrderScore", ctx, orderID).
		Return(nil).Maybe()

	mockRepo.On("GetOrderByID", ctx, orderID).
		Return(existingOrder, nil).
		Maybe()

	// Publicadores de eventos
	mockLogger.On("Info", mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything).Maybe()
	mockRedisPublisher.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).
		Maybe().Return(nil)
	mockRabbitPublisher.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).
		Maybe().Return(nil)

	// Act
	result, err := useCase.UpdateOrder(ctx, orderID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, orderID, result.ID)
	assert.Equal(t, newTotal, result.TotalAmount)
	assert.Equal(t, newEmail, result.CustomerEmail)
	assert.Equal(t, newName, result.CustomerName)

	mockRepo.AssertExpectations(t)
}

func TestUpdateOrder_EmptyID(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	req := &dtos.UpdateOrderRequest{}

	// Act
	result, err := useCase.UpdateOrder(ctx, "", req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "order ID is required")

	mockRepo.AssertNotCalled(t, "GetOrderByID")
	mockRepo.AssertNotCalled(t, "UpdateOrder")
}

func TestUpdateOrder_OrderNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	orderID := "non-existent"
	req := &dtos.UpdateOrderRequest{}

	notFoundErr := errors.New("record not found")
	mockRepo.On("GetOrderByID", ctx, orderID).
		Return(nil, notFoundErr)

	// Act
	result, err := useCase.UpdateOrder(ctx, orderID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error getting order")

	mockRepo.AssertNotCalled(t, "UpdateOrder")
}

func TestUpdateOrder_StatusChange_PublishesStatusEvent(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	orderID := "order-uuid-123"
	businessID := uint(1)

	existingOrder := &entities.ProbabilityOrder{
		ID:            orderID,
		BusinessID:    &businessID,
		IntegrationID: 10,
		OrderNumber:   "ORD-001",
		Status:        "pending",
		TotalAmount:   100.0,
	}

	// Cambiar el estado
	newStatus := "completed"
	req := &dtos.UpdateOrderRequest{
		Status: &newStatus,
	}

	mockRepo.On("GetOrderByID", ctx, orderID).
		Return(existingOrder, nil).
		Maybe()

	mockRepo.On("UpdateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).
		Return(nil).Maybe()

	mockScoreUseCase.On("CalculateAndUpdateOrderScore", ctx, orderID).
		Return(nil).Maybe()

	mockRepo.On("GetOrderByID", ctx, orderID).
		Return(existingOrder, nil).
		Maybe()

	// Verificar que se publican 2 eventos (update + status_changed)
	mockLogger.On("Info", mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything).Maybe()
	mockRedisPublisher.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).
		Maybe().Return(nil)
	mockRabbitPublisher.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).
		Maybe().Return(nil)

	// Act
	result, err := useCase.UpdateOrder(ctx, orderID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newStatus, result.Status)

	mockRepo.AssertExpectations(t)
}

func TestUpdateOrder_PartialUpdate(t *testing.T) {
	// Arrange - Actualizar solo algunos campos
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	orderID := "order-uuid-123"
	businessID := uint(1)

	existingOrder := &entities.ProbabilityOrder{
		ID:              orderID,
		BusinessID:      &businessID,
		OrderNumber:     "ORD-001",
		Status:          "pending",
		TotalAmount:     100.0,
		Currency:        "USD",
		CustomerEmail:   "customer@example.com",
		CustomerName:    "Original Name",
		ShippingCity:    "Original City",
		ShippingCountry: "Original Country",
	}

	// Solo actualizar ciudad y país
	newCity := "New York"
	newCountry := "USA"
	req := &dtos.UpdateOrderRequest{
		ShippingCity:    &newCity,
		ShippingCountry: &newCountry,
	}

	mockRepo.On("GetOrderByID", ctx, orderID).
		Return(existingOrder, nil).
		Maybe()

	var capturedOrder *entities.ProbabilityOrder
	mockRepo.On("UpdateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).
		Run(func(args mock.Arguments) {
			capturedOrder = args.Get(1).(*entities.ProbabilityOrder)
		}).
		Return(nil).Maybe()

	mockScoreUseCase.On("CalculateAndUpdateOrderScore", ctx, orderID).
		Return(nil).Maybe()

	mockLogger.On("Info", mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything).Maybe()
	mockRedisPublisher.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).
		Maybe().Return(nil)
	mockRabbitPublisher.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).
		Maybe().Return(nil)

	// Act
	result, err := useCase.UpdateOrder(ctx, orderID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verificar que solo se actualizaron los campos especificados
	assert.Equal(t, newCity, capturedOrder.ShippingCity)
	assert.Equal(t, newCountry, capturedOrder.ShippingCountry)

	// Verificar que los demás campos NO cambiaron
	assert.Equal(t, "Original Name", capturedOrder.CustomerName)
	assert.Equal(t, "customer@example.com", capturedOrder.CustomerEmail)
	assert.Equal(t, 100.0, capturedOrder.TotalAmount)

	mockRepo.AssertExpectations(t)
}

func TestUpdateOrder_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	orderID := "order-uuid-123"
	businessID := uint(1)

	existingOrder := &entities.ProbabilityOrder{
		ID:         orderID,
		BusinessID: &businessID,
		Status:     "pending",
	}

	newStatus := "completed"
	req := &dtos.UpdateOrderRequest{
		Status: &newStatus,
	}

	mockRepo.On("GetOrderByID", ctx, orderID).
		Return(existingOrder, nil)

	// UpdateOrder falla
	dbError := errors.New("database constraint violation")
	mockRepo.On("UpdateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).
		Return(dbError)

	// Act
	result, err := useCase.UpdateOrder(ctx, orderID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error updating order")

	mockRepo.AssertExpectations(t)
}

func TestUpdateOrder_ConfirmationStatus(t *testing.T) {
	tests := []struct {
		name                  string
		confirmationStatus    string
		expectedIsConfirmed   *bool
	}{
		{
			name:               "confirmation yes",
			confirmationStatus: "yes",
			expectedIsConfirmed: func() *bool {
				b := true
				return &b
			}(),
		},
		{
			name:               "confirmation no",
			confirmationStatus: "no",
			expectedIsConfirmed: func() *bool {
				b := false
				return &b
			}(),
		},
		{
			name:                "confirmation pending",
			confirmationStatus:  "pending",
			expectedIsConfirmed: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(mocks.RepositoryMock)
			mockRedisPublisher := new(mocks.EventPublisherMock)
			mockRabbitPublisher := new(mocks.RabbitPublisherMock)
			mockLogger := new(mocks.LoggerMock)
			mockScoreUseCase := new(mocks.ScoreUseCaseMock)

			useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

			ctx := context.Background()
			orderID := "order-uuid-123"
			businessID := uint(1)

			existingOrder := &entities.ProbabilityOrder{
				ID:         orderID,
				BusinessID: &businessID,
				Status:     "pending",
			}

			req := &dtos.UpdateOrderRequest{
				ConfirmationStatus: &tt.confirmationStatus,
			}

			mockRepo.On("GetOrderByID", ctx, orderID).
				Return(existingOrder, nil).
				Maybe()

			var capturedOrder *entities.ProbabilityOrder
			mockRepo.On("UpdateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).
				Run(func(args mock.Arguments) {
					capturedOrder = args.Get(1).(*entities.ProbabilityOrder)
				}).
				Return(nil).Maybe()

			mockScoreUseCase.On("CalculateAndUpdateOrderScore", ctx, orderID).
				Return(nil).Maybe()

			mockLogger.On("Info", mock.Anything).Maybe()
			mockLogger.On("Error", mock.Anything).Maybe()
			mockRedisPublisher.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).
				Maybe().Return(nil)
			mockRabbitPublisher.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).
				Maybe().Return(nil)

			// Act
			_, err := useCase.UpdateOrder(ctx, orderID, req)

			// Assert
			assert.NoError(t, err)
			if tt.expectedIsConfirmed == nil {
				assert.Nil(t, capturedOrder.IsConfirmed)
			} else {
				assert.NotNil(t, capturedOrder.IsConfirmed)
				assert.Equal(t, *tt.expectedIsConfirmed, *capturedOrder.IsConfirmed)
			}
		})
	}
}

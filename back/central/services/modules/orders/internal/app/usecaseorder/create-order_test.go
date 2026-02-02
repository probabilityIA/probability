package usecaseorder

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrder_Success(t *testing.T) {
	// Arrange - Configurar mocks
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	// Configurar el caso de uso
	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	businessID := uint(1)
	integrationID := uint(10)

	// Request de creación
	req := &dtos.CreateOrderRequest{
		BusinessID:      &businessID,
		IntegrationID:   integrationID,
		IntegrationType: "shopify",
		Platform:        "Shopify",
		ExternalID:      "EXT-12345",
		OrderNumber:     "ORD-001",
		InternalNumber:  "INT-001",
		Subtotal:        100.0,
		Tax:             10.0,
		Discount:        0.0,
		ShippingCost:    5.0,
		TotalAmount:     115.0,
		Currency:        "USD",
		CustomerName:    "John Doe",
		CustomerEmail:   "john@example.com",
		CustomerPhone:   "+1234567890",
		Status:          "pending",
		OriginalStatus:  "pending",
		PaymentMethodID: 1,
		IsPaid:          false,
		OccurredAt:      time.Now(),
	}

	// Configurar expectativas de los mocks

	// 1. OrderExists retorna false (la orden no existe)
	mockRepo.On("OrderExists", ctx, req.ExternalID, req.IntegrationID).
		Return(false, nil)

	// 2. CreateOrder debe ser llamado y retorna nil (éxito)
	mockRepo.On("CreateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).
		Run(func(args mock.Arguments) {
			// Simular que la BD asigna un ID
			order := args.Get(1).(*entities.ProbabilityOrder)
			order.ID = "generated-uuid-123"
			order.CreatedAt = time.Now()
			order.UpdatedAt = time.Now()
		}).
		Return(nil)

	// 3. Publicar eventos (los mocks no hacen nada, solo registran la llamada)
	mockLogger.On("Info", mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything).Maybe()

	// Los publicadores pueden ser llamados en goroutines, usamos Maybe
	mockRedisPublisher.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).
		Maybe().
		Return(nil)
	mockRabbitPublisher.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).
		Maybe().
		Return(nil)

	// Act - Ejecutar el caso de uso
	result, err := useCase.CreateOrder(ctx, req)

	// Assert - Verificar resultados
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "generated-uuid-123", result.ID)
	assert.Equal(t, req.OrderNumber, result.OrderNumber)
	assert.Equal(t, req.ExternalID, result.ExternalID)
	assert.Equal(t, req.TotalAmount, result.TotalAmount)
	assert.Equal(t, req.Currency, result.Currency)

	// Verificar que se llamaron los métodos esperados
	mockRepo.AssertExpectations(t)
	mockRepo.AssertCalled(t, "OrderExists", ctx, req.ExternalID, req.IntegrationID)
	mockRepo.AssertCalled(t, "CreateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder"))
}

func TestCreateOrder_OrderAlreadyExists(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	businessID := uint(1)

	req := &dtos.CreateOrderRequest{
		BusinessID:      &businessID,
		IntegrationID:   10,
		IntegrationType: "shopify",
		ExternalID:      "EXT-12345",
		OrderNumber:     "ORD-001",
		TotalAmount:     115.0,
		Currency:        "USD",
	}

	// OrderExists retorna true (la orden ya existe)
	mockRepo.On("OrderExists", ctx, req.ExternalID, req.IntegrationID).
		Return(true, nil)

	// Act
	result, err := useCase.CreateOrder(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists")

	// Verificar que NO se llamó CreateOrder
	mockRepo.AssertNotCalled(t, "CreateOrder", mock.Anything, mock.Anything)
}

func TestCreateOrder_OrderExistsCheckError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	businessID := uint(1)

	req := &dtos.CreateOrderRequest{
		BusinessID:    &businessID,
		IntegrationID: 10,
		ExternalID:    "EXT-12345",
		OrderNumber:   "ORD-001",
	}

	// OrderExists retorna error de BD
	dbError := errors.New("database connection failed")
	mockRepo.On("OrderExists", ctx, req.ExternalID, req.IntegrationID).
		Return(false, dbError)

	// Act
	result, err := useCase.CreateOrder(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error checking if order exists")

	mockRepo.AssertExpectations(t)
}

func TestCreateOrder_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	businessID := uint(1)

	req := &dtos.CreateOrderRequest{
		BusinessID:    &businessID,
		IntegrationID: 10,
		ExternalID:    "EXT-12345",
		OrderNumber:   "ORD-001",
		TotalAmount:   115.0,
	}

	// OrderExists retorna false
	mockRepo.On("OrderExists", ctx, req.ExternalID, req.IntegrationID).
		Return(false, nil)

	// CreateOrder retorna error
	dbError := errors.New("constraint violation")
	mockRepo.On("CreateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).
		Return(dbError)

	// Act
	result, err := useCase.CreateOrder(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error creating order")

	mockRepo.AssertExpectations(t)
}

func TestCreateOrder_ValidatesRequiredFields(t *testing.T) {
	// Esta prueba verifica el mapeo correcto de campos obligatorios
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	businessID := uint(1)

	// Request con todos los campos opcionales
	req := &dtos.CreateOrderRequest{
		BusinessID:      &businessID,
		IntegrationID:   10,
		IntegrationType: "shopify",
		ExternalID:      "EXT-12345",
		OrderNumber:     "ORD-001",
		InternalNumber:  "INT-001",
		Platform:        "Shopify",
		Subtotal:        100.0,
		TotalAmount:     115.0,
		Currency:        "USD",
		Status:          "pending",
		PaymentMethodID: 1,
	}

	mockRepo.On("OrderExists", ctx, req.ExternalID, req.IntegrationID).
		Return(false, nil)

	var capturedOrder *entities.ProbabilityOrder
	mockRepo.On("CreateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).
		Run(func(args mock.Arguments) {
			capturedOrder = args.Get(1).(*entities.ProbabilityOrder)
			capturedOrder.ID = "test-id"
		}).
		Return(nil)

	mockLogger.On("Info", mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything).Maybe()
	mockRedisPublisher.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).
		Maybe().Return(nil)
	mockRabbitPublisher.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).
		Maybe().Return(nil)

	// Act
	_, err := useCase.CreateOrder(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, capturedOrder)

	// Verificar que los campos se mapearon correctamente
	assert.Equal(t, req.BusinessID, capturedOrder.BusinessID)
	assert.Equal(t, req.IntegrationID, capturedOrder.IntegrationID)
	assert.Equal(t, req.ExternalID, capturedOrder.ExternalID)
	assert.Equal(t, req.OrderNumber, capturedOrder.OrderNumber)
	assert.Equal(t, req.TotalAmount, capturedOrder.TotalAmount)
	assert.Equal(t, req.Currency, capturedOrder.Currency)

	mockRepo.AssertExpectations(t)
}

package usecaseorder

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetOrderByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	orderID := "order-uuid-123"

	// Orden esperada del repositorio
	businessID := uint(1)
	expectedOrder := &entities.ProbabilityOrder{
		ID:              orderID,
		BusinessID:      &businessID,
		IntegrationID:   10,
		IntegrationType: "shopify",
		Platform:        "Shopify",
		ExternalID:      "EXT-12345",
		OrderNumber:     "ORD-001",
		InternalNumber:  "INT-001",
		Subtotal:        100.0,
		Tax:             10.0,
		ShippingCost:    5.0,
		TotalAmount:     115.0,
		Currency:        "USD",
		CustomerName:    "John Doe",
		CustomerEmail:   "john@example.com",
		Status:          "pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Configurar mock del repositorio
	mockRepo.On("GetOrderByID", ctx, orderID).
		Return(expectedOrder, nil)

	// Act
	result, err := useCase.GetOrderByID(ctx, orderID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedOrder.ID, result.ID)
	assert.Equal(t, expectedOrder.OrderNumber, result.OrderNumber)
	assert.Equal(t, expectedOrder.ExternalID, result.ExternalID)
	assert.Equal(t, expectedOrder.TotalAmount, result.TotalAmount)
	assert.Equal(t, expectedOrder.Currency, result.Currency)
	assert.Equal(t, expectedOrder.CustomerEmail, result.CustomerEmail)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertCalled(t, "GetOrderByID", ctx, orderID)
}

func TestGetOrderByID_EmptyID(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()

	// Act - Pasar ID vacío
	result, err := useCase.GetOrderByID(ctx, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "order ID is required")

	// Verificar que NO se llamó al repositorio
	mockRepo.AssertNotCalled(t, "GetOrderByID")
}

func TestGetOrderByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	orderID := "non-existent-uuid"

	// Repositorio retorna error de no encontrado
	notFoundError := errors.New("record not found")
	mockRepo.On("GetOrderByID", ctx, orderID).
		Return(nil, notFoundError)

	// Act
	result, err := useCase.GetOrderByID(ctx, orderID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error getting order")

	mockRepo.AssertExpectations(t)
}

func TestGetOrderByID_DatabaseError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	orderID := "order-uuid-123"

	// Repositorio retorna error de conexión
	dbError := errors.New("database connection timeout")
	mockRepo.On("GetOrderByID", ctx, orderID).
		Return(nil, dbError)

	// Act
	result, err := useCase.GetOrderByID(ctx, orderID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error getting order")

	mockRepo.AssertExpectations(t)
}

func TestGetOrderByID_WithCompleteData(t *testing.T) {
	// Arrange - Test con una orden completa con todos los campos
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	orderID := "order-uuid-complete"
	businessID := uint(1)
	customerID := uint(100)
	warehouseID := uint(5)

	paidAt := time.Now().Add(-24 * time.Hour)
	deliveredAt := time.Now()
	trackingNumber := "TRACK-123"
	guideID := "GUIDE-456"

	expectedOrder := &entities.ProbabilityOrder{
		ID:              orderID,
		BusinessID:      &businessID,
		IntegrationID:   10,
		IntegrationType: "shopify",
		Platform:        "Shopify",
		ExternalID:      "EXT-FULL-12345",
		OrderNumber:     "ORD-FULL-001",
		InternalNumber:  "INT-FULL-001",
		Subtotal:        100.0,
		Tax:             10.0,
		Discount:        5.0,
		ShippingCost:    10.0,
		TotalAmount:     115.0,
		Currency:        "USD",
		CustomerID:      &customerID,
		CustomerName:    "Jane Smith",
		CustomerEmail:   "jane@example.com",
		CustomerPhone:   "+1234567890",
		CustomerDNI:     "12345678A",
		ShippingStreet:  "123 Main St",
		ShippingCity:    "New York",
		ShippingState:   "NY",
		ShippingCountry: "USA",
		PaymentMethodID: 1,
		IsPaid:          true,
		PaidAt:          &paidAt,
		TrackingNumber:  &trackingNumber,
		GuideID:         &guideID,
		DeliveredAt:     &deliveredAt,
		WarehouseID:     &warehouseID,
		WarehouseName:   "Main Warehouse",
		Status:          "delivered",
		OriginalStatus:  "fulfilled",
		CreatedAt:       time.Now().Add(-48 * time.Hour),
		UpdatedAt:       time.Now(),
	}

	mockRepo.On("GetOrderByID", ctx, orderID).
		Return(expectedOrder, nil)

	// Act
	result, err := useCase.GetOrderByID(ctx, orderID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedOrder.ID, result.ID)
	assert.Equal(t, expectedOrder.OrderNumber, result.OrderNumber)
	assert.Equal(t, expectedOrder.CustomerEmail, result.CustomerEmail)
	assert.Equal(t, expectedOrder.ShippingCity, result.ShippingCity)
	assert.Equal(t, expectedOrder.Status, result.Status)
	assert.True(t, result.IsPaid)

	mockRepo.AssertExpectations(t)
}

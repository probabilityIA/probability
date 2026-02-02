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

func TestListOrders_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	page := 1
	pageSize := 10
	filters := map[string]interface{}{
		"status": "pending",
	}

	businessID := uint(1)

	// Órdenes de ejemplo
	orders := []entities.ProbabilityOrder{
		{
			ID:              "order-1",
			BusinessID:      &businessID,
			OrderNumber:     "ORD-001",
			ExternalID:      "EXT-001",
			TotalAmount:     100.0,
			Currency:        "USD",
			Status:          "pending",
			CustomerEmail:   "customer1@example.com",
			IntegrationID:   10,
			IntegrationType: "shopify",
			CreatedAt:       time.Now(),
		},
		{
			ID:              "order-2",
			BusinessID:      &businessID,
			OrderNumber:     "ORD-002",
			ExternalID:      "EXT-002",
			TotalAmount:     200.0,
			Currency:        "USD",
			Status:          "pending",
			CustomerEmail:   "customer2@example.com",
			IntegrationID:   10,
			IntegrationType: "shopify",
			CreatedAt:       time.Now(),
		},
	}

	totalRecords := int64(25) // Total de registros en BD

	// Configurar mock
	mockRepo.On("ListOrders", ctx, page, pageSize, filters).
		Return(orders, totalRecords, nil)

	// Act
	result, err := useCase.ListOrders(ctx, page, pageSize, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, totalRecords, result.Total)
	assert.Equal(t, page, result.Page)
	assert.Equal(t, pageSize, result.PageSize)
	assert.Equal(t, 3, result.TotalPages) // 25 / 10 = 3 páginas

	// Verificar datos mapeados
	assert.Equal(t, "order-1", result.Data[0].ID)
	assert.Equal(t, "ORD-001", result.Data[0].OrderNumber)
	assert.Equal(t, 100.0, result.Data[0].TotalAmount)

	mockRepo.AssertExpectations(t)
}

func TestListOrders_EmptyResult(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	page := 1
	pageSize := 10
	filters := map[string]interface{}{
		"status": "cancelled",
	}

	// No hay órdenes que coincidan con los filtros
	emptyOrders := []entities.ProbabilityOrder{}
	totalRecords := int64(0)

	mockRepo.On("ListOrders", ctx, page, pageSize, filters).
		Return(emptyOrders, totalRecords, nil)

	// Act
	result, err := useCase.ListOrders(ctx, page, pageSize, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 0)
	assert.Equal(t, int64(0), result.Total)
	assert.Equal(t, 0, result.TotalPages)

	mockRepo.AssertExpectations(t)
}

func TestListOrders_PaginationValidation(t *testing.T) {
	tests := []struct {
		name         string
		inputPage    int
		inputSize    int
		expectedPage int
		expectedSize int
	}{
		{
			name:         "página negativa se ajusta a 1",
			inputPage:    -5,
			inputSize:    10,
			expectedPage: 1,
			expectedSize: 10,
		},
		{
			name:         "página cero se ajusta a 1",
			inputPage:    0,
			inputSize:    10,
			expectedPage: 1,
			expectedSize: 10,
		},
		{
			name:         "pageSize negativo se ajusta a 10",
			inputPage:    1,
			inputSize:    -10,
			expectedPage: 1,
			expectedSize: 10,
		},
		{
			name:         "pageSize cero se ajusta a 10",
			inputPage:    1,
			inputSize:    0,
			expectedPage: 1,
			expectedSize: 10,
		},
		{
			name:         "pageSize mayor a 100 se ajusta a 10",
			inputPage:    1,
			inputSize:    150,
			expectedPage: 1,
			expectedSize: 10,
		},
		{
			name:         "valores válidos no se modifican",
			inputPage:    3,
			inputSize:    25,
			expectedPage: 3,
			expectedSize: 25,
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
			filters := map[string]interface{}{}

			// Capturar los valores que se pasan al repositorio
			mockRepo.On("ListOrders", ctx, tt.expectedPage, tt.expectedSize, filters).
				Return([]entities.ProbabilityOrder{}, int64(0), nil)

			// Act
			result, err := useCase.ListOrders(ctx, tt.inputPage, tt.inputSize, filters)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPage, result.Page)
			assert.Equal(t, tt.expectedSize, result.PageSize)

			mockRepo.AssertCalled(t, "ListOrders", ctx, tt.expectedPage, tt.expectedSize, filters)
		})
	}
}

func TestListOrders_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	page := 1
	pageSize := 10
	filters := map[string]interface{}{}

	// Repositorio retorna error
	dbError := errors.New("database timeout")
	mockRepo.On("ListOrders", ctx, page, pageSize, filters).
		Return(nil, int64(0), dbError)

	// Act
	result, err := useCase.ListOrders(ctx, page, pageSize, filters)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error listing orders")

	mockRepo.AssertExpectations(t)
}

func TestListOrders_TotalPagesCalculation(t *testing.T) {
	tests := []struct {
		name              string
		totalRecords      int64
		pageSize          int
		expectedTotalPage int
	}{
		{
			name:              "sin registros",
			totalRecords:      0,
			pageSize:          10,
			expectedTotalPage: 0,
		},
		{
			name:              "menos registros que pageSize",
			totalRecords:      5,
			pageSize:          10,
			expectedTotalPage: 1,
		},
		{
			name:              "registros exactos para una página",
			totalRecords:      10,
			pageSize:          10,
			expectedTotalPage: 1,
		},
		{
			name:              "registros que requieren múltiples páginas",
			totalRecords:      25,
			pageSize:          10,
			expectedTotalPage: 3,
		},
		{
			name:              "registros exactos para múltiples páginas",
			totalRecords:      100,
			pageSize:          25,
			expectedTotalPage: 4,
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
			page := 1
			filters := map[string]interface{}{}

			mockRepo.On("ListOrders", ctx, page, tt.pageSize, filters).
				Return([]entities.ProbabilityOrder{}, tt.totalRecords, nil)

			// Act
			result, err := useCase.ListOrders(ctx, page, tt.pageSize, filters)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedTotalPage, result.TotalPages)
		})
	}
}

func TestListOrders_WithFilters(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.RepositoryMock)
	mockRedisPublisher := new(mocks.EventPublisherMock)
	mockRabbitPublisher := new(mocks.RabbitPublisherMock)
	mockLogger := new(mocks.LoggerMock)
	mockScoreUseCase := new(mocks.ScoreUseCaseMock)

	useCase := New(mockRepo, mockRedisPublisher, mockRabbitPublisher, mockLogger, mockScoreUseCase)

	ctx := context.Background()
	page := 1
	pageSize := 10

	// Filtros complejos
	filters := map[string]interface{}{
		"status":          "pending",
		"business_id":     uint(1),
		"integration_id":  uint(10),
		"customer_email":  "john@example.com",
		"total_min":       50.0,
		"total_max":       500.0,
		"date_from":       "2024-01-01",
		"date_to":         "2024-12-31",
	}

	businessID := uint(1)
	orders := []entities.ProbabilityOrder{
		{
			ID:            "order-1",
			BusinessID:    &businessID,
			OrderNumber:   "ORD-001",
			Status:        "pending",
			TotalAmount:   100.0,
			CustomerEmail: "john@example.com",
			IntegrationID: 10,
		},
	}

	mockRepo.On("ListOrders", ctx, page, pageSize, filters).
		Return(orders, int64(1), nil)

	// Act
	result, err := useCase.ListOrders(ctx, page, pageSize, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, "john@example.com", result.Data[0].CustomerEmail)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertCalled(t, "ListOrders", ctx, page, pageSize, filters)
}

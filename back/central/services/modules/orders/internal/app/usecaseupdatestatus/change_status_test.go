package usecaseupdatestatus

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- helpers ---

func setupUseCase() (*UseCaseUpdateStatus, *mocks.RepositoryMock, *mocks.RabbitPublisherMock, *mocks.LoggerMock) {
	repo := new(mocks.RepositoryMock)
	rabbit := new(mocks.RabbitPublisherMock)
	logger := new(mocks.LoggerMock)

	uc := &UseCaseUpdateStatus{
		repo:                 repo,
		logger:               logger,
		rabbitEventPublisher: rabbit,
	}
	return uc, repo, rabbit, logger
}

func newOrder(id, status string) *entities.ProbabilityOrder {
	businessID := uint(1)
	return &entities.ProbabilityOrder{
		ID:            id,
		BusinessID:    &businessID,
		IntegrationID: 10,
		OrderNumber:   "ORD-001",
		InternalNumber: "INT-001",
		Status:        status,
		TotalAmount:   50000,
		Currency:      "COP",
		CustomerEmail: "test@example.com",
		Platform:      "shopify",
	}
}

func userID(id uint) *uint { return &id }

// --- Test: Validación de entrada ---

func TestChangeStatus_EmptyOrderID(t *testing.T) {
	uc, _, _, _ := setupUseCase()

	result, err := uc.ChangeStatus(context.Background(), "", &dtos.ChangeStatusRequest{Status: "picking"})

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order ID is required")
}

func TestChangeStatus_EmptyStatus(t *testing.T) {
	uc, _, _, _ := setupUseCase()

	result, err := uc.ChangeStatus(context.Background(), "order-1", &dtos.ChangeStatusRequest{Status: ""})

	assert.Nil(t, result)
	assert.ErrorIs(t, err, domainerrors.ErrInvalidStatus)
}

func TestChangeStatus_InvalidStatus(t *testing.T) {
	uc, _, _, _ := setupUseCase()

	result, err := uc.ChangeStatus(context.Background(), "order-1", &dtos.ChangeStatusRequest{Status: "nonexistent"})

	assert.Nil(t, result)
	assert.ErrorIs(t, err, domainerrors.ErrInvalidStatus)
}

// --- Test: Orden no encontrada ---

func TestChangeStatus_OrderNotFound(t *testing.T) {
	uc, repo, _, _ := setupUseCase()
	ctx := context.Background()

	repo.On("GetOrderByID", ctx, "order-1").Return(nil, errors.New("record not found"))

	result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{Status: "picking"})

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting order")
	repo.AssertExpectations(t)
}

// --- Test: Estado terminal ---

func TestChangeStatus_OrderInTerminalState(t *testing.T) {
	terminalStatuses := []string{"cancelled", "refunded"}

	for _, status := range terminalStatuses {
		t.Run("from_"+status, func(t *testing.T) {
			uc, repo, _, _ := setupUseCase()
			ctx := context.Background()
			order := newOrder("order-1", status)

			repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)

			result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{Status: "picking"})

			assert.Nil(t, result)
			assert.ErrorIs(t, err, domainerrors.ErrOrderInTerminalState)
			repo.AssertExpectations(t)
		})
	}
}

// --- Test: Transición no permitida ---

func TestChangeStatus_InvalidTransition(t *testing.T) {
	uc, repo, _, _ := setupUseCase()
	ctx := context.Background()
	order := newOrder("order-1", "pending")

	repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)

	// pending -> delivered is not a valid transition
	result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{Status: "delivered"})

	assert.Nil(t, result)
	assert.ErrorIs(t, err, domainerrors.ErrInvalidStatusTransition)
	repo.AssertExpectations(t)
}

// --- Test: Flujo exitoso simple (pending -> picking) ---

func TestChangeStatus_Success_PendingToPicking(t *testing.T) {
	uc, repo, rabbit, logger := setupUseCase()
	ctx := context.Background()
	order := newOrder("order-1", "pending")
	statusID := uint(3)

	repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
	repo.On("GetOrderStatusIDByCode", ctx, "picking").Return(&statusID, nil)
	repo.On("UpdateOrder", ctx, mock.MatchedBy(func(o *entities.ProbabilityOrder) bool {
		return o.Status == "picking" && o.StatusID != nil && *o.StatusID == statusID
	})).Return(nil)
	repo.On("CreateOrderHistory", ctx, mock.MatchedBy(func(h *entities.OrderHistory) bool {
		return h.OrderID == "order-1" &&
			h.PreviousStatus == "pending" &&
			h.NewStatus == "picking"
	})).Return(nil)
	logger.On("Info", mock.Anything).Maybe()
	logger.On("Error", mock.Anything).Maybe()
	rabbit.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{
		Status: "picking",
		UserID: userID(5),
		UserName: "admin",
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "picking", result.Status)
	repo.AssertExpectations(t)
}

// --- Test: Cancelled accesible desde cualquier estado no-terminal ---

func TestChangeStatus_Success_CancelledFromMultipleStatuses(t *testing.T) {
	statuses := []string{"pending", "picking", "packing", "ready_to_ship", "in_transit", "out_for_delivery"}

	for _, status := range statuses {
		t.Run("cancel_from_"+status, func(t *testing.T) {
			uc, repo, rabbit, logger := setupUseCase()
			ctx := context.Background()
			order := newOrder("order-1", status)
			statusID := uint(20)

			repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
			repo.On("GetOrderStatusIDByCode", ctx, "cancelled").Return(&statusID, nil)
			repo.On("UpdateOrder", ctx, mock.MatchedBy(func(o *entities.ProbabilityOrder) bool {
				return o.Status == "cancelled"
			})).Return(nil)
			repo.On("CreateOrderHistory", ctx, mock.Anything).Return(nil)
			logger.On("Info", mock.Anything).Maybe()
			logger.On("Warn", mock.Anything).Maybe()
			logger.On("Error", mock.Anything).Maybe()
			rabbit.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

			result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{
				Status:   "cancelled",
				Metadata: map[string]interface{}{"reason": "client request"},
			})

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "cancelled", result.Status)
			repo.AssertExpectations(t)
		})
	}
}

// --- Test: UpdateOrder falla ---

func TestChangeStatus_UpdateOrderError(t *testing.T) {
	uc, repo, _, logger := setupUseCase()
	ctx := context.Background()
	order := newOrder("order-1", "pending")
	statusID := uint(3)

	repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
	repo.On("GetOrderStatusIDByCode", ctx, "picking").Return(&statusID, nil)
	repo.On("UpdateOrder", ctx, mock.Anything).Return(errors.New("db connection lost"))
	logger.On("Warn", mock.Anything).Maybe()

	result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{Status: "picking"})

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error updating order")
	repo.AssertExpectations(t)
}

// --- Test: StatusID no se resuelve (warn, no error) ---

func TestChangeStatus_StatusIDNotResolved_ContinuesSuccessfully(t *testing.T) {
	uc, repo, rabbit, logger := setupUseCase()
	ctx := context.Background()
	order := newOrder("order-1", "pending")

	repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
	repo.On("GetOrderStatusIDByCode", ctx, "picking").Return(nil, errors.New("status not found"))
	repo.On("UpdateOrder", ctx, mock.MatchedBy(func(o *entities.ProbabilityOrder) bool {
		return o.Status == "picking" && o.StatusID == nil
	})).Return(nil)
	repo.On("CreateOrderHistory", ctx, mock.Anything).Return(nil)
	logger.On("Warn", mock.Anything).Maybe()
	logger.On("Info", mock.Anything).Maybe()
	logger.On("Error", mock.Anything).Maybe()
	rabbit.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{Status: "picking"})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	repo.AssertExpectations(t)
}

// --- Test: Strategies con metadata ---

func TestChangeStatus_Strategy_AssignedToDriver_ExtractsMetadata(t *testing.T) {
	uc, repo, rabbit, logger := setupUseCase()
	ctx := context.Background()
	order := newOrder("order-1", "ready_to_ship")
	statusID := uint(7)

	repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
	repo.On("GetOrderStatusIDByCode", ctx, "assigned_to_driver").Return(&statusID, nil)

	var captured *entities.ProbabilityOrder
	repo.On("UpdateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).
		Run(func(args mock.Arguments) { captured = args.Get(1).(*entities.ProbabilityOrder) }).
		Return(nil)
	repo.On("CreateOrderHistory", ctx, mock.Anything).Return(nil)
	logger.On("Info", mock.Anything).Maybe()
	logger.On("Error", mock.Anything).Maybe()
	rabbit.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{
		Status: "assigned_to_driver",
		Metadata: map[string]interface{}{
			"driver_id":   float64(42),
			"driver_name": "Carlos",
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "assigned_to_driver", captured.Status)
	assert.NotNil(t, captured.DriverID)
	assert.Equal(t, uint(42), *captured.DriverID)
	assert.Equal(t, "Carlos", captured.DriverName)
	repo.AssertExpectations(t)
}

func TestChangeStatus_Strategy_InTransit_ExtractsTrackingInfo(t *testing.T) {
	uc, repo, rabbit, logger := setupUseCase()
	ctx := context.Background()
	order := newOrder("order-1", "picked_up")
	statusID := uint(9)

	repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
	repo.On("GetOrderStatusIDByCode", ctx, "in_transit").Return(&statusID, nil)

	var captured *entities.ProbabilityOrder
	repo.On("UpdateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).
		Run(func(args mock.Arguments) { captured = args.Get(1).(*entities.ProbabilityOrder) }).
		Return(nil)
	repo.On("CreateOrderHistory", ctx, mock.Anything).Return(nil)
	logger.On("Info", mock.Anything).Maybe()
	logger.On("Error", mock.Anything).Maybe()
	rabbit.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{
		Status: "in_transit",
		Metadata: map[string]interface{}{
			"tracking_number": "TRACK-12345",
			"tracking_link":   "https://tracking.example.com/TRACK-12345",
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "in_transit", captured.Status)
	assert.NotNil(t, captured.TrackingNumber)
	assert.Equal(t, "TRACK-12345", *captured.TrackingNumber)
	assert.NotNil(t, captured.TrackingLink)
	assert.Equal(t, "https://tracking.example.com/TRACK-12345", *captured.TrackingLink)
	repo.AssertExpectations(t)
}

func TestChangeStatus_Strategy_Delivered_SetsDeliveredAt(t *testing.T) {
	uc, repo, rabbit, logger := setupUseCase()
	ctx := context.Background()
	order := newOrder("order-1", "out_for_delivery")
	statusID := uint(11)

	repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
	repo.On("GetOrderStatusIDByCode", ctx, "delivered").Return(&statusID, nil)

	var captured *entities.ProbabilityOrder
	repo.On("UpdateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).
		Run(func(args mock.Arguments) { captured = args.Get(1).(*entities.ProbabilityOrder) }).
		Return(nil)
	repo.On("CreateOrderHistory", ctx, mock.Anything).Return(nil)
	logger.On("Info", mock.Anything).Maybe()
	logger.On("Error", mock.Anything).Maybe()
	rabbit.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	before := time.Now()
	result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{Status: "delivered"})
	after := time.Now()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "delivered", captured.Status)
	assert.NotNil(t, captured.DeliveredAt)
	assert.True(t, captured.DeliveredAt.After(before) || captured.DeliveredAt.Equal(before))
	assert.True(t, captured.DeliveredAt.Before(after) || captured.DeliveredAt.Equal(after))
	repo.AssertExpectations(t)
}

func TestChangeStatus_Strategy_Cancelled_ExtractsReason(t *testing.T) {
	uc, repo, rabbit, logger := setupUseCase()
	ctx := context.Background()
	order := newOrder("order-1", "pending")
	statusID := uint(19)

	repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
	repo.On("GetOrderStatusIDByCode", ctx, "cancelled").Return(&statusID, nil)

	var captured *entities.ProbabilityOrder
	repo.On("UpdateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).
		Run(func(args mock.Arguments) { captured = args.Get(1).(*entities.ProbabilityOrder) }).
		Return(nil)
	repo.On("CreateOrderHistory", ctx, mock.MatchedBy(func(h *entities.OrderHistory) bool {
		return h.PreviousStatus == "pending" &&
			h.NewStatus == "cancelled" &&
			h.Reason != nil && *h.Reason == "client request"
	})).Return(nil)
	logger.On("Info", mock.Anything).Maybe()
	logger.On("Error", mock.Anything).Maybe()
	rabbit.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{
		Status:   "cancelled",
		Metadata: map[string]interface{}{"reason": "client request"},
		UserID:   userID(5),
		UserName: "admin",
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "cancelled", captured.Status)
	assert.NotNil(t, captured.Notes)
	assert.Equal(t, "client request", *captured.Notes)
	repo.AssertExpectations(t)
}

func TestChangeStatus_Strategy_DeliveryNovelty_SetsNovelty(t *testing.T) {
	uc, repo, rabbit, logger := setupUseCase()
	ctx := context.Background()
	order := newOrder("order-1", "out_for_delivery")
	statusID := uint(12)

	repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
	repo.On("GetOrderStatusIDByCode", ctx, "delivery_novelty").Return(&statusID, nil)

	var captured *entities.ProbabilityOrder
	repo.On("UpdateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).
		Run(func(args mock.Arguments) { captured = args.Get(1).(*entities.ProbabilityOrder) }).
		Return(nil)
	repo.On("CreateOrderHistory", ctx, mock.Anything).Return(nil)
	logger.On("Info", mock.Anything).Maybe()
	logger.On("Error", mock.Anything).Maybe()
	rabbit.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{
		Status:   "delivery_novelty",
		Metadata: map[string]interface{}{"reason": "client not home"},
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "delivery_novelty", captured.Status)
	assert.NotNil(t, captured.Novelty)
	assert.Equal(t, "client not home", *captured.Novelty)
	repo.AssertExpectations(t)
}

func TestChangeStatus_Strategy_InventoryIssue_SetsNovelty(t *testing.T) {
	uc, repo, rabbit, logger := setupUseCase()
	ctx := context.Background()
	order := newOrder("order-1", "picking")
	statusID := uint(6)

	repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
	repo.On("GetOrderStatusIDByCode", ctx, "inventory_issue").Return(&statusID, nil)

	var captured *entities.ProbabilityOrder
	repo.On("UpdateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).
		Run(func(args mock.Arguments) { captured = args.Get(1).(*entities.ProbabilityOrder) }).
		Return(nil)
	repo.On("CreateOrderHistory", ctx, mock.Anything).Return(nil)
	logger.On("Info", mock.Anything).Maybe()
	logger.On("Error", mock.Anything).Maybe()
	rabbit.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{
		Status:   "inventory_issue",
		Metadata: map[string]interface{}{"notes": "product damaged in warehouse"},
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "inventory_issue", captured.Status)
	assert.NotNil(t, captured.Novelty)
	assert.Equal(t, "product damaged in warehouse", *captured.Novelty)
	repo.AssertExpectations(t)
}

func TestChangeStatus_Strategy_OnHold_ExtractsReason(t *testing.T) {
	uc, repo, rabbit, logger := setupUseCase()
	ctx := context.Background()
	order := newOrder("order-1", "pending")
	statusID := uint(2)

	repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
	repo.On("GetOrderStatusIDByCode", ctx, "on_hold").Return(&statusID, nil)

	var captured *entities.ProbabilityOrder
	repo.On("UpdateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).
		Run(func(args mock.Arguments) { captured = args.Get(1).(*entities.ProbabilityOrder) }).
		Return(nil)
	repo.On("CreateOrderHistory", ctx, mock.Anything).Return(nil)
	logger.On("Info", mock.Anything).Maybe()
	logger.On("Error", mock.Anything).Maybe()
	rabbit.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{
		Status:   "on_hold",
		Metadata: map[string]interface{}{"reason": "awaiting payment confirmation"},
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "on_hold", captured.Status)
	assert.NotNil(t, captured.Notes)
	assert.Equal(t, "awaiting payment confirmation", *captured.Notes)
	repo.AssertExpectations(t)
}

// --- Test: Strategies sin metadata (no crashean) ---

func TestChangeStatus_Strategy_AssignedToDriver_NoMetadata(t *testing.T) {
	uc, repo, rabbit, logger := setupUseCase()
	ctx := context.Background()
	order := newOrder("order-1", "ready_to_ship")
	statusID := uint(7)

	repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
	repo.On("GetOrderStatusIDByCode", ctx, "assigned_to_driver").Return(&statusID, nil)

	var captured *entities.ProbabilityOrder
	repo.On("UpdateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).
		Run(func(args mock.Arguments) { captured = args.Get(1).(*entities.ProbabilityOrder) }).
		Return(nil)
	repo.On("CreateOrderHistory", ctx, mock.Anything).Return(nil)
	logger.On("Info", mock.Anything).Maybe()
	logger.On("Error", mock.Anything).Maybe()
	rabbit.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{
		Status: "assigned_to_driver",
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "assigned_to_driver", captured.Status)
	assert.Nil(t, captured.DriverID)
	assert.Empty(t, captured.DriverName)
	repo.AssertExpectations(t)
}

// --- Test: Historial registra usuario y metadata ---

func TestChangeStatus_History_RecordsUserAndMetadata(t *testing.T) {
	uc, repo, rabbit, logger := setupUseCase()
	ctx := context.Background()
	order := newOrder("order-1", "pending")
	statusID := uint(3)

	repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
	repo.On("GetOrderStatusIDByCode", ctx, "picking").Return(&statusID, nil)
	repo.On("UpdateOrder", ctx, mock.Anything).Return(nil)

	var capturedHistory *entities.OrderHistory
	repo.On("CreateOrderHistory", ctx, mock.AnythingOfType("*entities.OrderHistory")).
		Run(func(args mock.Arguments) { capturedHistory = args.Get(1).(*entities.OrderHistory) }).
		Return(nil)
	logger.On("Info", mock.Anything).Maybe()
	logger.On("Error", mock.Anything).Maybe()
	rabbit.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	uid := uint(10)
	_, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{
		Status:   "picking",
		UserID:   &uid,
		UserName: "warehouse_op",
	})

	assert.NoError(t, err)
	assert.NotNil(t, capturedHistory)
	assert.Equal(t, "order-1", capturedHistory.OrderID)
	assert.Equal(t, "pending", capturedHistory.PreviousStatus)
	assert.Equal(t, "picking", capturedHistory.NewStatus)
	assert.NotNil(t, capturedHistory.ChangedBy)
	assert.Equal(t, uint(10), *capturedHistory.ChangedBy)
	assert.Equal(t, "warehouse_op", capturedHistory.ChangedByName)
	assert.Nil(t, capturedHistory.Reason) // no reason in metadata
	repo.AssertExpectations(t)
}

func TestChangeStatus_History_RecordsReasonFromMetadata(t *testing.T) {
	uc, repo, rabbit, logger := setupUseCase()
	ctx := context.Background()
	order := newOrder("order-1", "out_for_delivery")
	statusID := uint(14)

	repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
	repo.On("GetOrderStatusIDByCode", ctx, "delivery_failed").Return(&statusID, nil)
	repo.On("UpdateOrder", ctx, mock.Anything).Return(nil)

	var capturedHistory *entities.OrderHistory
	repo.On("CreateOrderHistory", ctx, mock.AnythingOfType("*entities.OrderHistory")).
		Run(func(args mock.Arguments) { capturedHistory = args.Get(1).(*entities.OrderHistory) }).
		Return(nil)
	logger.On("Info", mock.Anything).Maybe()
	logger.On("Error", mock.Anything).Maybe()
	rabbit.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	_, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{
		Status:   "delivery_failed",
		Metadata: map[string]interface{}{"reason": "wrong address"},
	})

	assert.NoError(t, err)
	assert.NotNil(t, capturedHistory)
	assert.NotNil(t, capturedHistory.Reason)
	assert.Equal(t, "wrong address", *capturedHistory.Reason)
	assert.NotEmpty(t, capturedHistory.Metadata)
	repo.AssertExpectations(t)
}

// --- Test: History error no bloquea el flujo ---

func TestChangeStatus_HistoryError_DoesNotBlockFlow(t *testing.T) {
	uc, repo, rabbit, logger := setupUseCase()
	ctx := context.Background()
	order := newOrder("order-1", "pending")
	statusID := uint(3)

	repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
	repo.On("GetOrderStatusIDByCode", ctx, "picking").Return(&statusID, nil)
	repo.On("UpdateOrder", ctx, mock.Anything).Return(nil)
	repo.On("CreateOrderHistory", ctx, mock.Anything).Return(errors.New("history table locked"))
	logger.On("Info", mock.Anything).Maybe()
	logger.On("Error", mock.Anything).Maybe()
	rabbit.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{Status: "picking"})

	// Should succeed despite history error (it only logs the error)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	repo.AssertExpectations(t)
}

// --- Test: Flujo completo de la cadena logística ---

func TestChangeStatus_FullLogisticsFlow(t *testing.T) {
	flow := []struct {
		from string
		to   string
	}{
		{"pending", "picking"},
		{"picking", "packing"},
		{"packing", "ready_to_ship"},
		{"ready_to_ship", "assigned_to_driver"},
		{"assigned_to_driver", "picked_up"},
		{"picked_up", "in_transit"},
		{"in_transit", "out_for_delivery"},
		{"out_for_delivery", "delivered"},
		{"delivered", "completed"},
	}

	for _, step := range flow {
		t.Run(step.from+"_to_"+step.to, func(t *testing.T) {
			uc, repo, rabbit, logger := setupUseCase()
			ctx := context.Background()
			order := newOrder("order-1", step.from)
			statusID := uint(1)

			repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
			repo.On("GetOrderStatusIDByCode", ctx, step.to).Return(&statusID, nil)
			repo.On("UpdateOrder", ctx, mock.MatchedBy(func(o *entities.ProbabilityOrder) bool {
				return o.Status == step.to
			})).Return(nil)
			repo.On("CreateOrderHistory", ctx, mock.Anything).Return(nil)
			logger.On("Info", mock.Anything).Maybe()
			logger.On("Error", mock.Anything).Maybe()
			rabbit.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

			result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{Status: step.to})

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, step.to, result.Status)
			repo.AssertExpectations(t)
		})
	}
}

// --- Test: Flujo con novedad y reintento ---

func TestChangeStatus_DeliveryNovelty_RetryFlow(t *testing.T) {
	// delivery_novelty -> assigned_to_driver (retry)
	uc, repo, rabbit, logger := setupUseCase()
	ctx := context.Background()
	order := newOrder("order-1", "delivery_novelty")
	statusID := uint(7)

	repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
	repo.On("GetOrderStatusIDByCode", ctx, "assigned_to_driver").Return(&statusID, nil)
	repo.On("UpdateOrder", ctx, mock.Anything).Return(nil)
	repo.On("CreateOrderHistory", ctx, mock.Anything).Return(nil)
	logger.On("Info", mock.Anything).Maybe()
	logger.On("Error", mock.Anything).Maybe()
	rabbit.On("PublishOrderEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{
		Status: "assigned_to_driver",
		Metadata: map[string]interface{}{
			"driver_id":   float64(99),
			"driver_name": "Pedro",
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "assigned_to_driver", result.Status)
	repo.AssertExpectations(t)
}

// --- Test: Nil rabbit publisher (no crash) ---

func TestChangeStatus_NilRabbitPublisher_NoCrash(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	logger := new(mocks.LoggerMock)

	uc := &UseCaseUpdateStatus{
		repo:                 repo,
		logger:               logger,
		rabbitEventPublisher: nil, // no publisher
	}

	ctx := context.Background()
	order := newOrder("order-1", "pending")
	statusID := uint(3)

	repo.On("GetOrderByID", ctx, "order-1").Return(order, nil)
	repo.On("GetOrderStatusIDByCode", ctx, "picking").Return(&statusID, nil)
	repo.On("UpdateOrder", ctx, mock.Anything).Return(nil)
	repo.On("CreateOrderHistory", ctx, mock.Anything).Return(nil)
	logger.On("Info", mock.Anything).Maybe()
	logger.On("Error", mock.Anything).Maybe()

	result, err := uc.ChangeStatus(ctx, "order-1", &dtos.ChangeStatusRequest{Status: "picking"})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	repo.AssertExpectations(t)
}

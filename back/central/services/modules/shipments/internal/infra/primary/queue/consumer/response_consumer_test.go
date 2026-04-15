package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/mocks"
)

func newTestConsumer(repo domain.IRepository, ssePublisher domain.IShipmentSSEPublisher) *ResponseConsumer {
	return &ResponseConsumer{
		queue:        &mocks.RabbitMQMock{},
		repo:         repo,
		log:          mocks.NewLoggerMock(),
		ssePublisher: ssePublisher,
		redisClient:  nil,
	}
}

func shipmentIDPtr(id uint) *uint {
	return &id
}

func orderIDPtr(s string) *string {
	return &s
}

func buildCancelMessage(shipmentID *uint, businessID uint, status string, errMsg string) []byte {
	msg := TransportResponseMessage{
		ShipmentID:    shipmentID,
		BusinessID:    businessID,
		Provider:      "envioclick",
		Operation:     "cancel",
		Status:        status,
		CorrelationID: "corr-001",
		Timestamp:     time.Now(),
		Error:         errMsg,
	}
	b, _ := json.Marshal(msg)
	return b
}

func TestHandleCancelResponse_Success_UpdatesShipmentAndClearsOrder(t *testing.T) {
	shipmentID := uint(42)
	orderID := "order-uuid-123"
	businessID := uint(7)

	updatedStatus := ""
	clearedOrderID := ""
	cancelledShipmentID := uint(0)
	cancelledBusinessID := uint(0)

	repoMock := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) {
			return &domain.Shipment{
				ID:      id,
				OrderID: orderIDPtr(orderID),
				Status:  "pending",
			}, nil
		},
		UpdateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error {
			updatedStatus = shipment.Status
			return nil
		},
		ClearOrderGuideDataFn: func(ctx context.Context, id string) error {
			clearedOrderID = id
			return nil
		},
	}

	sseMock := &mocks.SSEPublisherMock{
		PublishShipmentCancelledFn: func(ctx context.Context, bID uint, sID uint) {
			cancelledBusinessID = bID
			cancelledShipmentID = sID
		},
	}

	consumer := newTestConsumer(repoMock, sseMock)
	msg := buildCancelMessage(shipmentIDPtr(shipmentID), businessID, "success", "")

	err := consumer.handleResponse(msg)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updatedStatus != "cancelled" {
		t.Errorf("expected shipment status 'cancelled', got '%s'", updatedStatus)
	}
	if clearedOrderID != orderID {
		t.Errorf("expected ClearOrderGuideData called with '%s', got '%s'", orderID, clearedOrderID)
	}
	if cancelledShipmentID != shipmentID {
		t.Errorf("expected SSE shipment ID %d, got %d", shipmentID, cancelledShipmentID)
	}
	if cancelledBusinessID != businessID {
		t.Errorf("expected SSE business ID %d, got %d", businessID, cancelledBusinessID)
	}
}

func TestHandleCancelResponse_Success_ClearsOrderGuideData(t *testing.T) {
	shipmentID := uint(10)
	orderID := "order-sync-456"

	clearCalled := false

	repoMock := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) {
			return &domain.Shipment{
				ID:      id,
				OrderID: orderIDPtr(orderID),
				Status:  "pending",
			}, nil
		},
		UpdateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error {
			return nil
		},
		ClearOrderGuideDataFn: func(ctx context.Context, id string) error {
			clearCalled = true
			return nil
		},
	}

	sseMock := &mocks.SSEPublisherMock{}

	consumer := newTestConsumer(repoMock, sseMock)
	msg := buildCancelMessage(shipmentIDPtr(shipmentID), 5, "success", "")

	consumer.handleResponse(msg)

	if !clearCalled {
		t.Error("expected ClearOrderGuideData to be called, but it was not")
	}
}

func TestHandleCancelResponse_Error_PublishesCancelFailed(t *testing.T) {
	shipmentID := uint(99)
	businessID := uint(3)
	errMsg := "carrier rejected cancellation"

	cancelFailedCalled := false
	cancelFailedShipmentID := uint(0)
	cancelFailedBusinessID := uint(0)
	cancelFailedErrMsg := ""
	cancelFailedCorrelationID := ""

	repoMock := &mocks.RepositoryMock{}

	sseMock := &mocks.SSEPublisherMock{
		PublishCancelFailedFn: func(ctx context.Context, bID uint, sID uint, corrID string, msg string) {
			cancelFailedCalled = true
			cancelFailedBusinessID = bID
			cancelFailedShipmentID = sID
			cancelFailedCorrelationID = corrID
			cancelFailedErrMsg = msg
		},
	}

	consumer := newTestConsumer(repoMock, sseMock)
	msg := buildCancelMessage(shipmentIDPtr(shipmentID), businessID, "error", errMsg)

	err := consumer.handleResponse(msg)

	if err != nil {
		t.Fatalf("expected no error from handleResponse, got %v", err)
	}
	if !cancelFailedCalled {
		t.Fatal("expected PublishCancelFailed to be called, but it was not")
	}
	if cancelFailedBusinessID != businessID {
		t.Errorf("expected business ID %d, got %d", businessID, cancelFailedBusinessID)
	}
	if cancelFailedShipmentID != shipmentID {
		t.Errorf("expected shipment ID %d, got %d", shipmentID, cancelFailedShipmentID)
	}
	if cancelFailedErrMsg != errMsg {
		t.Errorf("expected error msg '%s', got '%s'", errMsg, cancelFailedErrMsg)
	}
	if cancelFailedCorrelationID != "corr-001" {
		t.Errorf("expected correlation ID 'corr-001', got '%s'", cancelFailedCorrelationID)
	}
}

func TestHandleCancelResponse_Error_DoesNotUpdateRepository(t *testing.T) {
	updateCalled := false

	repoMock := &mocks.RepositoryMock{
		UpdateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error {
			updateCalled = true
			return nil
		},
		ClearOrderGuideDataFn: func(ctx context.Context, orderID string) error {
			updateCalled = true
			return nil
		},
	}

	sseMock := &mocks.SSEPublisherMock{}

	consumer := newTestConsumer(repoMock, sseMock)
	msg := buildCancelMessage(shipmentIDPtr(5), 1, "error", "some error")

	consumer.handleResponse(msg)

	if updateCalled {
		t.Error("expected repository not to be called on error response, but it was")
	}
}

func TestHandleCancelResponse_NoShipmentID_OnlyPublishesSSE(t *testing.T) {
	businessID := uint(8)

	getShipmentCalled := false
	updateShipmentCalled := false
	clearOrderCalled := false
	cancelledPublished := false
	cancelFailedPublished := false

	repoMock := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) {
			getShipmentCalled = true
			return nil, nil
		},
		UpdateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error {
			updateShipmentCalled = true
			return nil
		},
		ClearOrderGuideDataFn: func(ctx context.Context, orderID string) error {
			clearOrderCalled = true
			return nil
		},
	}

	sseMock := &mocks.SSEPublisherMock{
		PublishShipmentCancelledFn: func(ctx context.Context, bID uint, sID uint) {
			cancelledPublished = true
		},
		PublishCancelFailedFn: func(ctx context.Context, bID uint, sID uint, corrID string, msg string) {
			cancelFailedPublished = true
		},
	}

	consumer := newTestConsumer(repoMock, sseMock)
	msg := buildCancelMessage(nil, businessID, "success", "")

	err := consumer.handleResponse(msg)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if getShipmentCalled {
		t.Error("expected GetShipmentByID not to be called when ShipmentID is nil")
	}
	if updateShipmentCalled {
		t.Error("expected UpdateShipment not to be called when ShipmentID is nil")
	}
	if clearOrderCalled {
		t.Error("expected ClearOrderGuideData not to be called when ShipmentID is nil")
	}
	if cancelledPublished {
		t.Error("expected PublishShipmentCancelled not to be called when ShipmentID is nil")
	}
	if cancelFailedPublished {
		t.Error("expected PublishCancelFailed not to be called on success response")
	}
}

func TestHandleCancelResponse_RepositoryGetError_SkipsUpdateButPublishesSSE(t *testing.T) {
	shipmentID := uint(15)
	businessID := uint(2)

	updateCalled := false
	ssePublished := false

	repoMock := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) {
			return nil, errors.New("db connection lost")
		},
		UpdateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error {
			updateCalled = true
			return nil
		},
	}

	sseMock := &mocks.SSEPublisherMock{
		PublishShipmentCancelledFn: func(ctx context.Context, bID uint, sID uint) {
			ssePublished = true
		},
	}

	consumer := newTestConsumer(repoMock, sseMock)
	msg := buildCancelMessage(shipmentIDPtr(shipmentID), businessID, "success", "")

	err := consumer.handleResponse(msg)

	if err != nil {
		t.Fatalf("expected no error from handleResponse, got %v", err)
	}
	if updateCalled {
		t.Error("expected UpdateShipment not to be called when GetShipmentByID returns error")
	}
	if !ssePublished {
		t.Error("expected PublishShipmentCancelled to be called even when GetShipmentByID fails")
	}
}

func TestHandleCancelResponse_ShipmentWithoutOrderID_SkipsOrderClear(t *testing.T) {
	shipmentID := uint(20)

	clearOrderCalled := false

	repoMock := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) {
			return &domain.Shipment{
				ID:      id,
				OrderID: nil,
				Status:  "pending",
			}, nil
		},
		UpdateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error {
			return nil
		},
		ClearOrderGuideDataFn: func(ctx context.Context, orderID string) error {
			clearOrderCalled = true
			return nil
		},
	}

	sseMock := &mocks.SSEPublisherMock{}

	consumer := newTestConsumer(repoMock, sseMock)
	msg := buildCancelMessage(shipmentIDPtr(shipmentID), 1, "success", "")

	consumer.handleResponse(msg)

	if clearOrderCalled {
		t.Error("expected ClearOrderGuideData not to be called when OrderID is nil")
	}
}

func TestHandleCancelResponse_ShipmentWithEmptyOrderID_SkipsOrderClear(t *testing.T) {
	shipmentID := uint(21)

	clearOrderCalled := false

	repoMock := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) {
			emptyID := ""
			return &domain.Shipment{
				ID:      id,
				OrderID: &emptyID,
				Status:  "pending",
			}, nil
		},
		UpdateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error {
			return nil
		},
		ClearOrderGuideDataFn: func(ctx context.Context, orderID string) error {
			clearOrderCalled = true
			return nil
		},
	}

	sseMock := &mocks.SSEPublisherMock{}

	consumer := newTestConsumer(repoMock, sseMock)
	msg := buildCancelMessage(shipmentIDPtr(shipmentID), 1, "success", "")

	consumer.handleResponse(msg)

	if clearOrderCalled {
		t.Error("expected ClearOrderGuideData not to be called when OrderID is empty string")
	}
}

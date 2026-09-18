package consumer

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/mocks"
)

func buildGenerateSuccessMessage(shipmentID *uint, businessID uint, dataField map[string]interface{}) []byte {
	msg := TransportResponseMessage{
		ShipmentID:    shipmentID,
		BusinessID:    businessID,
		Provider:      "envioclick",
		Operation:     "generate",
		Status:        "success",
		CorrelationID: "corr-generate-cod",
		Timestamp:     time.Now(),
		Data:          dataField,
	}
	b, _ := json.Marshal(msg)
	return b
}

func floatPtr(v float64) *float64 {
	return &v
}

func TestNotificacionGuiaUsaLaComisionDeclaradaNoLaRecalibrada(t *testing.T) {
	shipmentID := uint(9001)
	declaredFee := 5365.0

	repoMock := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) {
			return &domain.Shipment{
				ID:            id,
				CustomerName:  "Cristian Camilo Herrera",
				CustomerPhone: "573000000000",
				OrderNumber:   "15789",
				CodTotal:      floatPtr(73367),
				CodCarrierFee: floatPtr(declaredFee),
			}, nil
		},
		UpdateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error {
			return nil
		},
	}

	var notification *domain.GuideNotificationData
	sseMock := &mocks.SSEPublisherMock{
		PublishGuideGeneratedFn: func(ctx context.Context, businessID uint, shipmentID uint, correlationID string, trackingNumber string, labelURL string, carrier string, n *domain.GuideNotificationData) {
			notification = n
		},
	}

	consumer := newTestConsumer(repoMock, sseMock)

	msg := buildGenerateSuccessMessage(shipmentIDPtr(shipmentID), 46, map[string]interface{}{
		"tracker":       "240061196214",
		"url":           "https://example.com/guide.pdf",
		"carrier":       "INTERRAPIDISIMO",
		"codCarrierFee": 10365.0,
	})

	if err := consumer.handleResponse(msg); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if notification == nil {
		t.Fatal("expected a guide notification to be published")
	}
	if notification.CodCarrierFee == nil {
		t.Fatal("expected CodCarrierFee to be set")
	}
	if *notification.CodCarrierFee != declaredFee {
		t.Errorf("valor a recaudar usa la comision recalibrada (%v) en vez de la declarada al transportador (%v)",
			*notification.CodCarrierFee, declaredFee)
	}
}

package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type SSEPublisherMock struct {
	PublishQuoteReceivedFn    func(ctx context.Context, businessID uint, correlationID string, data map[string]interface{})
	PublishQuoteFailedFn      func(ctx context.Context, businessID uint, correlationID string, errorMsg string)
	PublishGuideGeneratedFn   func(ctx context.Context, businessID uint, shipmentID uint, correlationID string, trackingNumber string, labelURL string, carrier string, notification *domain.GuideNotificationData)
	PublishGuideFailedFn      func(ctx context.Context, businessID uint, shipmentID uint, correlationID string, errorMsg string)
	PublishTrackingUpdatedFn  func(ctx context.Context, businessID uint, correlationID string, data map[string]interface{})
	PublishTrackingFailedFn   func(ctx context.Context, businessID uint, correlationID string, errorMsg string)
	PublishShipmentCancelledFn func(ctx context.Context, businessID uint, shipmentID uint)
	PublishCancelFailedFn     func(ctx context.Context, businessID uint, shipmentID uint, correlationID string, errorMsg string)
}

func (m *SSEPublisherMock) PublishQuoteReceived(ctx context.Context, businessID uint, correlationID string, data map[string]interface{}) {
	if m.PublishQuoteReceivedFn != nil {
		m.PublishQuoteReceivedFn(ctx, businessID, correlationID, data)
	}
}

func (m *SSEPublisherMock) PublishQuoteFailed(ctx context.Context, businessID uint, correlationID string, errorMsg string) {
	if m.PublishQuoteFailedFn != nil {
		m.PublishQuoteFailedFn(ctx, businessID, correlationID, errorMsg)
	}
}

func (m *SSEPublisherMock) PublishGuideGenerated(ctx context.Context, businessID uint, shipmentID uint, correlationID string, trackingNumber string, labelURL string, carrier string, notification *domain.GuideNotificationData) {
	if m.PublishGuideGeneratedFn != nil {
		m.PublishGuideGeneratedFn(ctx, businessID, shipmentID, correlationID, trackingNumber, labelURL, carrier, notification)
	}
}

func (m *SSEPublisherMock) PublishGuideFailed(ctx context.Context, businessID uint, shipmentID uint, correlationID string, errorMsg string) {
	if m.PublishGuideFailedFn != nil {
		m.PublishGuideFailedFn(ctx, businessID, shipmentID, correlationID, errorMsg)
	}
}

func (m *SSEPublisherMock) PublishTrackingUpdated(ctx context.Context, businessID uint, correlationID string, data map[string]interface{}) {
	if m.PublishTrackingUpdatedFn != nil {
		m.PublishTrackingUpdatedFn(ctx, businessID, correlationID, data)
	}
}

func (m *SSEPublisherMock) PublishTrackingFailed(ctx context.Context, businessID uint, correlationID string, errorMsg string) {
	if m.PublishTrackingFailedFn != nil {
		m.PublishTrackingFailedFn(ctx, businessID, correlationID, errorMsg)
	}
}

func (m *SSEPublisherMock) PublishShipmentCancelled(ctx context.Context, businessID uint, shipmentID uint) {
	if m.PublishShipmentCancelledFn != nil {
		m.PublishShipmentCancelledFn(ctx, businessID, shipmentID)
	}
}

func (m *SSEPublisherMock) PublishCancelFailed(ctx context.Context, businessID uint, shipmentID uint, correlationID string, errorMsg string) {
	if m.PublishCancelFailedFn != nil {
		m.PublishCancelFailedFn(ctx, businessID, shipmentID, correlationID, errorMsg)
	}
}

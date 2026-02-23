package mocks

import "context"

// EventPublisherMock implementa ports.IEventPublisher para tests unitarios
type EventPublisherMock struct {
	PublishOrderConfirmedFn    func(ctx context.Context, orderNumber, phoneNumber string, businessID uint) error
	PublishOrderCancelledFn    func(ctx context.Context, orderNumber, reason, phoneNumber string, businessID uint) error
	PublishNoveltyRequestedFn  func(ctx context.Context, orderNumber, noveltyType, phoneNumber string, businessID uint) error
	PublishHandoffRequestedFn  func(ctx context.Context, orderNumber, phoneNumber string, businessID uint, conversationID string) error
}

func (m *EventPublisherMock) PublishOrderConfirmed(ctx context.Context, orderNumber, phoneNumber string, businessID uint) error {
	if m.PublishOrderConfirmedFn != nil {
		return m.PublishOrderConfirmedFn(ctx, orderNumber, phoneNumber, businessID)
	}
	return nil
}

func (m *EventPublisherMock) PublishOrderCancelled(ctx context.Context, orderNumber, reason, phoneNumber string, businessID uint) error {
	if m.PublishOrderCancelledFn != nil {
		return m.PublishOrderCancelledFn(ctx, orderNumber, reason, phoneNumber, businessID)
	}
	return nil
}

func (m *EventPublisherMock) PublishNoveltyRequested(ctx context.Context, orderNumber, noveltyType, phoneNumber string, businessID uint) error {
	if m.PublishNoveltyRequestedFn != nil {
		return m.PublishNoveltyRequestedFn(ctx, orderNumber, noveltyType, phoneNumber, businessID)
	}
	return nil
}

func (m *EventPublisherMock) PublishHandoffRequested(ctx context.Context, orderNumber, phoneNumber string, businessID uint, conversationID string) error {
	if m.PublishHandoffRequestedFn != nil {
		return m.PublishHandoffRequestedFn(ctx, orderNumber, phoneNumber, businessID, conversationID)
	}
	return nil
}

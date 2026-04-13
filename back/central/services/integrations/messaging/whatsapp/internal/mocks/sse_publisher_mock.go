package mocks

import "context"

type SSEEventPublisherMock struct {
	PublishMessageReceivedFn      func(ctx context.Context, businessID uint, conversationID, phoneNumber, messageID, content string) error
	PublishConversationStartedFn  func(ctx context.Context, businessID uint, conversationID, phoneNumber string) error
	PublishMessageStatusUpdatedFn func(ctx context.Context, businessID uint, messageID, status string) error
}

func (m *SSEEventPublisherMock) PublishMessageReceived(ctx context.Context, businessID uint, conversationID, phoneNumber, messageID, content string) error {
	if m.PublishMessageReceivedFn != nil {
		return m.PublishMessageReceivedFn(ctx, businessID, conversationID, phoneNumber, messageID, content)
	}
	return nil
}

func (m *SSEEventPublisherMock) PublishConversationStarted(ctx context.Context, businessID uint, conversationID, phoneNumber string) error {
	if m.PublishConversationStartedFn != nil {
		return m.PublishConversationStartedFn(ctx, businessID, conversationID, phoneNumber)
	}
	return nil
}

func (m *SSEEventPublisherMock) PublishMessageStatusUpdated(ctx context.Context, businessID uint, messageID, status string) error {
	if m.PublishMessageStatusUpdatedFn != nil {
		return m.PublishMessageStatusUpdatedFn(ctx, businessID, messageID, status)
	}
	return nil
}

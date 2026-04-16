package whatsapp_conversation_consumer

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/mocks"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func TestStart_DeclaresQueueBeforeConsuming(t *testing.T) {
	declaredQueue := ""
	consumedQueue := ""
	callOrder := []string{}

	rabbitMock := &mocks.RabbitMQMock{
		DeclareQueueFn: func(queueName string, durable bool) error {
			declaredQueue = queueName
			callOrder = append(callOrder, "declare")
			if !durable {
				t.Error("expected durable=true")
			}
			return nil
		},
		ConsumeFn: func(ctx context.Context, queueName string, handler func([]byte) error) error {
			consumedQueue = queueName
			callOrder = append(callOrder, "consume")
			return nil
		},
	}

	consumer := New(rabbitMock, &mocks.WhatsAppPersisterMock{}, mocks.NewLoggerMock())

	err := consumer.Start(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if declaredQueue != rabbitmq.QueueWhatsAppConversationEvents {
		t.Errorf("expected declared queue %s, got %s", rabbitmq.QueueWhatsAppConversationEvents, declaredQueue)
	}
	if consumedQueue != rabbitmq.QueueWhatsAppConversationEvents {
		t.Errorf("expected consumed queue %s, got %s", rabbitmq.QueueWhatsAppConversationEvents, consumedQueue)
	}
	if len(callOrder) != 2 || callOrder[0] != "declare" || callOrder[1] != "consume" {
		t.Errorf("expected [declare, consume], got %v", callOrder)
	}
}

func TestStart_DeclareQueueError(t *testing.T) {
	expectedErr := errors.New("rabbitmq connection refused")

	rabbitMock := &mocks.RabbitMQMock{
		DeclareQueueFn: func(queueName string, durable bool) error {
			return expectedErr
		},
	}

	consumer := New(rabbitMock, &mocks.WhatsAppPersisterMock{}, mocks.NewLoggerMock())

	err := consumer.Start(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestStart_ConsumeError(t *testing.T) {
	expectedErr := errors.New("channel closed")

	rabbitMock := &mocks.RabbitMQMock{
		DeclareQueueFn: func(queueName string, durable bool) error {
			return nil
		},
		ConsumeFn: func(ctx context.Context, queueName string, handler func([]byte) error) error {
			return expectedErr
		},
	}

	consumer := New(rabbitMock, &mocks.WhatsAppPersisterMock{}, mocks.NewLoggerMock())

	err := consumer.Start(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

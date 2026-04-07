package publisher

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/mocks"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func sampleResult() *entities.DeliveryResult {
	return &entities.DeliveryResult{
		Channel:       "email",
		BusinessID:    1,
		IntegrationID: 5,
		ConfigID:      10,
		To:            "user@example.com",
		Subject:       "Notificación: order.created",
		EventType:     "order.created",
		Status:        "sent",
		SentAt:        time.Date(2026, 3, 2, 10, 0, 0, 0, time.UTC),
	}
}


func TestPublishResult_Success(t *testing.T) {
	rmq := &mocks.RabbitMQMock{}
	pub := New(rmq, mocks.NewLoggerMock())

	err := pub.PublishResult(context.Background(), sampleResult())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(rmq.Published) != 1 {
		t.Fatalf("expected 1 published message, got %d", len(rmq.Published))
	}

	msg := rmq.Published[0]
	if msg.Queue != rabbitmq.QueueNotificationDeliveryResults {
		t.Errorf("expected queue=%s, got %s", rabbitmq.QueueNotificationDeliveryResults, msg.Queue)
	}

	// Verify JSON content
	var decoded map[string]interface{}
	if err := json.Unmarshal(msg.Body, &decoded); err != nil {
		t.Fatalf("published body is not valid JSON: %v", err)
	}
	if decoded["Channel"] != "email" {
		t.Errorf("expected Channel=email, got %v", decoded["Channel"])
	}
	if decoded["Status"] != "sent" {
		t.Errorf("expected Status=sent, got %v", decoded["Status"])
	}
	if decoded["To"] != "user@example.com" {
		t.Errorf("expected To=user@example.com, got %v", decoded["To"])
	}
}

func TestPublishResult_RabbitMQError_ReturnsError(t *testing.T) {
	rmq := &mocks.RabbitMQMock{
		PublishFn: func(ctx context.Context, queue string, body []byte) error {
			return errors.New("connection refused")
		},
	}
	pub := New(rmq, mocks.NewLoggerMock())

	err := pub.PublishResult(context.Background(), sampleResult())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, err) {
		t.Errorf("expected wrapped error, got %v", err)
	}
}

func TestPublishResult_FailedStatus_PublishesCorrectly(t *testing.T) {
	rmq := &mocks.RabbitMQMock{}
	pub := New(rmq, mocks.NewLoggerMock())

	result := sampleResult()
	result.Status = "failed"
	result.ErrorMessage = "SES timeout"

	err := pub.PublishResult(context.Background(), result)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	var decoded map[string]interface{}
	json.Unmarshal(rmq.Published[0].Body, &decoded)

	if decoded["Status"] != "failed" {
		t.Errorf("expected Status=failed, got %v", decoded["Status"])
	}
	if decoded["ErrorMessage"] != "SES timeout" {
		t.Errorf("expected ErrorMessage='SES timeout', got %v", decoded["ErrorMessage"])
	}
}

func TestPublishResult_SerializesAllFields(t *testing.T) {
	rmq := &mocks.RabbitMQMock{}
	pub := New(rmq, mocks.NewLoggerMock())

	result := &entities.DeliveryResult{
		Channel:       "email",
		BusinessID:    42,
		IntegrationID: 7,
		ConfigID:      99,
		To:            "test@test.com",
		Subject:       "Test Subject",
		EventType:     "order.shipped",
		Status:        "sent",
		SentAt:        time.Date(2026, 3, 2, 15, 30, 0, 0, time.UTC),
	}

	_ = pub.PublishResult(context.Background(), result)

	var decoded entities.DeliveryResult
	json.Unmarshal(rmq.Published[0].Body, &decoded)

	if decoded.Channel != "email" {
		t.Errorf("Channel: got %s, want email", decoded.Channel)
	}
	if decoded.BusinessID != 42 {
		t.Errorf("BusinessID: got %d, want 42", decoded.BusinessID)
	}
	if decoded.IntegrationID != 7 {
		t.Errorf("IntegrationID: got %d, want 7", decoded.IntegrationID)
	}
	if decoded.ConfigID != 99 {
		t.Errorf("ConfigID: got %d, want 99", decoded.ConfigID)
	}
	if decoded.To != "test@test.com" {
		t.Errorf("To: got %s, want test@test.com", decoded.To)
	}
	if decoded.Subject != "Test Subject" {
		t.Errorf("Subject: got %s, want Test Subject", decoded.Subject)
	}
	if decoded.EventType != "order.shipped" {
		t.Errorf("EventType: got %s, want order.shipped", decoded.EventType)
	}
}

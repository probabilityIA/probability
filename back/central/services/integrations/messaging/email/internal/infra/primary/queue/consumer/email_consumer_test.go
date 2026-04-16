package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/infra/primary/queue/consumer/request"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/mocks"
)

func newTestConsumer(uc *mocks.UseCaseMock) *emailConsumer {
	return &emailConsumer{
		rabbitMQ: &mocks.RabbitMQMock{},
		useCase:  uc,
		logger:   mocks.NewLoggerMock(),
	}
}

func makeEventJSON(event request.EmailNotificationEvent) []byte {
	data, _ := json.Marshal(event)
	return data
}


func TestHandleMessage_ValidEvent_CallsUseCase(t *testing.T) {
	uc := &mocks.UseCaseMock{}
	c := newTestConsumer(uc)

	event := request.EmailNotificationEvent{
		EventType:     "order.created",
		BusinessID:    1,
		IntegrationID: 5,
		ConfigID:      10,
		CustomerEmail: "user@example.com",
		EventData:     map[string]interface{}{"key": "value"},
	}

	err := c.handleMessage(makeEventJSON(event))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(uc.Calls) != 1 {
		t.Fatalf("expected 1 usecase call, got %d", len(uc.Calls))
	}

	call := uc.Calls[0]
	if call.CustomerEmail != "user@example.com" {
		t.Errorf("expected email=user@example.com, got %s", call.CustomerEmail)
	}
	if call.EventType != "order.created" {
		t.Errorf("expected event_type=order.created, got %s", call.EventType)
	}
	if call.BusinessID != 1 {
		t.Errorf("expected business_id=1, got %d", call.BusinessID)
	}
	if call.IntegrationID != 5 {
		t.Errorf("expected integration_id=5, got %d", call.IntegrationID)
	}
	if call.ConfigID != 10 {
		t.Errorf("expected config_id=10, got %d", call.ConfigID)
	}
}

func TestHandleMessage_MalformedJSON_ReturnsNil(t *testing.T) {
	uc := &mocks.UseCaseMock{}
	c := newTestConsumer(uc)

	err := c.handleMessage([]byte("not json"))
	if err != nil {
		t.Fatalf("expected nil (no retry for malformed), got %v", err)
	}

	if len(uc.Calls) != 0 {
		t.Errorf("expected 0 usecase calls, got %d", len(uc.Calls))
	}
}

func TestHandleMessage_EmptyCustomerEmail_Discarded(t *testing.T) {
	uc := &mocks.UseCaseMock{}
	c := newTestConsumer(uc)

	event := request.EmailNotificationEvent{
		EventType:     "order.created",
		BusinessID:    1,
		CustomerEmail: "",
	}

	err := c.handleMessage(makeEventJSON(event))
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if len(uc.Calls) != 0 {
		t.Errorf("expected 0 usecase calls (message discarded), got %d", len(uc.Calls))
	}
}

func TestHandleMessage_UseCaseError_StillReturnsNil(t *testing.T) {
	uc := &mocks.UseCaseMock{
		SendNotificationEmailFn: func(ctx context.Context, dto dtos.SendEmailDTO) error {
			return errors.New("send failed")
		},
	}
	c := newTestConsumer(uc)

	event := request.EmailNotificationEvent{
		EventType:     "order.created",
		BusinessID:    1,
		CustomerEmail: "user@example.com",
	}

	err := c.handleMessage(makeEventJSON(event))
	if err != nil {
		t.Fatalf("expected nil (fire-and-forget), got %v", err)
	}
}

func TestHandleMessage_EmptyBody_ReturnsNil(t *testing.T) {
	uc := &mocks.UseCaseMock{}
	c := newTestConsumer(uc)

	err := c.handleMessage([]byte(""))
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if len(uc.Calls) != 0 {
		t.Errorf("expected 0 usecase calls, got %d", len(uc.Calls))
	}
}

func TestHandleMessage_MapsAllFieldsToDTO(t *testing.T) {
	uc := &mocks.UseCaseMock{}
	c := newTestConsumer(uc)

	event := request.EmailNotificationEvent{
		EventType:     "order.shipped",
		BusinessID:    42,
		IntegrationID: 7,
		ConfigID:      99,
		CustomerEmail: "test@test.com",
		EventData:     map[string]interface{}{"tracking": "TR123"},
	}

	_ = c.handleMessage(makeEventJSON(event))

	call := uc.Calls[0]
	if call.EventType != "order.shipped" {
		t.Errorf("event_type: got %s, want order.shipped", call.EventType)
	}
	if call.BusinessID != 42 {
		t.Errorf("business_id: got %d, want 42", call.BusinessID)
	}
	if call.IntegrationID != 7 {
		t.Errorf("integration_id: got %d, want 7", call.IntegrationID)
	}
	if call.ConfigID != 99 {
		t.Errorf("config_id: got %d, want 99", call.ConfigID)
	}
	if call.CustomerEmail != "test@test.com" {
		t.Errorf("customer_email: got %s, want test@test.com", call.CustomerEmail)
	}
}

package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/mocks"
)

func newTestUseCase(client *mocks.EmailClientMock, pub *mocks.ResultPublisherMock) *useCase {
	return &useCase{
		emailClient: client,
		resultPub:   pub,
		logger:      mocks.NewLoggerMock(),
	}
}

func validDTO() dtos.SendEmailDTO {
	return dtos.SendEmailDTO{
		EventType:     "order.created",
		BusinessID:    1,
		IntegrationID: 5,
		ConfigID:      10,
		CustomerEmail: "user@example.com",
		EventData:     map[string]interface{}{"order_id": "ABC123"},
	}
}


func TestSendNotificationEmail_Success(t *testing.T) {
	client := &mocks.EmailClientMock{}
	pub := &mocks.ResultPublisherMock{}
	uc := newTestUseCase(client, pub)

	err := uc.SendNotificationEmail(context.Background(), validDTO())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// Verifica que se envió el email
	if len(client.Calls) != 1 {
		t.Fatalf("expected 1 client call, got %d", len(client.Calls))
	}
	call := client.Calls[0]
	if call.To != "user@example.com" {
		t.Errorf("expected to=user@example.com, got %s", call.To)
	}
	if call.Subject != "Notificación: order.created" {
		t.Errorf("expected subject 'Notificación: order.created', got %s", call.Subject)
	}

	// Verifica que se publicó resultado con status "sent"
	if len(pub.Results) != 1 {
		t.Fatalf("expected 1 published result, got %d", len(pub.Results))
	}
	result := pub.Results[0]
	if result.Status != "sent" {
		t.Errorf("expected status=sent, got %s", result.Status)
	}
	if result.Channel != "email" {
		t.Errorf("expected channel=email, got %s", result.Channel)
	}
	if result.To != "user@example.com" {
		t.Errorf("expected to=user@example.com, got %s", result.To)
	}
	if result.BusinessID != 1 {
		t.Errorf("expected business_id=1, got %d", result.BusinessID)
	}
	if result.ConfigID != 10 {
		t.Errorf("expected config_id=10, got %d", result.ConfigID)
	}
	if result.ErrorMessage != "" {
		t.Errorf("expected empty error message, got %s", result.ErrorMessage)
	}
}

func TestSendNotificationEmail_EmptyEmail_ReturnsErrMissingRecipient(t *testing.T) {
	client := &mocks.EmailClientMock{}
	pub := &mocks.ResultPublisherMock{}
	uc := newTestUseCase(client, pub)

	dto := validDTO()
	dto.CustomerEmail = ""

	err := uc.SendNotificationEmail(context.Background(), dto)
	if !errors.Is(err, domainerrors.ErrMissingRecipient) {
		t.Fatalf("expected ErrMissingRecipient, got %v", err)
	}

	// No se envió email ni se publicó resultado
	if len(client.Calls) != 0 {
		t.Errorf("expected 0 client calls, got %d", len(client.Calls))
	}
	if len(pub.Results) != 0 {
		t.Errorf("expected 0 published results, got %d", len(pub.Results))
	}
}

func TestSendNotificationEmail_ClientError_ReturnsError_PublishesFailedResult(t *testing.T) {
	sendErr := errors.New("SES connection failed")
	client := &mocks.EmailClientMock{
		SendHTMLFn: func(ctx context.Context, to, subject, html string) error {
			return sendErr
		},
	}
	pub := &mocks.ResultPublisherMock{}
	uc := newTestUseCase(client, pub)

	err := uc.SendNotificationEmail(context.Background(), validDTO())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "SES connection failed" {
		t.Errorf("expected 'SES connection failed', got %s", err.Error())
	}

	// Verifica que se publicó resultado con status "failed"
	if len(pub.Results) != 1 {
		t.Fatalf("expected 1 published result, got %d", len(pub.Results))
	}
	result := pub.Results[0]
	if result.Status != "failed" {
		t.Errorf("expected status=failed, got %s", result.Status)
	}
	if result.ErrorMessage != "SES connection failed" {
		t.Errorf("expected error message 'SES connection failed', got %s", result.ErrorMessage)
	}
}

func TestSendNotificationEmail_PublisherError_DoesNotFailSend(t *testing.T) {
	client := &mocks.EmailClientMock{}
	pub := &mocks.ResultPublisherMock{
		PublishResultFn: func(ctx context.Context, result *entities.DeliveryResult) error {
			return errors.New("rabbitmq down")
		},
	}
	uc := newTestUseCase(client, pub)

	// El envío debe ser exitoso aunque el publisher falle (best-effort)
	err := uc.SendNotificationEmail(context.Background(), validDTO())
	if err != nil {
		t.Fatalf("expected nil error (publisher failure is best-effort), got %v", err)
	}

	// El email sí se envió
	if len(client.Calls) != 1 {
		t.Errorf("expected 1 client call, got %d", len(client.Calls))
	}
}

func TestSendNotificationEmail_BothClientAndPublisherFail_ReturnsClientError(t *testing.T) {
	client := &mocks.EmailClientMock{
		SendHTMLFn: func(ctx context.Context, to, subject, html string) error {
			return errors.New("SES error")
		},
	}
	pub := &mocks.ResultPublisherMock{
		PublishResultFn: func(ctx context.Context, result *entities.DeliveryResult) error {
			return errors.New("rabbitmq error")
		},
	}
	uc := newTestUseCase(client, pub)

	err := uc.SendNotificationEmail(context.Background(), validDTO())
	if err == nil || err.Error() != "SES error" {
		t.Fatalf("expected 'SES error', got %v", err)
	}
}

func TestSendNotificationEmail_SetsCorrectDeliveryResultFields(t *testing.T) {
	client := &mocks.EmailClientMock{}
	pub := &mocks.ResultPublisherMock{}
	uc := newTestUseCase(client, pub)

	dto := dtos.SendEmailDTO{
		EventType:     "order.shipped",
		BusinessID:    42,
		IntegrationID: 7,
		ConfigID:      99,
		CustomerEmail: "test@test.com",
		EventData:     map[string]interface{}{"tracking": "TR123"},
	}

	_ = uc.SendNotificationEmail(context.Background(), dto)

	result := pub.Results[0]
	if result.Channel != "email" {
		t.Errorf("expected channel=email, got %s", result.Channel)
	}
	if result.BusinessID != 42 {
		t.Errorf("expected business_id=42, got %d", result.BusinessID)
	}
	if result.IntegrationID != 7 {
		t.Errorf("expected integration_id=7, got %d", result.IntegrationID)
	}
	if result.ConfigID != 99 {
		t.Errorf("expected config_id=99, got %d", result.ConfigID)
	}
	if result.To != "test@test.com" {
		t.Errorf("expected to=test@test.com, got %s", result.To)
	}
	if result.EventType != "order.shipped" {
		t.Errorf("expected event_type=order.shipped, got %s", result.EventType)
	}
	if result.Subject != "Notificación: order.shipped" {
		t.Errorf("expected subject, got %s", result.Subject)
	}
	if result.SentAt.IsZero() {
		t.Error("expected non-zero SentAt")
	}
}

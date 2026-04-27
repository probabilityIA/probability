package app

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	boldErrors "github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type webhookUseCase struct {
	repo      ports.IIntegrationRepository
	publisher ports.IWebhookPublisher
	log       log.ILogger
}

func NewWebhookUseCase(repo ports.IIntegrationRepository, publisher ports.IWebhookPublisher, logger log.ILogger) ports.IWebhookUseCase {
	return &webhookUseCase{
		repo:      repo,
		publisher: publisher,
		log:       logger.WithModule("bold.webhook.usecase"),
	}
}

type cloudEventsEnvelope struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Subject string `json:"subject"`
	Source  string `json:"source"`
	Time    int64  `json:"time"`
	Data    struct {
		PaymentID         string  `json:"payment_id"`
		Amount            float64 `json:"amount"`
		Currency          string  `json:"currency"`
		PaymentMethod     string  `json:"payment_method"`
		MerchantReference string  `json:"merchant_reference"`
		PayerEmail        string  `json:"payer_email"`
	} `json:"data"`
}

func (uc *webhookUseCase) HandleIncomingWebhook(ctx context.Context, signatureHeader string, body []byte) error {
	if len(body) == 0 {
		return fmt.Errorf("empty body")
	}

	cfg, err := uc.repo.GetBoldConfig(ctx)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("bold webhook: cannot load config")
		return err
	}

	if !verifySignature(body, signatureHeader, cfg) {
		uc.log.Warn(ctx).
			Str("signature_header", signatureHeader).
			Msg("bold webhook: invalid signature")
		return boldErrors.ErrInvalidSignature
	}

	var envelope cloudEventsEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return fmt.Errorf("decode bold webhook envelope: %w", err)
	}
	if envelope.ID == "" {
		return fmt.Errorf("bold webhook missing event id")
	}

	var occurredAt *time.Time
	if envelope.Time > 0 {
		t := time.Unix(0, envelope.Time)
		occurredAt = &t
	}

	msg := &ports.BoldWebhookMessage{
		BoldEventID:       envelope.ID,
		Type:              envelope.Type,
		Subject:           envelope.Subject,
		Source:            envelope.Source,
		OccurredAt:        occurredAt,
		PaymentID:         envelope.Data.PaymentID,
		MerchantReference: envelope.Data.MerchantReference,
		Amount:            envelope.Data.Amount,
		Currency:          envelope.Data.Currency,
		PaymentMethod:     envelope.Data.PaymentMethod,
		PayerEmail:        envelope.Data.PayerEmail,
		RawPayload:        body,
	}

	if err := uc.publisher.PublishWebhookEvent(ctx, msg); err != nil {
		return fmt.Errorf("publish bold webhook event: %w", err)
	}

	uc.log.Info(ctx).
		Str("bold_event_id", envelope.ID).
		Str("type", envelope.Type).
		Msg("bold webhook validated and published")

	return nil
}

func verifySignature(body []byte, signatureHeader string, cfg *ports.BoldConfig) bool {
	if signatureHeader == "" {
		return false
	}

	secret, ok := cfg.SecretKey()
	if !ok {
		return false
	}

	bodyB64 := base64.StdEncoding.EncodeToString(body)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(bodyB64))
	expected := hex.EncodeToString(mac.Sum(nil))

	provided := strings.TrimSpace(signatureHeader)
	return hmac.Equal([]byte(expected), []byte(provided))
}

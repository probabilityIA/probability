package app

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	bcErrors "github.com/secamc93/probability/back/central/services/integrations/pay/bancolombia/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bancolombia/internal/domain/ports"
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
		log:       logger.WithModule("bancolombia.webhook.usecase"),
	}
}

func (uc *webhookUseCase) HandleIncomingWebhook(ctx context.Context, signatureHeader string, body []byte, isTest bool) error {
	if len(body) == 0 {
		return fmt.Errorf("empty body")
	}

	cfg, err := uc.repo.GetConfigForMode(ctx, isTest)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("bancolombia webhook: cannot load config")
		return err
	}

	signatureValid := false
	secret, hasSecret := cfg.WebhookSecretKey()
	if hasSecret {
		signatureValid = verifySignature(body, signatureHeader, secret)
		if !signatureValid {
			uc.log.Warn(ctx).
				Str("environment", cfg.Environment).
				Bool("is_test", isTest).
				Int("body_len", len(body)).
				Msg("bancolombia webhook: invalid signature")
			return bcErrors.ErrInvalidSignature
		}
	} else {
		uc.log.Warn(ctx).
			Str("environment", cfg.Environment).
			Bool("is_test", isTest).
			Msg("bancolombia webhook: no webhook_secret configured, accepting without signature validation")
	}

	msg, err := parseWebhookPayload(body)
	if err != nil {
		return fmt.Errorf("decode bancolombia webhook envelope: %w", err)
	}
	msg.IsTest = isTest
	msg.SignatureValid = signatureValid
	msg.RawPayload = body

	if err := uc.publisher.PublishWebhookEvent(ctx, msg); err != nil {
		return fmt.Errorf("publish bancolombia webhook event: %w", err)
	}

	uc.log.Info(ctx).
		Str("event_id", msg.EventID).
		Str("type", msg.Type).
		Str("state", msg.TransferState).
		Str("reference", msg.Reference).
		Msg("bancolombia webhook validated and published")

	return nil
}

func verifySignature(body []byte, signatureHeader, secret string) bool {
	provided := strings.TrimSpace(signatureHeader)
	if provided == "" {
		return false
	}
	provided = strings.TrimPrefix(provided, "sha256=")

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	sum := mac.Sum(nil)
	if hmac.Equal([]byte(hex.EncodeToString(sum)), []byte(provided)) {
		return true
	}
	if hmac.Equal([]byte(base64.StdEncoding.EncodeToString(sum)), []byte(provided)) {
		return true
	}

	macB64 := hmac.New(sha256.New, []byte(secret))
	macB64.Write([]byte(base64.StdEncoding.EncodeToString(body)))
	sumB64 := macB64.Sum(nil)
	if hmac.Equal([]byte(hex.EncodeToString(sumB64)), []byte(provided)) {
		return true
	}
	return hmac.Equal([]byte(base64.StdEncoding.EncodeToString(sumB64)), []byte(provided))
}

func parseWebhookPayload(body []byte) (*ports.BancolombiaWebhookMessage, error) {
	var root map[string]interface{}
	if err := json.Unmarshal(body, &root); err != nil {
		return nil, err
	}

	flat := map[string]interface{}{}
	mergeMap(flat, root)
	for _, key := range []string{"header", "data", "transaction", "transfer", "payload"} {
		switch nested := root[key].(type) {
		case map[string]interface{}:
			mergeMap(flat, nested)
		case []interface{}:
			if len(nested) > 0 {
				if first, ok := nested[0].(map[string]interface{}); ok {
					mergeMap(flat, first)
					for _, innerKey := range []string{"transaction", "transfer", "merchant"} {
						if inner, ok := first[innerKey].(map[string]interface{}); ok {
							mergeMap(flat, inner)
						}
					}
				}
			}
		}
	}

	msg := &ports.BancolombiaWebhookMessage{
		EventID:       pickString(flat, "eventId", "event_id", "messageId", "message_id", "id", "notificationId"),
		Type:          pickString(flat, "eventType", "event_type", "type", "event", "notificationType"),
		TransferState: pickString(flat, "transferState", "transfer_state", "status", "state", "transactionState"),
		Reference:     pickString(flat, "transferReference", "paymentReference", "payment_reference", "reference", "transactionReference", "commerceReference", "invoiceNumber"),
		TransferCode:  pickString(flat, "transferCode", "transfer_code", "voucher", "cus", "authorizationCode", "approvalNumber", "transactionId"),
		Amount:        pickFloat(flat, "transferAmount", "transfer_amount", "amount", "value", "totalAmount"),
		Currency:      pickString(flat, "currency", "currencyCode"),
		PayerAccount:  pickString(flat, "originAccount", "payerAccount", "sourceAccount"),
		OccurredAt:    pickTime(flat, "transferDate", "transfer_date", "date", "occurredAt", "timestamp", "time"),
	}

	if msg.EventID == "" {
		sum := sha256.Sum256(body)
		msg.EventID = "sha256:" + hex.EncodeToString(sum[:])
	}
	if msg.Type == "" {
		msg.Type = "qr.transfer.notification"
	}
	return msg, nil
}

func mergeMap(dst, src map[string]interface{}) {
	for k, v := range src {
		switch v.(type) {
		case map[string]interface{}, []interface{}:
			continue
		default:
			if _, exists := dst[k]; !exists {
				dst[k] = v
			}
		}
	}
}

func pickString(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			switch typed := v.(type) {
			case string:
				if strings.TrimSpace(typed) != "" {
					return strings.TrimSpace(typed)
				}
			case float64:
				return strconv.FormatFloat(typed, 'f', -1, 64)
			}
		}
	}
	return ""
}

func pickFloat(m map[string]interface{}, keys ...string) float64 {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			switch typed := v.(type) {
			case float64:
				return typed
			case string:
				if f, err := strconv.ParseFloat(strings.TrimSpace(typed), 64); err == nil {
					return f
				}
			}
		}
	}
	return 0
}

func pickTime(m map[string]interface{}, keys ...string) *time.Time {
	raw := pickString(m, keys...)
	if raw == "" {
		return nil
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"} {
		if t, err := time.Parse(layout, raw); err == nil {
			return &t
		}
	}
	if unix, err := strconv.ParseInt(raw, 10, 64); err == nil && unix > 0 {
		var t time.Time
		if unix > 1_000_000_000_000 {
			t = time.UnixMilli(unix)
		} else {
			t = time.Unix(unix, 0)
		}
		return &t
	}
	return nil
}

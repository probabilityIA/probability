package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
	wooqueue "github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/infra/primary/queue"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func (h *wooCommerceHandler) HandleWebhook(c *gin.Context) {
	ctx := c.Request.Context()

	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Failed to read webhook body")
		c.Status(http.StatusBadRequest)
		return
	}

	topic := c.GetHeader("X-WC-Webhook-Topic")
	source := c.GetHeader("X-WC-Webhook-Source")
	signature := c.GetHeader("X-WC-Webhook-Signature")
	integrationID := c.Query("integration_id")

	h.logger.Info(ctx).
		Str("topic", topic).
		Str("source", source).
		Str("integration_id", integrationID).
		Int("body_size", len(rawBody)).
		Msg("WooCommerce webhook received")

	webhookSecret := os.Getenv("WOOCOMMERCE_WEBHOOK_SECRET")
	signatureChecked := webhookSecret != "" && signature != ""
	signatureOK := true
	if signatureChecked {
		signatureOK = verifyWebhookHMAC(rawBody, signature, webhookSecret)
	}

	if !signatureOK {
		h.logger.Warn(ctx).
			Str("topic", topic).
			Str("source", source).
			Msg("WooCommerce webhook invalid HMAC signature")
		h.saveWebhookLog(topic, integrationID, c, rawBody, domain.WebhookLogStatusRejected,
			http.StatusUnauthorized, strPtr("Firma HMAC invalida"), signatureChecked, signatureOK)
		c.Status(http.StatusUnauthorized)
		return
	}

	c.Status(http.StatusOK)

	logStatus := domain.WebhookLogStatusReceived
	var logError *string
	if topic == "" || len(rawBody) == 0 {
		logStatus = domain.WebhookLogStatusIgnored
		logError = strPtr("Webhook sin topic o sin cuerpo")
	}
	h.saveWebhookLog(topic, integrationID, c, rawBody, logStatus, http.StatusOK, logError, signatureChecked, signatureOK)

	if topic != "" && len(rawBody) > 0 {
		h.enqueueWebhook(topic, source, integrationID, rawBody)
	}
}

func strPtr(s string) *string { return &s }

func (h *wooCommerceHandler) saveWebhookLog(
	topic string,
	integrationID string,
	c *gin.Context,
	rawBody []byte,
	status string,
	responseCode int,
	errorMessage *string,
	signatureChecked bool,
	signatureOK bool,
) {
	if h.webhookLogRepo == nil {
		return
	}

	headers := map[string]string{}
	for k := range c.Request.Header {
		headers[k] = c.GetHeader(k)
	}
	headersJSON, _ := json.Marshal(headers)

	entry := &domain.WebhookLog{
		EventType:     topic,
		URL:           c.Request.URL.String(),
		Method:        c.Request.Method,
		Headers:       headersJSON,
		RequestBody:   rawBody,
		RemoteIP:      c.ClientIP(),
		Status:        status,
		ResponseCode:  responseCode,
		ErrorMessage:  errorMessage,
		IntegrationID: integrationID,
	}

	if signatureChecked {
		valid := signatureOK
		entry.SignatureValid = &valid
	}

	var payload struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(rawBody, &payload); err == nil && payload.Status != "" {
		raw := payload.Status
		entry.RawStatus = &raw
		mapped := mapper.MapWooStatus(raw)
		entry.MappedStatus = &mapped
	}

	retention := time.Now().AddDate(0, 0, domain.WebhookLogRetentionDays)
	entry.RetentionUntil = &retention

	go func() {
		bgCtx := context.Background()
		if err := h.webhookLogRepo.Save(bgCtx, entry); err != nil {
			h.logger.Warn(bgCtx).Err(err).Str("topic", topic).Msg("No se pudo guardar el log del webhook de WooCommerce")
		}
	}()
}

func (h *wooCommerceHandler) enqueueWebhook(topic, source, integrationID string, rawBody []byte) {
	ctx := context.Background()

	if h.rabbit != nil {
		msg := wooqueue.WebhookMessage{
			Topic:         topic,
			Source:        source,
			IntegrationID: integrationID,
			Body:          rawBody,
		}
		payload, err := json.Marshal(msg)
		if err == nil {
			if err := h.rabbit.Publish(ctx, rabbitmq.QueueWebhooksWoocommerceReceived, payload); err == nil {
				return
			}
			h.logger.Warn(ctx).Str("topic", topic).Str("source", source).
				Msg("No se pudo encolar el webhook, procesando inline como fallback")
		}
	}

	go h.processWebhookAsync(topic, source, integrationID, rawBody)
}

func (h *wooCommerceHandler) processWebhookAsync(topic, source, integrationID string, rawBody []byte) {
	ctx := context.Background()

	if err := h.useCase.ProcessWebhookOrder(ctx, topic, source, integrationID, rawBody); err != nil {
		h.logger.Error(ctx).Err(err).
			Str("topic", topic).
			Str("source", source).
			Str("integration_id", integrationID).
			Msg("Failed to process WooCommerce webhook order")
	}
}

func verifyWebhookHMAC(body []byte, signature string, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedSig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expectedSig), []byte(signature))
}

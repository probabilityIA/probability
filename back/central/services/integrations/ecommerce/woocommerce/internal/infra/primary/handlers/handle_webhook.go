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

	"github.com/gin-gonic/gin"
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
	if webhookSecret != "" && signature != "" {
		if !verifyWebhookHMAC(rawBody, signature, webhookSecret) {
			h.logger.Warn(ctx).
				Str("topic", topic).
				Str("source", source).
				Msg("WooCommerce webhook invalid HMAC signature")
			c.Status(http.StatusUnauthorized)
			return
		}
	}

	c.Status(http.StatusOK)

	if topic != "" && len(rawBody) > 0 {
		h.enqueueWebhook(topic, source, integrationID, rawBody)
	}
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

package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// HandleWebhook recibe webhooks de eventos de WooCommerce.
// WooCommerce envía un POST con:
//   - X-WC-Webhook-Topic: tipo de evento (order.created, order.updated, etc.)
//   - X-WC-Webhook-Source: URL de la tienda origen
//   - X-WC-Webhook-Signature: HMAC-SHA256 del body con el webhook secret
//
// Referencia: https://woocommerce.com/document/webhooks/
func (h *wooCommerceHandler) HandleWebhook(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. Leer body crudo
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Failed to read webhook body")
		c.Status(http.StatusBadRequest)
		return
	}

	// 2. Extraer headers
	topic := c.GetHeader("X-WC-Webhook-Topic")
	source := c.GetHeader("X-WC-Webhook-Source")
	signature := c.GetHeader("X-WC-Webhook-Signature")

	h.logger.Info(ctx).
		Str("topic", topic).
		Str("source", source).
		Int("body_size", len(rawBody)).
		Msg("WooCommerce webhook received")

	// 3. Validar firma HMAC (si hay secret configurado)
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

	// 4. Responder 200 inmediatamente (WooCommerce requiere respuesta rápida)
	c.Status(http.StatusOK)

	// 5. Procesar asincrónicamente
	if topic != "" && len(rawBody) > 0 {
		go h.processWebhookAsync(topic, source, rawBody)
	}
}

func (h *wooCommerceHandler) processWebhookAsync(topic, source string, rawBody []byte) {
	ctx := context.Background()

	if err := h.useCase.ProcessWebhookOrder(ctx, topic, source, rawBody); err != nil {
		h.logger.Error(ctx).Err(err).
			Str("topic", topic).
			Str("source", source).
			Msg("Failed to process WooCommerce webhook order")
	}
}

// verifyWebhookHMAC valida la firma HMAC-SHA256 del webhook de WooCommerce.
// WooCommerce firma: base64(HMAC-SHA256(secret, body))
func verifyWebhookHMAC(body []byte, signature string, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedSig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expectedSig), []byte(signature))
}

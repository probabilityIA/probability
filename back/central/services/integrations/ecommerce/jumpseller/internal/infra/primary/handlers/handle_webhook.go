package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

func (h *jumpsellerHandler) HandleWebhook(c *gin.Context) {
	ctx := c.Request.Context()

	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Failed to read Jumpseller webhook body")
		c.Status(http.StatusBadRequest)
		return
	}

	event := c.GetHeader(domain.HeaderEvent)
	storeCode := c.GetHeader(domain.HeaderStoreCode)
	signature := c.GetHeader(domain.HeaderHmac)
	integrationID := c.Query("integration_id")

	h.logger.Info(ctx).
		Str("event", event).
		Str("store_code", storeCode).
		Str("integration_id", integrationID).
		Int("body_size", len(rawBody)).
		Msg("Jumpseller webhook received")

	if integrationID == "" {
		h.logger.Warn(ctx).Msg("Jumpseller webhook sin integration_id en la URL")
		c.Status(http.StatusBadRequest)
		return
	}

	hooksToken, err := h.useCase.ResolveHooksToken(ctx, integrationID)
	if err != nil {
		h.logger.Error(ctx).Err(err).
			Str("integration_id", integrationID).
			Msg("No se pudo resolver el hooks_token para validar el webhook de Jumpseller")
		c.Status(http.StatusUnauthorized)
		return
	}

	if !verifyWebhookHMAC(rawBody, signature, hooksToken) {
		h.logger.Warn(ctx).
			Str("event", event).
			Str("store_code", storeCode).
			Str("integration_id", integrationID).
			Msg("Jumpseller webhook invalid HMAC signature")
		c.Status(http.StatusUnauthorized)
		return
	}

	c.Status(http.StatusOK)

	if event != "" && len(rawBody) > 0 {
		go h.processWebhookAsync(event, storeCode, integrationID, rawBody)
	}
}

func (h *jumpsellerHandler) processWebhookAsync(event, storeCode, integrationID string, rawBody []byte) {
	ctx := context.Background()

	if err := h.useCase.ProcessWebhookOrder(ctx, event, storeCode, integrationID, rawBody); err != nil {
		h.logger.Error(ctx).Err(err).
			Str("event", event).
			Str("store_code", storeCode).
			Str("integration_id", integrationID).
			Msg("Failed to process Jumpseller webhook order")
	}
}

func verifyWebhookHMAC(body []byte, signature, hooksToken string) bool {
	if hooksToken == "" || signature == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(hooksToken))
	mac.Write(body)
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

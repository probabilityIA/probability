package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/infra/secondary/client/response"
)

// HandleWebhook recibe webhooks de eventos de VTEX.
// VTEX envía un POST con el payload del evento en el body (Hook v1).
// Se responde 200 inmediatamente y se procesa en background.
//
// Referencia: https://developers.vtex.com/docs/guides/orders-feed
func (h *vtexHandler) HandleWebhook(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. Leer body crudo
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Failed to read VTEX webhook body")
		c.Status(http.StatusBadRequest)
		return
	}

	// 2. Deserializar el payload del webhook
	var webhookBody response.VTEXWebhookBody
	if err := json.Unmarshal(rawBody, &webhookBody); err != nil {
		h.logger.Warn(ctx).Err(err).
			Str("raw_body", string(rawBody)).
			Msg("Invalid VTEX webhook payload")
		c.Status(http.StatusBadRequest)
		return
	}

	h.logger.Info(ctx).
		Str("order_id", webhookBody.OrderID).
		Str("state", webhookBody.State).
		Str("last_state", webhookBody.LastState).
		Str("domain", webhookBody.Domain).
		Msg("VTEX webhook received")

	// 3. Responder 200 inmediatamente (VTEX espera respuesta rápida)
	c.Status(http.StatusOK)

	// 4. Convertir a dominio y procesar en background
	payload := webhookBody.ToDomain()
	go h.processWebhookAsync(&payload)
}

func (h *vtexHandler) processWebhookAsync(payload *domain.VTEXWebhookPayload) {
	ctx := context.Background()

	if err := h.useCase.ProcessWebhook(ctx, payload); err != nil {
		h.logger.Error(ctx).Err(err).
			Str("order_id", payload.OrderID).
			Str("state", payload.State).
			Msg("Failed to process VTEX webhook")
	}
}

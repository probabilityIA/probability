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

func (h *vtexHandler) HandleWebhook(c *gin.Context) {
	ctx := c.Request.Context()

	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Failed to read VTEX webhook body")
		c.Status(http.StatusBadRequest)
		return
	}

	integrationID := c.Query("integration_id")
	if integrationID == "" {
		h.logger.Warn(ctx).Msg("VTEX webhook sin integration_id en la query")
		c.Status(http.StatusBadRequest)
		return
	}

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
		Str("integration_id", integrationID).
		Msg("VTEX webhook received")

	c.Status(http.StatusOK)

	payload := webhookBody.ToDomain()
	payload.IntegrationID = integrationID
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

package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/app/usecases"
)

type tiendanubeWebhookPayload struct {
	StoreID json.Number `json:"store_id"`
	Event   string      `json:"event"`
	ID      json.Number `json:"id"`
}

func (h *tiendanubeHandler) HandleWebhook(c *gin.Context) {
	ctx := c.Request.Context()

	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("No se pudo leer el cuerpo del webhook de Tiendanube")
		c.Status(http.StatusOK)
		return
	}

	var payload tiendanubeWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		h.logger.Warn(ctx).Err(err).
			Str("body", truncateBody(body)).
			Msg("Webhook de Tiendanube con payload ilegible, se descarta")
		c.Status(http.StatusOK)
		return
	}

	integrationID := c.Query("integration_id")
	if integrationID == "" {
		h.logger.Warn(ctx).
			Str("event", payload.Event).
			Str("store_id", payload.StoreID.String()).
			Msg("Webhook de Tiendanube sin integration_id en la URL, se descarta")
		c.Status(http.StatusOK)
		return
	}

	if err := h.verificarWebhook(ctx, body, c.GetHeader(hmacHeaderName), integrationID, payload.StoreID.String()); err != nil {
		h.logger.Warn(ctx).Err(err).
			Str("event", payload.Event).
			Str("integration_id", integrationID).
			Str("store_id", payload.StoreID.String()).
			Msg("Webhook de Tiendanube rechazado: no supero la validacion de firma")
		c.Status(http.StatusUnauthorized)
		return
	}

	c.Status(http.StatusOK)

	event := strings.TrimSpace(payload.Event)
	resourceID := payload.ID.String()

	if !usecases.IsOrderEvent(event) {
		h.logger.Info(ctx).
			Str("event", event).
			Str("integration_id", integrationID).
			Msg("Webhook de Tiendanube ignorado: no es un evento de orden")
		return
	}

	if resourceID == "" || resourceID == "0" {
		h.logger.Warn(ctx).
			Str("event", event).
			Str("integration_id", integrationID).
			Msg("Webhook de orden de Tiendanube sin id de recurso, se descarta")
		return
	}

	go func() {
		bgCtx := context.WithoutCancel(ctx)
		if err := h.useCase.ProcessOrderEvent(bgCtx, integrationID, event, resourceID); err != nil {
			h.logger.Error(bgCtx).Err(err).
				Str("event", event).
				Str("integration_id", integrationID).
				Str("order_id", resourceID).
				Msg("Error procesando el webhook de orden de Tiendanube")
		}
	}()
}

func truncateBody(body []byte) string {
	const limit = 300
	if len(body) <= limit {
		return string(body)
	}
	return string(body[:limit]) + "..."
}

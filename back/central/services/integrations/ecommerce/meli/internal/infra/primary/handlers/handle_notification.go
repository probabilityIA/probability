package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/infra/secondary/client/response"
)

// HandleNotification recibe notificaciones IPN de MercadoLibre.
// MercadoLibre envía un POST con el topic y el resource en el body JSON.
// Se responde 200 inmediatamente y se procesa en background.
//
// Referencia: https://developers.mercadolibre.com/es_ar/recibir-notificaciones
func (h *meliHandler) HandleNotification(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. Leer body crudo
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Failed to read notification body")
		c.Status(http.StatusBadRequest)
		return
	}

	// 2. Deserializar la notificación
	var notifBody response.MeliNotificationBody
	if err := json.Unmarshal(rawBody, &notifBody); err != nil {
		h.logger.Warn(ctx).Err(err).
			Str("raw_body", string(rawBody)).
			Msg("Invalid MeLi notification payload")
		c.Status(http.StatusBadRequest)
		return
	}

	h.logger.Info(ctx).
		Str("topic", notifBody.Topic).
		Str("resource", notifBody.Resource).
		Int64("user_id", notifBody.UserID).
		Int("attempts", notifBody.Attempts).
		Msg("MercadoLibre IPN notification received")

	// 3. Responder 200 inmediatamente (MeLi espera respuesta rápida, si no, re-envía)
	c.Status(http.StatusOK)

	// 4. Convertir a dominio y procesar en background
	notification := notifBody.ToDomain()
	go h.processNotificationAsync(&notification)
}

func (h *meliHandler) processNotificationAsync(notification *domain.MeliNotification) {
	ctx := context.Background()

	if err := h.useCase.ProcessNotification(ctx, notification); err != nil {
		h.logger.Error(ctx).Err(err).
			Str("topic", notification.Topic).
			Str("resource", notification.Resource).
			Int64("user_id", notification.UserID).
			Msg("Failed to process MercadoLibre notification")
	}
}

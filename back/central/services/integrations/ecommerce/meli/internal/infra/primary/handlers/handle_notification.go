package handlers

import (
	"context"
	"io"
	"net/http"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/infra/secondary/client/response"
)

func (h *meliHandler) HandleNotification(c *gin.Context) {
	ctx := c.Request.Context()

	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Failed to read notification body")
		c.Status(http.StatusBadRequest)
		return
	}

	var notifBody response.MeliNotificationBody
	if err := json.Unmarshal(rawBody, &notifBody); err != nil {
		h.logger.Warn(ctx).Err(err).
			Str("raw_body", string(rawBody)).
			Msg("Invalid MeLi notification payload")
		c.Status(http.StatusBadRequest)
		return
	}

	h.verifyNotificationSignature(ctx, c, notifBody.Resource)

	h.logger.Info(ctx).
		Str("topic", notifBody.Topic).
		Str("resource", notifBody.Resource).
		Int64("user_id", notifBody.UserID).
		Int("attempts", notifBody.Attempts).
		Msg("MercadoLibre IPN notification received")

	c.Status(http.StatusOK)

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

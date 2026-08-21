package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type privacyWebhookPayload struct {
	StoreID        json.Number   `json:"store_id"`
	CustomerID     json.Number   `json:"customer_id"`
	OrdersToRedact []json.Number `json:"orders_to_redact"`
}

func (h *tiendanubeHandler) handlePrivacyWebhook(c *gin.Context, topic string) {
	ctx := c.Request.Context()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error(ctx).Err(err).
			Str("topic", topic).
			Msg("No se pudo leer el body del webhook de privacidad de Tiendanube")
		c.Status(http.StatusOK)
		return
	}

	var payload privacyWebhookPayload
	if len(body) > 0 {
		if err := json.Unmarshal(body, &payload); err != nil {
			h.logger.Warn(ctx).Err(err).
				Str("topic", topic).
				Str("body", string(body)).
				Msg("Body invalido en webhook de privacidad de Tiendanube")
		}
	}

	h.logger.Info(ctx).
		Str("topic", topic).
		Str("store_id", payload.StoreID.String()).
		Str("customer_id", payload.CustomerID.String()).
		Int("orders_to_redact", len(payload.OrdersToRedact)).
		Str("hmac", c.GetHeader("x-linkedstore-hmac-sha256")).
		Msg("Solicitud de privacidad recibida desde Tiendanube")

	c.Status(http.StatusOK)
}

func (h *tiendanubeHandler) HandleStoreRedact(c *gin.Context) {
	h.handlePrivacyWebhook(c, "store/redact")
}

func (h *tiendanubeHandler) HandleCustomersRedact(c *gin.Context) {
	h.handlePrivacyWebhook(c, "customers/redact")
}

func (h *tiendanubeHandler) HandleCustomersDataRequest(c *gin.Context) {
	h.handlePrivacyWebhook(c, "customers/data_request")
}

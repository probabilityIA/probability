package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleWebhook recibe webhooks de eventos de WooCommerce.
// WooCommerce envía un POST con el evento en el header X-WC-Webhook-Topic
// y el payload de la orden en el body.
//
// Referencia: https://woocommerce.com/document/webhooks/
func (h *wooCommerceHandler) HandleWebhook(c *gin.Context) {
	// TODO: implementar procesamiento de webhooks de WooCommerce
	// 1. Validar la firma HMAC del header X-WC-Webhook-Signature
	// 2. Leer el topic del header X-WC-Webhook-Topic
	// 3. Deserializar el payload de la orden
	// 4. Publicar evento a la cola de órdenes
	h.logger.Info(c.Request.Context()).
		Str("topic", c.GetHeader("X-WC-Webhook-Topic")).
		Str("source", c.GetHeader("X-WC-Webhook-Source")).
		Msg("WooCommerce webhook received")

	c.Status(http.StatusOK)
}

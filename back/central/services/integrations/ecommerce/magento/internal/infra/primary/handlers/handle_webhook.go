package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleWebhook recibe webhooks de eventos de Magento.
// Magento puede enviar webhooks via extensiones o integrations API.
//
// Referencia: https://developer.adobe.com/commerce/php/development/components/message-queues/
func (h *magentoHandler) HandleWebhook(c *gin.Context) {
	// TODO: implementar procesamiento de webhooks de Magento
	// 1. Validar la autenticación del webhook (token o HMAC)
	// 2. Leer el tipo de evento del payload
	// 3. Deserializar el payload de la orden
	// 4. Publicar evento a la cola de órdenes
	h.logger.Info(c.Request.Context()).
		Str("content_type", c.GetHeader("Content-Type")).
		Str("user_agent", c.GetHeader("User-Agent")).
		Msg("Magento webhook received")

	c.Status(http.StatusOK)
}

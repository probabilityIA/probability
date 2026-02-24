package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleWebhook recibe webhooks de eventos de Tiendanube.
// Tiendanube envia un POST con el evento en el body incluyendo el topic
// y el payload del recurso afectado.
//
// Referencia: https://tiendanube.github.io/api-documentation/resources/webhook
func (h *tiendanubeHandler) HandleWebhook(c *gin.Context) {
	// TODO: implementar procesamiento de webhooks de Tiendanube
	// 1. Validar la autenticidad del webhook
	// 2. Leer el evento/topic del payload
	// 3. Deserializar el payload de la orden
	// 4. Publicar evento a la cola de ordenes
	h.logger.Info(c.Request.Context()).
		Str("content_type", c.GetHeader("Content-Type")).
		Str("user_agent", c.GetHeader("User-Agent")).
		Msg("Tiendanube webhook received")

	c.Status(http.StatusOK)
}

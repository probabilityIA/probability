package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleWebhook recibe webhooks de eventos de Falabella Seller Center.
// Falabella envía un POST con el payload del evento en el body.
//
// TODO: implementar procesamiento de webhooks de Falabella Seller Center
func (h *falabellaHandler) HandleWebhook(c *gin.Context) {
	// TODO: implementar procesamiento de webhooks de Falabella
	// 1. Validar la autenticación del webhook
	// 2. Leer el tipo de evento del payload
	// 3. Deserializar el payload de la orden
	// 4. Publicar evento a la cola de órdenes
	h.logger.Info(c.Request.Context()).
		Str("content_type", c.GetHeader("Content-Type")).
		Msg("Falabella webhook received")

	c.Status(http.StatusOK)
}

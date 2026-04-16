package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleWebhook recibe webhooks de eventos de Exito marketplace.
//
// TODO: implementar procesamiento de webhooks de Exito
// 1. Validar la firma/autenticacion del webhook
// 2. Leer el tipo de evento
// 3. Deserializar el payload de la orden
// 4. Publicar evento a la cola de ordenes
func (h *exitoHandler) HandleWebhook(c *gin.Context) {
	h.logger.Info(c.Request.Context()).
		Str("content_type", c.GetHeader("Content-Type")).
		Msg("Exito webhook received")

	c.Status(http.StatusOK)
}

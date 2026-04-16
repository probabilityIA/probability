package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleNotification recibe notificaciones de eventos de Amazon.
// Amazon SP-API usa SQS/SNS para enviar notificaciones de cambios en ordenes,
// inventario y listings.
//
// TODO: Implementar procesamiento de notificaciones de Amazon SP-API
// 1. Validar la firma SNS del mensaje
// 2. Determinar el tipo de notificacion (ORDER_CHANGE, LISTINGS_ITEM_STATUS_CHANGE)
// 3. Deserializar el payload
// 4. Publicar evento a la cola de ordenes
//
// Referencia: https://developer-docs.amazon.com/sp-api/docs/notifications-api-v1-reference
func (h *amazonHandler) HandleNotification(c *gin.Context) {
	h.logger.Info(c.Request.Context()).
		Str("content_type", c.ContentType()).
		Str("method", c.Request.Method).
		Msg("Amazon notification received")

	c.Status(http.StatusOK)
}

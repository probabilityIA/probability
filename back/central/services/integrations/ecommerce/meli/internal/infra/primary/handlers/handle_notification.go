package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleNotification recibe notificaciones IPN de MercadoLibre.
// MercadoLibre envía una notificación POST con el topic y el resource ID.
//
// Referencia: https://developers.mercadolibre.com/es_ar/recibir-notificaciones
func (h *meliHandler) HandleNotification(c *gin.Context) {
	// TODO: implementar procesamiento de notificaciones IPN
	// 1. Validar el x-signature header
	// 2. Leer topic y resource_id del query string
	// 3. Obtener detalles del recurso vía API de MercadoLibre
	// 4. Publicar evento a la cola de órdenes
	h.logger.Info(c.Request.Context()).
		Str("topic", c.Query("topic")).
		Str("resource", c.Query("resource")).
		Msg("MercadoLibre IPN notification received")

	c.Status(http.StatusOK)
}

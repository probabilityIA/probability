package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

// GetShipmentByTrackingNumber godoc
// @Summary      Obtener envío por número de tracking
// @Description  Obtiene un envío específico por su número de tracking
// @Tags         Shipments
// @Accept       json
// @Produce      json
// @Param        tracking_number   path      string  true  "Número de tracking"
// @Security     BearerAuth
// @Success      200  {object}  domain.ShipmentResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /shipments/tracking/{tracking_number} [get]
func (h *Handlers) GetShipmentByTrackingNumber(c *gin.Context) {
	trackingNumber := c.Param("tracking_number")

	if trackingNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Número de tracking es requerido",
			"error":   "El tracking_number es requerido",
		})
		return
	}

	// Llamar al caso de uso
	shipment, err := h.uc.GetShipmentByTrackingNumber(c.Request.Context(), trackingNumber)
	if err != nil {
		if err == domain.ErrShipmentNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Envío no encontrado",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener envío",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Envío obtenido exitosamente",
		"data":    shipment,
	})
}


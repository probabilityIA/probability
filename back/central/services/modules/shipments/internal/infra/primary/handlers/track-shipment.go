package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TrackShipment godoc
// @Summary      Rastrear envío en EnvioClick
// @Description  Obtiene el rastreo en tiempo real desde la API de EnvioClick
// @Tags         Shipments
// @Accept       json
// @Produce      json
// @Param        tracking_number   path      string  true  "Número de tracking"
// @Security     BearerAuth
// @Success      200  {object}  domain.EnvioClickTrackingResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /shipments/tracking/{tracking_number}/track [post]
func (h *Handlers) TrackEnvioclickShipment(c *gin.Context) {
	trackingNumber := c.Param("tracking_number")

	if trackingNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Número de tracking es requerido",
		})
		return
	}

	resp, err := h.envioClickUC.TrackShipment(c.Request.Context(), trackingNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al rastrear envío en EnvioClick",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Rastreo obtenido exitosamente",
		"data":    resp.Data,
	})
}

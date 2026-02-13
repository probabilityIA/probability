package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CancelEnvioclickShipment godoc
// @Summary      Cancelar envío en EnvioClick
// @Description  Cancela un envío directamente en la API de EnvioClick
// @Tags         Shipments
// @Accept       json
// @Produce      json
// @Param        id                path      string  true  "ID de Envío (shipment id o tracking)"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /shipments/:id/cancel [post]
func (h *Handlers) CancelEnvioclickShipment(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de envío es requerido",
		})
		return
	}

	resp, err := h.envioClickUC.CancelShipment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al cancelar envío en EnvioClick",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": resp.Message,
		"data":    resp,
	})
}

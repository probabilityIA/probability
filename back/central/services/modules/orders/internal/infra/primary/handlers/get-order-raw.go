package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetOrderRaw godoc
// @Summary      Obtener datos crudos de orden
// @Description  Obtiene el JSON original recibido del canal de venta
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID de la orden (UUID)"
// @Security     BearerAuth
// @Success      200  {object}  domain.OrderRawResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders/{id}/raw [get]
func (h *Handlers) GetOrderRaw(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de orden inválido",
			"error":   "El ID de la orden es requerido",
		})
		return
	}

	// Llamar al caso de uso
	rawResponse, err := h.orderCRUD.GetOrderRaw(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "raw data not found for this order" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Datos crudos no encontrados",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener datos crudos",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Datos crudos obtenidos exitosamente",
		"data":    rawResponse,
	})
}

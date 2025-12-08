package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetOrderByID godoc
// @Summary      Obtener orden por ID
// @Description  Obtiene una orden específica por su ID
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID de la orden (UUID)"
// @Security     BearerAuth
// @Success      200  {object}  domain.OrderResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders/{id} [get]
func (h *Handlers) GetOrderByID(c *gin.Context) {
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
	order, err := h.orderCRUD.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Orden no encontrada",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener orden",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Orden obtenida exitosamente",
		"data":    order,
	})
}

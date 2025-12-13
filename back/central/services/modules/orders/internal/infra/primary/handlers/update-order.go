package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders/domain"
)

// UpdateOrder godoc
// @Summary      Actualizar orden
// @Description  Actualiza una orden existente
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id     path      string                     true  "ID de la orden (UUID)"
// @Param        order  body      domain.UpdateOrderRequest  true  "Datos a actualizar"
// @Security     BearerAuth
// @Success      200  {object}  domain.OrderResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders/{id} [put]
func (h *Handlers) UpdateOrder(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de orden inválido",
			"error":   "El ID de la orden es requerido",
		})
		return
	}

	var req domain.UpdateOrderRequest

	// Validar el request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

	// Llamar al caso de uso
	order, err := h.orderCRUD.UpdateOrder(c.Request.Context(), id, &req)
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
			"message": "Error al actualizar orden",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Orden actualizada exitosamente",
		"data":    order,
	})
}

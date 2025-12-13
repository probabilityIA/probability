package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders/domain"
)

// CreateOrder godoc
// @Summary      Crear orden
// @Description  Crea una nueva orden en el sistema
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        order  body      domain.CreateOrderRequest  true  "Datos de la orden"
// @Security     BearerAuth
// @Success      201  {object}  domain.OrderResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders [post]
func (h *Handlers) CreateOrder(c *gin.Context) {
	var req domain.CreateOrderRequest

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
	order, err := h.orderCRUD.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		// Verificar si es un error de duplicado
		if err.Error() == "order with this external_id already exists for this integration" {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Orden ya existe",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al crear orden",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Orden creada exitosamente",
		"data":    order,
	})
}

package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	orderErrors "github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/request"
)

// UpdateOrder godoc
// @Summary      Actualizar orden
// @Description  Actualiza una orden existente
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id     path      string                              true  "ID de la orden (UUID)"
// @Param        order  body      request.UpdateOrder                 true  "Datos a actualizar"
// @Security     BearerAuth
// @Success      200  {object}  dtos.OrderResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders/{id} [put]
func (h *Handlers) UpdateOrder(c *gin.Context) {
	fmt.Printf("\n\n🚀🚀🚀 [HANDLER] UpdateOrder INICIADO 🚀🚀🚀\n\n")
	id := c.Param("id")
	fmt.Printf("📌 Order ID: %s\n", id)

	if id == "" {
		fmt.Printf("❌ Order ID vacío\n")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de orden inválido",
			"error":   "El ID de la orden es requerido",
		})
		return
	}

	var req request.UpdateOrder

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

	h.logger.Info(c.Request.Context()).
		Str("order_id", id).
		Int("items_received_in_handler", len(req.Items)).
		Msg("🔍 UpdateOrder handler: items received")

	fmt.Printf("🔍 [HANDLER] UpdateOrder iniciado - order_id=%s, items_received=%d\n", id, len(req.Items))
	if len(req.Items) > 0 {
		fmt.Printf("   Items recibidos:\n")
		for i, item := range req.Items {
			fmt.Printf("   [%d] %v\n", i+1, item)
		}
	}

	c.Request.Header.Set("X-Debug-Items-Count", fmt.Sprintf("%d", len(req.Items)))

	domainReq := mappers.MapUpdateOrderRequestToDomain(&req)

	order, err := h.orderCRUD.UpdateOrder(c.Request.Context(), id, domainReq)
	if err != nil {
		if errors.Is(err, orderErrors.ErrOrderNotFound) {
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
		"data":    mappers.OrderToResponse(order),
	})
}

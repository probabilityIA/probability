package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	orderErrors "github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/request"
)

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

	var req request.UpdateOrder

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

	if !h.ensureOrderOwnership(c, id) {
		return
	}

	domainReq := mappers.MapUpdateOrderRequestToDomain(&req)
	if uid, ok := middleware.GetUserID(c); ok && uid > 0 {
		domainReq.UpdatedBy = &uid
	}

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

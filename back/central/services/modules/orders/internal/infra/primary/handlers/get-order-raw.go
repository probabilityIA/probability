package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/mappers"
)

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

	if !h.ensureOrderOwnership(c, id) {
		return
	}

	domainResp, err := h.orderCRUD.GetOrderRaw(c.Request.Context(), id)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "raw data not found for this order") {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Esta orden no tiene datos crudos guardados. Los datos crudos solo están disponibles para órdenes creadas después de la implementación de esta funcionalidad.",
				"error":   "raw data not found",
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

	httpResp := mappers.OrderRawToResponse(domainResp)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Datos crudos obtenidos exitosamente",
		"data":    httpResp,
	})
}

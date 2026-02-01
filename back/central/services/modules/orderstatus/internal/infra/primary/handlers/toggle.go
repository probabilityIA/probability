package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/primary/handlers/mappers"
)

// Toggle godoc
// @Summary      Alternar estado activo/inactivo
// @Description  Cambia el estado activo/inactivo de un mapeo
// @Tags         Order Status Mappings
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del mapeo"
// @Success      200  {object}  response.OrderStatusMappingResponse
// @Failure      400  {object}  map[string]string
// @Router       /order-status-mappings/{id}/toggle [patch]
func (h *handler) Toggle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	result, err := h.uc.ToggleOrderStatusMappingActive(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mappers.DomainToResponse(result, h.getImageURLBase()))
}

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TogglePaymentMapping godoc
// @Summary      Activar/Desactivar mapeo
// @Description  Cambia el estado activo/inactivo de un mapeo
// @Tags         Payment Mappings
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del mapeo"
// @Success      200  {object}  domain.PaymentMappingResponse
// @Failure      400  {object}  map[string]string
// @Router       /payments/mappings/{id}/toggle [patch]
func (h *PaymentHandlers) TogglePaymentMapping(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	response, err := h.uc.TogglePaymentMappingActive(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

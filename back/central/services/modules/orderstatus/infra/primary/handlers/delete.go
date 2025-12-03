package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary      Eliminar mapeo de estado de orden
// @Description  Elimina un mapeo de estado de orden
// @Tags         Order Status Mappings
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del mapeo"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]string
// @Router       /order-status-mappings/{id} [delete]
func (h *OrderStatusMappingHandlers) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.uc.DeleteOrderStatusMapping(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

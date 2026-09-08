package scheduled_rule

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) GetByID(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
		return
	}

	id, ok := parseID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}

	rule, err := h.useCase.GetRule(c.Request.Context(), id, businessID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": err.Error()})
		return
	}
	if rule == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "regla no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": rule})
}

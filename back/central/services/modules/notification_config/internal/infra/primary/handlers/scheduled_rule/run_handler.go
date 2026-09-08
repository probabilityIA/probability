package scheduled_rule

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) RunNow(c *gin.Context) {
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

	run, err := h.useCase.RunRuleNow(c.Request.Context(), id, businessID)
	if err != nil {
		h.logger.Error().Err(err).Uint("id", id).Msg("Error ejecutando la regla programada")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error(), "data": run})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": run})
}

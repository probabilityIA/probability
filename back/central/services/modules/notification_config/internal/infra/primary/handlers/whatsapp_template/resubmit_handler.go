package whatsapp_template

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) Resubmit(c *gin.Context) {
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

	template, err := h.useCase.Resubmit(c.Request.Context(), id, businessID)
	if err != nil {
		h.logger.Error().Err(err).Uint("id", id).Msg("Error reenviando plantilla a Meta")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": template})
}

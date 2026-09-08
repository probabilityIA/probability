package whatsapp_template

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) Delete(c *gin.Context) {
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

	if err := h.useCase.Delete(c.Request.Context(), id, businessID); err != nil {
		h.logger.Error().Err(err).Uint("id", id).Msg("Error eliminando plantilla de WhatsApp")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

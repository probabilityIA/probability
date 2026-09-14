package whatsapp_template

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) SyncStatuses(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
		return
	}

	pending, err := h.useCase.SyncPendingStatuses(c.Request.Context(), businessID)
	if err != nil {
		h.logger.Error().Err(err).Uint("business_id", businessID).Msg("Error consultando el estado de las plantillas en Meta")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"pending": pending}})
}

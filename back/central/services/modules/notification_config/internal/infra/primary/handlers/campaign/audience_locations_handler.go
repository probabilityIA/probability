package campaign

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) ListAudienceLocations(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
		return
	}

	locations, err := h.useCase.ListAudienceLocations(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": locations})
}

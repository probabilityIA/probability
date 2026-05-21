package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) DeleteClientGroup(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	groupID, ok := parseUintParam(c, "id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	if err := h.uc.DeleteClientGroup(c.Request.Context(), businessID, groupID); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "client group deleted"})
}

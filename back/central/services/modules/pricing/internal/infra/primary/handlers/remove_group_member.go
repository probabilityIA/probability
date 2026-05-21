package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) RemoveGroupMember(c *gin.Context) {
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

	clientID, ok := parseUintParam(c, "clientId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
		return
	}

	if err := h.uc.RemoveGroupMember(c.Request.Context(), businessID, groupID, clientID); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member removed from group"})
}

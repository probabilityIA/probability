package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) Delete(c *gin.Context) {
	id, ok := h.parseIDParam(c)
	if !ok {
		return
	}

	if err := h.uc.DeleteAnnouncement(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "announcement deleted"})
}

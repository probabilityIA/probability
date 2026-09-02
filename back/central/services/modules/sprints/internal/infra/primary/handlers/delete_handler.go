package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) Delete(c *gin.Context) {
	if _, ok := h.requireSuperAdmin(c); !ok {
		return
	}
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

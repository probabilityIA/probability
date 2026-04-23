package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) ListEvents(c *gin.Context) {
	events := h.uc.ListEvents(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"events": events})
}

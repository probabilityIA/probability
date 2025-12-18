package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) GenerateGuide(c *gin.Context) {
	var req interface{} // Receive arbitrary JSON payload for now
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.envioClickUC.GenerateGuide(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetModuleCodes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": h.uc.GetModuleCodes()})
}

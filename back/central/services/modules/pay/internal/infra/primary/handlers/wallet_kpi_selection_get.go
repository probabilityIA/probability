package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *walletHandler) GetKPISelection(c *gin.Context) {
	selection, err := h.walletUC.GetWalletKPISelection(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    selection,
	})
}

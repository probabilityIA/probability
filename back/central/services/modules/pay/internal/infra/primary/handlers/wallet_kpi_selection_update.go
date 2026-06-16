package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
)

func (h *walletHandler) UpdateKPISelection(c *gin.Context) {
	var req dtos.UpdateWalletKPISelectionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid request"})
		return
	}

	selection, err := h.walletUC.UpdateWalletKPISelection(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    selection,
	})
}

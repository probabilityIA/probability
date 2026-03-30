package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BoldGenerateSignature genera la firma de integridad para Bold.co
func (h *walletHandler) BoldGenerateSignature(c *gin.Context) {
	var req struct {
		Amount   float64 `json:"amount" binding:"required"`
		Currency string  `json:"currency" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: amount and currency are required"})
		return
	}

	resp, err := h.walletUC.BoldGenerateSignature(c.Request.Context(), req.Amount, req.Currency)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Msg("Failed to generate Bold signature")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate integrity signature"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetBoldStatus consulta el estado de una orden en Bold.co
func (h *walletHandler) GetBoldStatus(c *gin.Context) {
	boldOrderID := c.Param("id")
	if boldOrderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	resp, err := h.walletUC.GetBoldStatus(c.Request.Context(), boldOrderID)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Str("bold_order_id", boldOrderID).Msg("Failed to fetch Bold status")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction status"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

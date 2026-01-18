package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

type RechargeRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

func (h *WalletHandlers) RechargeWallet(c *gin.Context) {
	businessID, exists := middleware.GetBusinessID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount < 15000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El monto mínimo de recarga es de $15,000"})
		return
	}

	qr, err := h.uc.Recharge(c.Request.Context(), businessID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"qr_code": qr})
}

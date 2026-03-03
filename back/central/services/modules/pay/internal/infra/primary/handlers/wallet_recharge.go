package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/primary/handlers/request"
)

// RechargeWallet maneja POST /pay/wallet/recharge
func (h *walletHandler) RechargeWallet(c *gin.Context) {
	var req request.RechargeWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var businessID uint
	if middleware.IsSuperAdmin(c) {
		if req.BusinessID == nil || *req.BusinessID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido para super admin"})
			return
		}
		businessID = *req.BusinessID
	} else {
		jwtBusinessID, exists := middleware.GetBusinessID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		businessID = jwtBusinessID
	}

	tx, err := h.walletUC.RechargeWallet(c.Request.Context(), &dtos.RechargeWalletDTO{
		BusinessID: businessID,
		Amount:     req.Amount,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"qr_code": tx.QrCode})
}

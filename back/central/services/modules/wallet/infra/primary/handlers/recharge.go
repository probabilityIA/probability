package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

type RechargeRequest struct {
	Amount     float64 `json:"amount" binding:"required"`
	BusinessID *uint   `json:"business_id"` // Requerido cuando super admin actúa en nombre de un negocio
}

func (h *WalletHandlers) RechargeWallet(c *gin.Context) {
	var req RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var businessID uint
	if middleware.IsSuperAdmin(c) {
		// Super admin debe proveer business_id en el body
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

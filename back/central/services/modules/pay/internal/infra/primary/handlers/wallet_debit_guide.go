package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/primary/handlers/request"
)

// DebitForGuide maneja POST /pay/wallet/debit-guide
func (h *walletHandler) DebitForGuide(c *gin.Context) {
	var req request.DebitForGuideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	businessID, exists := middleware.GetBusinessID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Business ID not found in token"})
		return
	}

	if err := h.walletUC.DebitForGuide(c.Request.Context(), &dtos.DebitForGuideDTO{
		BusinessID:     businessID,
		Amount:         req.Amount,
		TrackingNumber: req.TrackingNumber,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet debited successfully", "tracking_number": req.TrackingNumber})
}

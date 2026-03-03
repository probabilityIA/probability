package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/primary/handlers/mappers"
)

// GetBalance maneja GET /pay/wallet/balance
func (h *walletHandler) GetBalance(c *gin.Context) {
	businessID, ok := resolveBusinessID(c)
	if !ok {
		return
	}

	wallet, err := h.walletUC.GetWallet(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mappers.WalletToResponse(wallet))
}

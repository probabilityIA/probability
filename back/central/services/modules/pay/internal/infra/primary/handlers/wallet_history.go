package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/primary/handlers/mappers"
)

// GetHistory maneja GET /pay/wallet/history
func (h *walletHandler) GetHistory(c *gin.Context) {
	businessID, ok := resolveBusinessID(c)
	if !ok {
		return
	}

	txs, err := h.walletUC.GetTransactionsByBusinessID(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mappers.WalletTxListToResponse(txs))
}

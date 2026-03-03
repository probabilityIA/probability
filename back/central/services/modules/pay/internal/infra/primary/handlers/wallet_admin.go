package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/primary/handlers/request"
)

// GetAllWallets maneja GET /pay/wallet/all (admin)
func (h *walletHandler) GetAllWallets(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		roles, _ := middleware.GetUserRoles(c)
		isAdmin := false
		for _, r := range roles {
			if r == "Admin" || r == "Administrador" {
				isAdmin = true
				break
			}
		}
		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
	}

	wallets, err := h.walletUC.GetAllWallets(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mappers.WalletListToResponse(wallets))
}

// GetPendingTransactions maneja GET /pay/wallet/admin/pending-requests
func (h *walletHandler) GetPendingTransactions(c *gin.Context) {
	txs, err := h.walletUC.GetPendingTransactions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mappers.WalletTxListToResponse(txs))
}

// GetProcessedTransactions maneja GET /pay/wallet/admin/processed-requests
func (h *walletHandler) GetProcessedTransactions(c *gin.Context) {
	txs, err := h.walletUC.GetProcessedTransactions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mappers.WalletTxListToResponse(txs))
}

// ApproveTransaction maneja POST /pay/wallet/admin/requests/:id/approve
func (h *walletHandler) ApproveTransaction(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	}

	if err := h.walletUC.ApproveTransaction(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction approved"})
}

// RejectTransaction maneja POST /pay/wallet/admin/requests/:id/reject
func (h *walletHandler) RejectTransaction(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	}

	if err := h.walletUC.RejectTransaction(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction rejected"})
}

// ManualDebit maneja POST /pay/wallet/admin/manual-debit
func (h *walletHandler) ManualDebit(c *gin.Context) {
	var req request.ManualDebitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.walletUC.ManualDebit(c.Request.Context(), &dtos.ManualDebitDTO{
		BusinessID: req.BusinessID,
		Amount:     req.Amount,
		Reference:  req.Reference,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Balance subtracted successfully"})
}

// ClearRechargeHistory maneja DELETE /pay/wallet/admin/history/:business_id
func (h *walletHandler) ClearRechargeHistory(c *gin.Context) {
	businessIDStr := c.Param("business_id")
	if businessIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Business ID is required"})
		return
	}

	var businessID uint
	if _, err := fmt.Sscanf(businessIDStr, "%d", &businessID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid business ID"})
		return
	}

	if err := h.walletUC.ClearRechargeHistory(c.Request.Context(), businessID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recharge history cleared successfully"})
}

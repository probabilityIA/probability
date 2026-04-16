package handlers

import (
	"fmt"
	"net/http"
	"time"

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

// AdminAdjustBalance maneja POST /pay/wallet/admin/adjust-balance (sin restricción de mínimo)
func (h *walletHandler) AdminAdjustBalance(c *gin.Context) {
	var req struct {
		BusinessID uint    `json:"business_id" binding:"required"`
		Amount     float64 `json:"amount" binding:"required"`
		Reference  string  `json:"reference" binding:"required,max=255"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.walletUC.AdminAdjustBalance(c.Request.Context(), &dtos.AdminAdjustBalanceDTO{
		BusinessID: req.BusinessID,
		Amount:     req.Amount,
		Reference:  req.Reference,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Balance adjusted successfully"})
}

// GetFinancialStats maneja GET /pay/wallet/admin/financial-stats
func (h *walletHandler) GetFinancialStats(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// Parsear query params
	businessIDStr := c.Query("business_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	month := c.Query("month")

	// Si no hay fechas explícitas, usar el mes actual o el especificado
	if startDate == "" && endDate == "" && month == "" {
		// Usar mes actual (default)
		now := time.Now()
		startDate = now.Format("2006-01-02")[:7] + "-01"
		endDate = now.AddDate(0, 1, -now.Day()).Format("2006-01-02")
	} else if month != "" {
		// Usar mes especificado (YYYY-MM)
		startDate = month + "-01"
		endDate = time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location()).
			AddDate(0, 1, -1).Format("2006-01-02")
		if len(month) >= 7 {
			// Parsear año-mes
			year := month[:4]
			monthStr := month[5:]
			endDate = fmt.Sprintf("%s-%s-31", year, monthStr)
			// Ajustar al último día del mes
			parsedDate, _ := time.Parse("2006-01", year+"-"+monthStr)
			endDate = parsedDate.AddDate(0, 1, -parsedDate.Day()).Format("2006-01-02")
		}
	}

	// Parsear business_id si existe
	var businessID *uint
	if businessIDStr != "" {
		var id uint
		if _, err := fmt.Sscanf(businessIDStr, "%d", &id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid business_id"})
			return
		}
		businessID = &id
	}

	// Llamar usecase
	stats, err := h.walletUC.GetFinancialStats(c.Request.Context(), &dtos.FinancialStatsDTO{
		BusinessID: businessID,
		StartDate:  startDate,
		EndDate:    endDate,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Helper para conversión de fechas de mes
func parseMonth(month string) (start, end string, err error) {
	parsedDate, err := time.Parse("2006-01", month)
	if err != nil {
		return "", "", err
	}
	start = parsedDate.Format("2006-01-02")
	end = parsedDate.AddDate(0, 1, -parsedDate.Day()).Format("2006-01-02")
	return
}

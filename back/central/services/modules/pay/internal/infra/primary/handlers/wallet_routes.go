package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterWalletRoutes registra las rutas de wallet bajo /pay/wallet
func (h *walletHandler) RegisterWalletRoutes(router *gin.RouterGroup) {
	wallet := router.Group("/pay/wallet")
	wallet.Use(middleware.JWT())
	{
		// Rutas de negocio
		wallet.GET("/balance", h.GetBalance)
		wallet.POST("/recharge", h.RechargeWallet)
		wallet.GET("/history", h.GetHistory)
		wallet.POST("/debit-guide", h.DebitForGuide)

		// Rutas de admin
		wallet.GET("/all", h.GetAllWallets)
		wallet.GET("/admin/pending-requests", h.GetPendingTransactions)
		wallet.GET("/admin/processed-requests", h.GetProcessedTransactions)
		wallet.POST("/admin/requests/:id/approve", h.ApproveTransaction)
		wallet.POST("/admin/requests/:id/reject", h.RejectTransaction)
		wallet.POST("/admin/manual-debit", h.ManualDebit)
		wallet.DELETE("/admin/history/:business_id", h.ClearRechargeHistory)
	}
}

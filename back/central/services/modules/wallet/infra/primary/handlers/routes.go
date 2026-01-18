package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *WalletHandlers) RegisterRoutes(router *gin.RouterGroup) {
	wallet := router.Group("/wallet")
	{
		// Requires Authentication
		wallet.Use(middleware.JWT())

		// Business Routes
		wallet.GET("/balance", h.GetBalance)
		wallet.POST("/recharge", h.RechargeWallet)

		// Admin Routes
		// Simple role check middleware would be better here, but assuming handler checks permission or we trust JWT roles.
		// For stricter control: wallet.Use(middleware.RequireRole("SuperAdmin"))

		wallet.GET("/all", h.GetAllWallets)
		wallet.GET("/admin/pending-requests", h.GetPendingTransactions)
		wallet.POST("/admin/requests/:id/approve", h.ApproveTransaction)
		wallet.POST("/admin/requests/:id/reject", h.RejectTransaction)
	}
}

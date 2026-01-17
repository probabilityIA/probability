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
		// wallet.GET("/all", middleware.RequireRole("Admin"), h.GetAllWallets) // "Admin" or "SuperAdmin"?
		// The middleware `IsSuperAdmin` checks businessID == 0.
		// I will use `IsSuperAdmin` check inside GetAllWallets or use middleware if available for specific roles.
		// Prompt says "para los que son admin".
		wallet.GET("/all", h.GetAllWallets)
	}
}

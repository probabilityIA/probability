package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *WalletHandlers) GetAllWallets(c *gin.Context) {
	// Check if user is Admin
	// middleware.IsSuperAdmin(c) or check role
	if !middleware.IsSuperAdmin(c) {
		// Also allow "Admin" role?
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

	wallets, err := h.uc.GetAllWallets(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wallets)
}

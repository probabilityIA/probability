package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	pub := router.Group("/public/siigo-referrals")
	{
		pub.POST("", h.CreateReferral)
	}

	admin := router.Group("/siigo-referrals")
	{
		admin.GET("", middleware.JWT(), middleware.RequireSuperAdmin(), h.ListReferrals)
	}
}

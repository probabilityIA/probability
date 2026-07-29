package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	pub := router.Group("/public/marketing-leads")
	{
		pub.POST("", h.CreateLead)
	}

	admin := router.Group("/marketing-leads")
	{
		admin.GET("", middleware.JWT(), middleware.RequireSuperAdmin(), h.ListLeads)
	}
}

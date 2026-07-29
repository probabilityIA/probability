package handlers

import "github.com/gin-gonic/gin"

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	pub := router.Group("/public/marketing-leads")
	{
		pub.POST("", h.CreateLead)
	}
}

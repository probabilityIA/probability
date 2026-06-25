package handlers

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(v1Group *gin.RouterGroup) {
	demoGroup := v1Group.Group("/auth")
	{
		demoGroup.POST("/demo-register", h.DemoRegisterHandler)
		demoGroup.POST("/verify-email", h.VerifyEmailHandler)
	}
}

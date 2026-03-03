package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/reference-data", h.GetReferenceData)
	router.POST("/generate", h.GenerateOrders)
}

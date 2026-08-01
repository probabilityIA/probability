package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandler interface {
	SyncShopInfo(c *gin.Context)
	ListSnapshots(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}

type tiktokHandler struct {
	useCase usecases.ITikTokUseCase
	logger  log.ILogger
}

func New(useCase usecases.ITikTokUseCase, logger log.ILogger) IHandler {
	return &tiktokHandler{
		useCase: useCase,
		logger:  logger.WithModule("tiktok"),
	}
}

func (h *tiktokHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	tiktok := router.Group("/tiktok")
	{
		tiktok.POST("/shop-info/sync", middleware.JWT(), h.SyncShopInfo)
		tiktok.GET("/shop-info/snapshots", middleware.JWT(), h.ListSnapshots)
	}
}

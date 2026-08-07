package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/commercial/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandlers interface {
	ListProspects(c *gin.Context)
	GetStats(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

type Handlers struct {
	uc     app.IUseCase
	logger log.ILogger
}

func New(uc app.IUseCase, logger log.ILogger) IHandlers {
	return &Handlers{uc: uc, logger: logger}
}

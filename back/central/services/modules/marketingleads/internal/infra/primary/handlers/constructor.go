package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandlers interface {
	CreateLead(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

type Handlers struct {
	uc     app.IUseCase
	logger log.ILogger
}

func New(uc app.IUseCase, logger log.ILogger) IHandlers {
	return &Handlers{uc: uc, logger: logger}
}

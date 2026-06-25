package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/demo/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandler interface {
	DemoRegisterHandler(c *gin.Context)
	VerifyEmailHandler(c *gin.Context)
	RegisterRoutes(v1Group *gin.RouterGroup)
}

type Handler struct {
	usecase app.IUseCase
	logger  log.ILogger
}

func New(usecase app.IUseCase, logger log.ILogger) IHandler {
	return &Handler{
		usecase: usecase,
		logger:  logger.WithModule("demo"),
	}
}

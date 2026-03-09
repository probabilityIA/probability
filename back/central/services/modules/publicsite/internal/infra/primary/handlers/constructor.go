package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/app"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandlers interface {
	GetBusinessPage(c *gin.Context)
	ListCatalog(c *gin.Context)
	GetProduct(c *gin.Context)
	SubmitContact(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

type Handlers struct {
	uc     app.IUseCase
	logger log.ILogger
	env    env.IConfig
}

func New(uc app.IUseCase, logger log.ILogger, environment env.IConfig) IHandlers {
	return &Handlers{uc: uc, logger: logger, env: environment}
}

func (h *Handlers) getImageURLBase() string {
	return h.env.Get("URL_BASE_DOMAIN_S3")
}

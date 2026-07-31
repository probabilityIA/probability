package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/app"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/storage"
)

type IHandlers interface {
	GetConfig(c *gin.Context)
	UpdateConfig(c *gin.Context)
	UploadImage(c *gin.Context)
	DeleteImage(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

type Handlers struct {
	uc     app.IUseCase
	logger log.ILogger
	s3     storage.IS3Service
	env    env.IConfig
}

func New(uc app.IUseCase, logger log.ILogger, s3 storage.IS3Service, environment env.IConfig) IHandlers {
	return &Handlers{uc: uc, logger: logger, s3: s3, env: environment}
}

func (h *Handlers) resolveBusinessID(c *gin.Context) (uint, bool) {
	businessID := c.GetUint("business_id")
	if businessID > 0 {
		return businessID, true
	}
	if param := c.Query("business_id"); param != "" {
		if id, err := strconv.ParseUint(param, 10, 64); err == nil && id > 0 {
			return uint(id), true
		}
	}
	return 0, false
}

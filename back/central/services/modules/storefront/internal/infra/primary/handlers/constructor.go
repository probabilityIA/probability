package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/app"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/storage"
)

type IHandlers interface {
	ListCatalog(c *gin.Context)
	GetCatalogFilters(c *gin.Context)
	UpdateCatalogLayout(c *gin.Context)
	UpdateCatalogBanner(c *gin.Context)
	UploadCatalogBannerImage(c *gin.Context)
	GetProduct(c *gin.Context)
	CreateOrder(c *gin.Context)
	ListMyOrders(c *gin.Context)
	GetMyOrder(c *gin.Context)
	CreateClient(c *gin.Context)
	ListClients(c *gin.Context)
	Register(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

type Handlers struct {
	uc     app.IUseCase
	logger log.ILogger
	env    env.IConfig
	s3     storage.IS3Service
}

func New(uc app.IUseCase, logger log.ILogger, environment env.IConfig, s3 storage.IS3Service) IHandlers {
	return &Handlers{uc: uc, logger: logger, env: environment, s3: s3}
}

func (h *Handlers) getImageURLBase() string {
	return h.env.Get("URL_BASE_DOMAIN_S3")
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

package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/app"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IHandlers defines the storefront handlers interface
type IHandlers interface {
	ListCatalog(c *gin.Context)
	GetProduct(c *gin.Context)
	CreateOrder(c *gin.Context)
	ListMyOrders(c *gin.Context)
	GetMyOrder(c *gin.Context)
	Register(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

// Handlers contains the use case
type Handlers struct {
	uc     app.IUseCase
	logger log.ILogger
	env    env.IConfig
}

// New creates a new instance of the storefront handlers
func New(uc app.IUseCase, logger log.ILogger, environment env.IConfig) IHandlers {
	return &Handlers{uc: uc, logger: logger, env: environment}
}

// getImageURLBase returns the S3 base URL for building full image URLs
func (h *Handlers) getImageURLBase() string {
	return h.env.Get("URL_BASE_DOMAIN_S3")
}

// resolveBusinessID gets business_id from JWT context, falls back to query param (for super admins)
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

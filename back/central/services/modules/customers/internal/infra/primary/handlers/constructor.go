package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/app"
)

type IHandlers interface {
	ListClients(c *gin.Context)
	GetClient(c *gin.Context)
	CreateClient(c *gin.Context)
	UpdateClient(c *gin.Context)
	DeleteClient(c *gin.Context)
	GetCustomerSummary(c *gin.Context)
	ListCustomerAddresses(c *gin.Context)
	ListCustomerProducts(c *gin.Context)
	ListCustomerOrderItems(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

type Handlers struct {
	uc app.IUseCase
}

func New(uc app.IUseCase) IHandlers {
	return &Handlers{uc: uc}
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

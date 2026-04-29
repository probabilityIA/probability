package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/app"
)

type IHandlers interface {
	List(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
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

func (h *Handlers) requireSuperAdmin(c *gin.Context) bool {
	if c.GetBool("is_super_admin") {
		return true
	}
	c.JSON(http.StatusForbidden, gin.H{"error": "super admin access required"})
	c.Abort()
	return false
}

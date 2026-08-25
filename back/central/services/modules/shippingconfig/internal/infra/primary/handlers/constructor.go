package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shippingconfig/internal/domain"
)

type Handlers struct {
	uc domain.IUseCase
}

func New(uc domain.IUseCase) *Handlers {
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

func (h *Handlers) resolveActor(c *gin.Context) (uint, string) {
	userID := c.GetUint("user_id")
	name, _ := c.Get("user_name")
	userName, _ := name.(string)
	return userID, userName
}

func parseUintParam(c *gin.Context, key string) (uint, bool) {
	raw := c.Param(key)
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return 0, false
	}
	return uint(id), true
}

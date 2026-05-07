package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/ports"
)

type IHandlers interface {
	List(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Bulk(c *gin.Context)
	Lookup(c *gin.Context)
	Delete(c *gin.Context)
	Display(c *gin.Context)
	FlushDisplayCache(c *gin.Context)
	Probability(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

type Handlers struct {
	uc            app.IUseCase
	probabilityUC ports.IProbabilityUseCase
}

func New(uc app.IUseCase, probabilityUC ports.IProbabilityUseCase) IHandlers {
	return &Handlers{uc: uc, probabilityUC: probabilityUC}
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

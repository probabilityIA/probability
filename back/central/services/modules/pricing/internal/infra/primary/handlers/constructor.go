package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/app"
)

// Handlers contiene el use case
type Handlers struct {
	uc app.IUseCase
}

// New crea una nueva instancia de los handlers
func New(uc app.IUseCase) *Handlers {
	return &Handlers{uc: uc}
}

// resolveBusinessID obtiene el business_id efectivo.
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

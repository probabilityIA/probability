package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Handlers struct {
	useCase ports.IUseCase
	logger  log.ILogger
}

func New(useCase ports.IUseCase, logger log.ILogger) *Handlers {
	return &Handlers{
		useCase: useCase,
		logger:  logger.WithModule("orderscompare.handlers"),
	}
}

func (h *Handlers) resolveBusinessID(c *gin.Context, bodyBusinessID *uint) (uint, bool) {
	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "contexto de negocio no encontrado"})
		return 0, false
	}
	if businessID > 0 {
		return businessID, true
	}

	if bodyBusinessID != nil && *bodyBusinessID > 0 {
		return *bodyBusinessID, true
	}
	if param := c.Query("business_id"); param != "" {
		if id, err := strconv.ParseUint(param, 10, 64); err == nil && id > 0 {
			return uint(id), true
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
	return 0, false
}

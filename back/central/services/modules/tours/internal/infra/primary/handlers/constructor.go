package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/tours/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Handler struct {
	uc  app.IUseCase
	log log.ILogger
}

func New(useCase app.IUseCase, logger log.ILogger) *Handler {
	return &Handler{uc: useCase, log: logger}
}

func (h *Handler) resolveIdentity(c *gin.Context) (uint, uint, bool) {
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "usuario no identificado"})
		return 0, 0, false
	}

	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "contexto de negocio no encontrado"})
		return 0, 0, false
	}

	if businessID == 0 {
		if param := c.Query("business_id"); param != "" {
			parsed, err := strconv.ParseUint(param, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "business_id invalido"})
				return 0, 0, false
			}
			businessID = uint(parsed)
		}
	}

	return userID, businessID, true
}

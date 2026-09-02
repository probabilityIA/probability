package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/app"
	dom "github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandlers interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type Handlers struct {
	uc  app.IUseCase
	log log.ILogger
}

func New(uc app.IUseCase, logger log.ILogger) IHandlers {
	return &Handlers{uc: uc, log: logger}
}

func (h *Handlers) requireSuperAdmin(c *gin.Context) (uint, bool) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return 0, false
	}
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return 0, false
	}
	return userID, true
}

func (h *Handlers) parseUintParam(c *gin.Context, key string) (uint, bool) {
	v, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil || v == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return uint(v), true
}

func (h *Handlers) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, dom.ErrSprintNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, dom.ErrNameRequired), errors.Is(err, dom.ErrInvalidStatus),
		errors.Is(err, dom.ErrInvalidStartDate), errors.Is(err, dom.ErrInvalidEndDate),
		errors.Is(err, dom.ErrInvalidDateRange), errors.Is(err, dom.ErrCreatorRequired):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/app"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/errors"
)

type Handlers struct {
	uc app.IUseCase
}

func New(uc app.IUseCase) *Handlers {
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

func parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func parseUintParam(c *gin.Context, name string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || value == 0 {
		return 0, false
	}
	return uint(value), true
}

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domainerrors.ErrGroupNotFound),
		errors.Is(err, domainerrors.ErrClientNotFound),
		errors.Is(err, domainerrors.ErrProductNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domainerrors.ErrGroupNameDuplicate):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, domainerrors.ErrGroupNameRequired),
		errors.Is(err, domainerrors.ErrInvalidPrice),
		errors.Is(err, domainerrors.ErrTargetRequired),
		errors.Is(err, domainerrors.ErrTargetAmbiguous),
		errors.Is(err, domainerrors.ErrNoClients):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

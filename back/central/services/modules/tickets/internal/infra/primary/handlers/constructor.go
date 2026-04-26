package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/app"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
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

func (h *Handlers) parseUint(s string) (uint, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

func (h *Handlers) parseUintParam(c *gin.Context, key string) (uint, bool) {
	id, err := h.parseUint(c.Param(key))
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return id, true
}

func (h *Handlers) requesterContext(c *gin.Context) (userID uint, businessID *uint, isSuperAdmin bool) {
	userID, _ = middleware.GetUserID(c)
	bid, _ := middleware.GetBusinessID(c)
	isSuperAdmin = middleware.IsSuperAdmin(c)
	if !isSuperAdmin && bid > 0 {
		businessID = &bid
	}
	if isSuperAdmin {
		if param := c.Query("business_id"); param != "" {
			if v, err := strconv.ParseUint(param, 10, 64); err == nil && v > 0 {
				u := uint(v)
				businessID = &u
			}
		}
	}
	return
}

func (h *Handlers) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, dom.ErrTicketNotFound), errors.Is(err, dom.ErrCommentNotFound), errors.Is(err, dom.ErrAttachmentNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, dom.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, dom.ErrInvalidStatus), errors.Is(err, dom.ErrInvalidPriority),
		errors.Is(err, dom.ErrInvalidType), errors.Is(err, dom.ErrInvalidSeverity),
		errors.Is(err, dom.ErrTitleRequired), errors.Is(err, dom.ErrDescriptionRequired),
		errors.Is(err, dom.ErrAssigneeNotFound):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (h *Handlers) splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

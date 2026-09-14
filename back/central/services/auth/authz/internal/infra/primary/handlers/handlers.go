package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/app"
	domainerrors "github.com/secamc93/probability/back/central/services/auth/authz/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/auth/authz/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Handlers struct {
	uc  app.IUseCase
	log log.ILogger
}

func New(uc app.IUseCase, logger log.ILogger) *Handlers {
	return &Handlers{uc: uc, log: logger}
}

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/auth/me/access", middleware.JWT(), h.GetMyAccess)
}

func (h *Handlers) GetMyAccess(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "sesion invalida"})
		return
	}
	tokenBusinessID, _ := middleware.GetBusinessIDFromContext(c)

	var requestedBusinessID uint
	if tokenBusinessID == 0 {
		if raw := c.Query("business_id"); raw != "" {
			if parsed, err := strconv.ParseUint(raw, 10, 64); err == nil {
				requestedBusinessID = uint(parsed)
			}
		}
	}

	access, err := h.uc.GetEffectiveAccess(c.Request.Context(), userID, tokenBusinessID, requestedBusinessID)
	if err != nil {
		if errors.Is(err, domainerrors.ErrNoBusinessRelation) {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "error": err.Error()})
			return
		}
		h.log.Error(c.Request.Context()).Err(err).Uint("user_id", userID).Msg("[authz] error resolviendo acceso efectivo")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "no se pudo resolver el acceso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response.FromAccess(access, h.uc.Navigation(access)),
	})
}

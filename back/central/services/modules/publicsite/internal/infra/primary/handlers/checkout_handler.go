package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/infra/primary/handlers/request"
)

func (h *Handlers) CreateCheckoutSession(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug es requerido"})
		return
	}

	var req request.CreateCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos invalidos: " + err.Error()})
		return
	}

	session, err := h.uc.CreateCheckoutSession(c.Request.Context(), slug, req.ToDTO(), h.optionalUserID(c))
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrBusinessNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrPublicSiteNotActive):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrEmptyCart), errors.Is(err, domainerrors.ErrInvalidCartItem):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, domainerrors.ErrOnlinePayNotReady):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": session})
}

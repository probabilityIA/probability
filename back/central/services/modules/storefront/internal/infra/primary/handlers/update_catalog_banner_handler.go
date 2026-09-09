package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/infra/primary/handlers/response"
)

func (h *Handlers) UpdateCatalogBanner(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido"})
		return
	}
	requesterUserID := c.GetUint("user_id")

	var req request.UpdateCatalogBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	banner, err := h.uc.UpdateCatalogBanner(c.Request.Context(), businessID, requesterUserID, req.Enabled, nil)
	if err != nil {
		respondCatalogBannerError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.CatalogBannerFromEntity(banner, h.getImageURLBase()))
}

func respondCatalogBannerError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domainerrors.ErrStorefrontNotActive):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, domainerrors.ErrRoleNotAllowed):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

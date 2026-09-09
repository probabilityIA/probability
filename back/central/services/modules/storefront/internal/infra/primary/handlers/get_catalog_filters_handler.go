package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/infra/primary/handlers/response"
)

func (h *Handlers) GetCatalogFilters(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido"})
		return
	}

	filters, err := h.uc.GetCatalogFilters(c.Request.Context(), businessID)
	if err != nil {
		if errors.Is(err, domainerrors.ErrStorefrontNotActive) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.CatalogFiltersFromEntity(filters, h.getImageURLBase()))
}

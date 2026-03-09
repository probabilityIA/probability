package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/infra/primary/handlers/response"
)

func (h *Handlers) GetBusinessPage(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug es requerido"})
		return
	}

	business, err := h.uc.GetBusinessPage(c.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, domainerrors.ErrBusinessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domainerrors.ErrPublicSiteNotActive) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	imageURLBase := h.getImageURLBase()

	// Also get featured products for the landing page
	featured, _ := h.uc.GetFeaturedProducts(c.Request.Context(), slug, 8)

	c.JSON(http.StatusOK, response.BusinessPageFromEntity(business, featured, imageURLBase))
}

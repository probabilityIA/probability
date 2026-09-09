package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/infra/primary/handlers/response"
)

func (h *Handlers) UploadCatalogBannerImage(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido"})
		return
	}
	requesterUserID := c.GetUint("user_id")

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "la imagen es requerida"})
		return
	}

	folder := fmt.Sprintf("storefront-banners/%d", businessID)
	relativePath, err := h.s3.UploadImage(c.Request.Context(), file, folder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	enabled := true
	banner, err := h.uc.UpdateCatalogBanner(c.Request.Context(), businessID, requesterUserID, &enabled, &relativePath)
	if err != nil {
		respondCatalogBannerError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.CatalogBannerFromEntity(banner, h.getImageURLBase()))
}

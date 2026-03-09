package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/infra/primary/handlers/response"
)

func (h *Handlers) GetProduct(c *gin.Context) {
	slug := c.Param("slug")
	productID := c.Param("id")

	if slug == "" || productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug y id del producto son requeridos"})
		return
	}

	product, err := h.uc.GetProduct(c.Request.Context(), slug, productID)
	if err != nil {
		if errors.Is(err, domainerrors.ErrBusinessNotFound) || errors.Is(err, domainerrors.ErrProductNotFound) {
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
	c.JSON(http.StatusOK, response.ProductFromEntity(product, imageURLBase))
}

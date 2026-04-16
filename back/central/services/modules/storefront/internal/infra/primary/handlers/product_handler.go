package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/infra/primary/handlers/response"
)

func (h *Handlers) GetProduct(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido"})
		return
	}

	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id del producto es requerido"})
		return
	}

	product, err := h.uc.GetProduct(c.Request.Context(), businessID, productID)
	if err != nil {
		if errors.Is(err, domainerrors.ErrStorefrontNotActive) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domainerrors.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	imageURLBase := h.getImageURLBase()
	c.JSON(http.StatusOK, response.ProductFromEntity(product, imageURLBase))
}

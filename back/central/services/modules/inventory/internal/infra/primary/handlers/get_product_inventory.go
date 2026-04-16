package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers/response"
)

func (h *handlers) GetProductInventory(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	productID := c.Param("productId")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id is required"})
		return
	}

	params := dtos.GetProductInventoryParams{
		ProductID:  productID,
		BusinessID: businessID,
	}

	levels, err := h.uc.GetProductInventory(c.Request.Context(), params)
	if err != nil {
		if errors.Is(err, domainerrors.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]response.InventoryLevelResponse, len(levels))
	for i, l := range levels {
		data[i] = response.InventoryLevelFromEntity(&l)
	}

	c.JSON(http.StatusOK, data)
}

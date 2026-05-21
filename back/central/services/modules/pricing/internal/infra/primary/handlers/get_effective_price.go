package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/infra/primary/handlers/response"
)

func (h *Handlers) GetEffectivePrice(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	productID := c.Query("product_id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id is required"})
		return
	}

	clientID := parseOptionalUint(c, "client_id")
	if clientID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id is required"})
		return
	}

	price, err := h.uc.GetEffectivePrice(c.Request.Context(), dtos.EffectivePriceParams{
		BusinessID: businessID,
		ProductID:  productID,
		ClientID:   *clientID,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.FromEffectivePrice(price))
}

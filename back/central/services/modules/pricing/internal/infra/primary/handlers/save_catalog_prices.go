package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/infra/primary/handlers/request"
)

func (h *Handlers) SaveCatalogPrices(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	var req request.SaveCatalogPricesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items := make([]dtos.SaveCatalogPriceItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = dtos.SaveCatalogPriceItem{
			ProductID: item.ProductID,
			Price:     item.Price,
		}
	}

	err := h.uc.SaveCatalogPrices(c.Request.Context(), dtos.SaveCatalogPricesDTO{
		Target: dtos.CatalogPriceTarget{
			BusinessID:    businessID,
			ClientGroupID: req.ClientGroupID,
			ClientID:      req.ClientID,
		},
		Items: items,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "catalog prices saved"})
}

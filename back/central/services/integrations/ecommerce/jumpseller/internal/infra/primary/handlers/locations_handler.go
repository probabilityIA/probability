package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

func (h *jumpsellerHandler) GetLocations(c *gin.Context) {
	integrationIDParam := c.Query("integration_id")
	if integrationIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}

	var bodyBusinessID *uint
	if raw := c.Query("business_id"); raw != "" {
		if parsed, err := strconv.ParseUint(raw, 10, 64); err == nil {
			value := uint(parsed)
			bodyBusinessID = &value
		}
	}

	businessID, ok := h.resolveBusinessID(c, bodyBusinessID)
	if !ok {
		return
	}

	info, err := h.useCase.GetLocations(c.Request.Context(), integrationIDParam, businessID)
	if err != nil {
		if errors.Is(err, domain.ErrIntegrationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	locations := make([]gin.H, 0, len(info.Locations))
	for _, location := range info.Locations {
		locations = append(locations, gin.H{
			"id":              location.ID,
			"name":            location.Name,
			"main":            location.Main,
			"is_stock_origin": location.IsStockOrigin,
			"pickup_point":    location.PickupPoint,
			"city":            location.City,
			"country":         location.Country,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":           true,
		"locations":         locations,
		"subscription_plan": info.SubscriptionPlan,
		"multi_location":    info.MultiLocation,
		"stock_origin_name": info.StockOriginName,
	})
}

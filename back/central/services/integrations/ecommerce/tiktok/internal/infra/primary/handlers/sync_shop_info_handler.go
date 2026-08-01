package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/domain"
)

type syncShopInfoRequest struct {
	IntegrationID uint  `json:"integration_id" binding:"required"`
	BusinessID    *uint `json:"business_id"`
}

func (h *tiktokHandler) SyncShopInfo(c *gin.Context) {
	var req syncShopInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}

	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "contexto de negocio no encontrado"})
		return
	}
	if businessID == 0 {
		if req.BusinessID == nil || *req.BusinessID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
			return
		}
		businessID = *req.BusinessID
	}

	snapshot, err := h.useCase.SyncShopInfo(c.Request.Context(), strconv.FormatUint(uint64(req.IntegrationID), 10), businessID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrIntegrationNotFound) || errors.Is(err, domain.ErrShopNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Shop info sincronizada y snapshot encolado",
		"data": gin.H{
			"shop_id":                snapshot.ShopID,
			"shop_name":              snapshot.ShopName,
			"shop_rating":            snapshot.ShopRating,
			"shop_review_count":      snapshot.ShopReviewCount,
			"item_sold_count":        snapshot.ItemSoldCount,
			"shop_performance_value": snapshot.ShopPerformanceValue,
		},
	})
}

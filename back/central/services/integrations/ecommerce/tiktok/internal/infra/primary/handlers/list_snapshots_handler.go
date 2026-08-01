package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *tiktokHandler) ListSnapshots(c *gin.Context) {
	integrationID, err := strconv.ParseUint(c.Query("integration_id"), 10, 64)
	if err != nil || integrationID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id es requerido"})
		return
	}

	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "contexto de negocio no encontrado"})
		return
	}
	if businessID == 0 {
		param, err := strconv.ParseUint(c.Query("business_id"), 10, 64)
		if err != nil || param == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
			return
		}
		businessID = uint(param)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	snapshots, total, err := h.useCase.ListSnapshots(c.Request.Context(), businessID, uint(integrationID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	data := make([]gin.H, 0, len(snapshots))
	for _, s := range snapshots {
		data = append(data, gin.H{
			"id":                     s.ID,
			"integration_id":         s.IntegrationID,
			"shop_id":                s.ShopID,
			"shop_name":              s.ShopName,
			"shop_rating":            s.ShopRating,
			"shop_review_count":      s.ShopReviewCount,
			"item_sold_count":        s.ItemSoldCount,
			"shop_performance_value": s.ShopPerformanceValue,
			"created_at":             s.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"data":        data,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": int(math.Ceil(float64(total) / float64(pageSize))),
	})
}

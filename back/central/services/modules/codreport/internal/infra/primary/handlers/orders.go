package handlers

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
)

func (h *Handlers) ListOrders(c *gin.Context) {
	businessID, err := resolveBusinessID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	start, end := parseDateRange(c)
	page, pageSize := parsePagination(c)

	f := dtos.OrdersFilter{
		BusinessID: businessID,
		StartDate:  start,
		EndDate:    end,
		Carrier:    strings.TrimSpace(c.Query("carrier")),
		Status:     strings.TrimSpace(c.Query("status")),
		Search:     strings.TrimSpace(c.Query("search")),
		Page:       page,
		PageSize:   pageSize,
	}
	if collected := c.Query("collected"); collected != "" {
		if v, errParse := strconv.ParseBool(collected); errParse == nil {
			f.Collected = &v
		}
	}
	if hasGuide := c.Query("has_guide"); hasGuide != "" {
		if v, errParse := strconv.ParseBool(hasGuide); errParse == nil {
			f.HasGuide = &v
		}
	}

	orders, total, err := h.uc.ListOrders(c.Request.Context(), f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener las ordenes contra entrega",
			"error":   err.Error(),
		})
		return
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(pageSize)))
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Ordenes contra entrega obtenidas exitosamente",
		"data":        mapOrders(orders),
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

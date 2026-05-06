package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (h *Handlers) StatsByGeozone(c *gin.Context) {
	businessID := c.GetUint("business_id")
	if businessID == 0 {
		if v := c.Query("business_id"); v != "" {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil {
				businessID = uint(id)
			}
		}
	}
	if businessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "business_id is required"})
		return
	}

	filter := domain.ShipmentStatsFilter{
		BusinessID: businessID,
		Carrier:    c.Query("carrier"),
		Type:       c.Query("type"),
	}

	if v := c.Query("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.From = &t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.To = &t
		}
	}
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.Limit = n
		}
	}

	stats, err := h.uc.GetStatsByGeozone(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

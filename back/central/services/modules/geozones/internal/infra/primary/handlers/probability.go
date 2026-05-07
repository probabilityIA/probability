package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/geozones/internal/domain/dtos"
)

func (h *Handlers) Probability(c *gin.Context) {
	if h.probabilityUC == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "error": "probability use case not initialized"})
		return
	}
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id is required"})
		return
	}
	req := dtos.ProbabilityRequest{
		BusinessID: businessID,
		Carrier:    c.Query("carrier"),
		OrderID:    c.Query("order_id"),
	}
	if latStr := c.Query("lat"); latStr != "" {
		if v, err := strconv.ParseFloat(latStr, 64); err == nil {
			req.Lat = &v
		}
	}
	if lngStr := c.Query("lng"); lngStr != "" {
		if v, err := strconv.ParseFloat(lngStr, 64); err == nil {
			req.Lng = &v
		}
	}
	if req.OrderID == "" && (req.Lat == nil || req.Lng == nil) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "order_id or (lat, lng) required"})
		return
	}
	res, err := h.probabilityUC.GetProbability(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}

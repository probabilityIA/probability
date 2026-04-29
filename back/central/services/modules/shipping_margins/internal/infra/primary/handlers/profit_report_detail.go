package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/dtos"
)

func (h *Handlers) ProfitReportDetail(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	carrier := c.Query("carrier")
	if carrier == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "carrier is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	params := dtos.ProfitReportDetailParams{
		BusinessID: businessID,
		Carrier:    carrier,
		Page:       page,
		PageSize:   pageSize,
	}

	if from := c.Query("from"); from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			params.From = &t
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from (YYYY-MM-DD)"})
			return
		}
	}
	if to := c.Query("to"); to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
			tEnd := t.Add(24 * time.Hour)
			params.To = &tEnd
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to (YYYY-MM-DD)"})
			return
		}
	}

	resp, err := h.uc.ProfitReportDetail(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

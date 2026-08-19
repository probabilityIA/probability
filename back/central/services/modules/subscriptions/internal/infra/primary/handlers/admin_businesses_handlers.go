package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListAdminBusinesses(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	status := c.Query("status")

	rows, total, err := h.uc.ListAdminBusinesses(c.Request.Context(), page, pageSize, search, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list businesses"})
		return
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	c.JSON(http.StatusOK, gin.H{
		"data":        response.FromAdminBusinessRows(rows),
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

func (h *Handlers) GetAdminKPIs(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}

	kpis, err := h.uc.GetAdminKPIs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch kpis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response.FromAdminKPIs(kpis)})
}

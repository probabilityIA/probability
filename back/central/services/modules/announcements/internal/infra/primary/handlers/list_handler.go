package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/infra/primary/handlers/mappers"
)

func (h *handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	params := dtos.ListAnnouncementsParams{
		Status:   c.Query("status"),
		Search:   c.Query("search"),
		Page:     page,
		PageSize: pageSize,
	}

	if catIDStr := c.Query("category_id"); catIDStr != "" {
		if catID, err := strconv.ParseUint(catIDStr, 10, 64); err == nil {
			id := uint(catID)
			params.CategoryID = &id
		}
	}

	items, total, err := h.uc.ListAnnouncements(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	resp := mappers.EntityListToResponse(items, total, page, pageSize)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp.Data, "total": resp.Total, "page": resp.Page, "page_size": resp.PageSize, "total_pages": resp.TotalPages})
}

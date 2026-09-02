package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/infra/primary/handlers/response"
)

func (h *Handlers) List(c *gin.Context) {
	if _, ok := h.requireSuperAdmin(c); !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	params := dtos.ListSprintsParams{
		Page:     page,
		PageSize: pageSize,
		Status:   c.Query("status"),
	}

	items, total, err := h.uc.List(c.Request.Context(), params)
	if err != nil {
		h.handleError(c, err)
		return
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}

	out := make([]response.SprintResponse, 0, len(items))
	for i := range items {
		out = append(out, response.FromSprint(&items[i]))
	}
	totalPages := total / int64(params.PageSize)
	if total%int64(params.PageSize) > 0 {
		totalPages++
	}
	c.JSON(http.StatusOK, gin.H{
		"data":        out,
		"total":       total,
		"page":        params.Page,
		"page_size":   params.PageSize,
		"total_pages": totalPages,
	})
}

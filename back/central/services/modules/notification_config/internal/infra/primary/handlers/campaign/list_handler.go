package campaign

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) List(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
		return
	}

	page, pageSize := paging(c)

	items, total, err := h.useCase.List(c.Request.Context(), businessID, c.Query("status"), page, pageSize)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error listando campanas")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "error interno"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"data":        items,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages(total, pageSize),
	})
}

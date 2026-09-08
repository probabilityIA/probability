package scheduled_rule

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

	items, total, err := h.useCase.ListRules(c.Request.Context(), businessID, page, pageSize)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error listando reglas programadas")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "error interno"})
		return
	}

	totalPages := 0
	if pageSize > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"data":        items,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

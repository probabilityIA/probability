package scheduled_rule

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) ListRuns(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
		return
	}

	id, ok := parseID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}

	page, pageSize := paging(c)

	items, total, err := h.useCase.ListRuns(c.Request.Context(), id, businessID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

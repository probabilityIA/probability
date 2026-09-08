package whatsapp_template

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *handler) List(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	items, total, err := h.useCase.List(c.Request.Context(), businessID, c.Query("scope"), c.Query("status"), page, pageSize)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error listando plantillas de WhatsApp")
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

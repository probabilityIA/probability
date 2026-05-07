package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) ListSKUs(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	prefix := strings.TrimSpace(c.Query("prefix"))

	skus, err := h.uc.ListSKUs(c.Request.Context(), businessID, prefix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener SKUs",
			"error":   err.Error(),
		})
		return
	}

	if skus == nil {
		skus = []string{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "SKUs obtenidos exitosamente",
		"data":    skus,
	})
}

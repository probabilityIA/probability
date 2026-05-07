package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetNextSKUBatch(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	prefix := strings.TrimSpace(c.Query("prefix"))
	if prefix == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "El parámetro 'prefix' es requerido",
			"error":   "prefix is required",
		})
		return
	}

	countStr := c.DefaultQuery("count", "1")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "El parámetro 'count' debe ser un número mayor a 0",
			"error":   "invalid count parameter",
		})
		return
	}

	if count > 1000 {
		count = 1000
	}

	nextSKUs, err := h.uc.GetNextSKUBatch(c.Request.Context(), businessID, prefix, count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener siguientes SKUs",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Siguientes SKUs obtenidos exitosamente",
		"data":    nextSKUs,
	})
}

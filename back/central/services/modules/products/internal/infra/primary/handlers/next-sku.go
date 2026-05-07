package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetNextSKU(c *gin.Context) {
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

	nextSKU, err := h.uc.GetNextSKU(c.Request.Context(), businessID, prefix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener siguiente SKU",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Siguiente SKU obtenido exitosamente",
		"data":    nextSKU,
	})
}

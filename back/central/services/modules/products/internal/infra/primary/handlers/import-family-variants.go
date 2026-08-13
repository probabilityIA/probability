package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

func (h *Handlers) ImportFamilyVariants(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	familyID, ok := parseFamilyID(c)
	if !ok {
		return
	}

	if _, err := h.uc.GetProductFamilyByID(c.Request.Context(), businessID, familyID); err != nil {
		if errors.Is(err, domain.ErrProductFamilyNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Familia no encontrada", "error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al validar familia", "error": err.Error()})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No se recibió un archivo válido",
			"error":   err.Error(),
		})
		return
	}
	defer file.Close()

	rows, err := parseDimensionRows(file, false)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Error al leer el archivo: %s", err.Error()),
			"error":   err.Error(),
		})
		return
	}

	preview, err := h.uc.MatchFamilyImportRows(c.Request.Context(), businessID, rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al procesar el archivo",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("%d productos encontrados, %d sin coincidencia", len(preview.Matched), len(preview.NotFound)),
		"data":    preview,
	})
}

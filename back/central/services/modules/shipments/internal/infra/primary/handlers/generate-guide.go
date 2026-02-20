package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (h *Handlers) GenerateGuide(c *gin.Context) {
	var req domain.EnvioClickQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.envioClickUC.GenerateGuide(c.Request.Context(), req)
	if err != nil {
		// Detectar si es un error de validación (422) para mostrarlo amigablemente
		statusCode := http.StatusInternalServerError
		errMsg := err.Error()
		if strings.Contains(strings.ToLower(errMsg), "error:") ||
			strings.Contains(strings.ToLower(errMsg), "inválido") ||
			strings.Contains(strings.ToLower(errMsg), "unprocessed entity") ||
			strings.Contains(strings.ToLower(errMsg), "falta") {
			statusCode = http.StatusUnprocessableEntity
		}

		c.JSON(statusCode, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, resp)
}

package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
)

func (h *Handlers) Summary(c *gin.Context) {
	businessID, err := resolveBusinessID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	start, end := parseDateRange(c)
	res, err := h.uc.Summary(c.Request.Context(), dtos.ReportFilter{
		BusinessID: businessID,
		StartDate:  start,
		EndDate:    end,
		Carrier:    strings.TrimSpace(c.Query("carrier")),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener el resumen de recaudo",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Resumen de recaudo obtenido exitosamente",
		"data":    mapSummary(res),
	})
}

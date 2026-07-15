package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
)

func (h *Handlers) SelectableCutOrders(c *gin.Context) {
	businessID, err := resolveBusinessID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	start, err1 := time.Parse("2006-01-02", c.Query("period_start"))
	end, err2 := time.Parse("2006-01-02", c.Query("period_end"))
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Fechas del periodo invalidas"})
		return
	}
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)

	orders, err := h.uc.SelectableOrders(c.Request.Context(), dtos.SelectableOrdersFilter{
		BusinessID:  businessID,
		PeriodStart: start,
		PeriodEnd:   end,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener las ordenes de la semana",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ordenes seleccionables obtenidas exitosamente",
		"data":    mapOrders(orders),
	})
}

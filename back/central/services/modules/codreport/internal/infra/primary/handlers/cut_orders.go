package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) DeleteCut(c *gin.Context) {
	if !isAdminUser(c) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Solo un administrador puede eliminar cortes de pago"})
		return
	}

	businessID, err := resolveBusinessID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	cutID, err := strconv.ParseUint(c.Query("cut_id"), 10, 64)
	if err != nil || cutID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "cut_id invalido"})
		return
	}

	if err := h.uc.DeleteCut(c.Request.Context(), businessID, uint(cutID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al eliminar el corte de pago",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Corte de pago eliminado exitosamente"})
}

func (h *Handlers) CutOrders(c *gin.Context) {
	businessID, err := resolveBusinessID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	cutID, err := strconv.ParseUint(c.Query("cut_id"), 10, 64)
	if err != nil || cutID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "cut_id invalido"})
		return
	}

	orders, err := h.uc.CutOrders(c.Request.Context(), businessID, uint(cutID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener las ordenes del corte",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ordenes del corte obtenidas exitosamente",
		"data":    mapOrders(orders),
	})
}

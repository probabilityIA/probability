package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) AcceptShipmentOverage(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	if err := h.uc.AcceptShipmentOverage(c.Request.Context(), businessID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cargo extra por excedente de envios aceptado. Se facturara al cierre del periodo.",
	})
}

func (h *Handlers) PayShipmentOverage(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	if err := h.uc.PayShipmentOverage(c.Request.Context(), businessID, c.GetUint("user_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cargo de excedente pagado. El plan gratuito quedo renovado.",
	})
}

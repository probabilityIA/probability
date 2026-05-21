package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) ListCuts(c *gin.Context) {
	businessID, err := resolveBusinessID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	admin := isAdminUser(c)
	cuts, err := h.uc.ListCuts(c.Request.Context(), businessID, admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener los cortes de pago",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Cortes de pago obtenidos exitosamente",
		"data":        mapCuts(cuts),
		"can_confirm": admin,
	})
}

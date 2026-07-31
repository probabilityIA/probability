package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
)

func (h *Handlers) GetMyPrices(c *gin.Context) {
	slug := c.Param("slug")
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "sesion no valida"})
		return
	}

	prices, err := h.uc.GetMyPrices(c.Request.Context(), slug, userID)
	if err != nil {
		if errors.Is(err, domainerrors.ErrBusinessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "negocio no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "error consultando precios"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"prices": prices}})
}

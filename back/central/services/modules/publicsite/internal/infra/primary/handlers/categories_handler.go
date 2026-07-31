package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
)

func (h *Handlers) GetCategories(c *gin.Context) {
	slug := c.Param("slug")

	categories, err := h.uc.GetCategories(c.Request.Context(), slug)
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrBusinessNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "negocio no encontrado"})
		case errors.Is(err, domainerrors.ErrPublicSiteNotActive):
			c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "tienda no activa"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "error consultando categorias"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": categories})
}

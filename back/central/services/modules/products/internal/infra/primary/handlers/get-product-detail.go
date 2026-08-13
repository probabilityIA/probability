package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

func (h *Handlers) GetProductDetail(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		h.respondBusinessIDRequired(c)
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de producto invalido",
			"error":   "El ID es requerido",
		})
		return
	}

	detail, err := h.uc.GetProductDetail(c.Request.Context(), businessID, id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Producto no encontrado",
				"error":   err.Error(),
			})
			return
		}
		h.log.Error(c.Request.Context()).Err(err).
			Uint("business_id", businessID).Str("product_id", id).
			Msg("Error al obtener el detalle del producto")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener el detalle del producto",
			"error":   err.Error(),
		})
		return
	}

	detail.ImageURL = h.buildImageURL(detail.ImageURL)
	if detail.Family != nil {
		detail.Family.ImageURL = h.buildImageURL(detail.Family.ImageURL)
	}
	for i := range detail.Channels {
		detail.Channels[i].IntegrationImageURL = h.buildImageURL(detail.Channels[i].IntegrationImageURL)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detalle del producto obtenido exitosamente",
		"data":    detail,
	})
}

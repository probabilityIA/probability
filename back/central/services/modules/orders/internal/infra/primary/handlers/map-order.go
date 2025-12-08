package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
)

// MapAndSaveOrder godoc
// @Summary      Mapear y guardar orden canónica
// @Description  Recibe una orden en formato canónico (después de mapeo) y la guarda en todas las tablas relacionadas
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        order  body      domain.CanonicalOrderDTO  true  "Orden en formato canónico"
// @Security     BearerAuth
// @Success      201  {object}  domain.OrderResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders/map [post]
func (h *Handlers) MapAndSaveOrder(c *gin.Context) {
	var req domain.CanonicalOrderDTO

	// Validar el request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

	// Llamar al caso de uso de mapeo
	order, err := h.orderMapping.MapAndSaveOrder(c.Request.Context(), &req)
	if err != nil {
		// Verificar si es un error de duplicado
		// Verificar si es un error de duplicado
		if errors.Is(err, domain.ErrOrderAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Orden ya existe",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al mapear y guardar orden",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Orden mapeada y guardada exitosamente",
		"data":    order,
	})
}

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

// UpdateShipment godoc
// @Summary      Actualizar envío
// @Description  Actualiza un envío existente
// @Tags         Shipments
// @Accept       json
// @Produce      json
// @Param        id        path      int                        true  "ID del envío"
// @Param        shipment  body      domain.UpdateShipmentRequest true  "Datos a actualizar"
// @Security     BearerAuth
// @Success      200  {object}  domain.ShipmentResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /shipments/{id} [put]
func (h *Handlers) UpdateShipment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de envío inválido",
			"error":   "El ID debe ser un número válido",
		})
		return
	}

	var req domain.UpdateShipmentRequest

	// Validar el request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

	// Llamar al caso de uso
	shipment, err := h.uc.UpdateShipment(c.Request.Context(), uint(id), &req)
	if err != nil {
		if err == domain.ErrShipmentNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Envío no encontrado",
				"error":   err.Error(),
			})
			return
		}

		if err == domain.ErrShipmentAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Envío con este número de tracking ya existe",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al actualizar envío",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Envío actualizado exitosamente",
		"data":    shipment,
	})
}


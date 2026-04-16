package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

// CreateShipment godoc
// @Summary      Crear envío
// @Description  Crea un nuevo envío en el sistema
// @Tags         Shipments
// @Accept       json
// @Produce      json
// @Param        shipment  body      domain.CreateShipmentRequest  true  "Datos del envío"
// @Security     BearerAuth
// @Success      201  {object}  domain.ShipmentResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /shipments [post]
func (h *Handlers) CreateShipment(c *gin.Context) {
	var req domain.CreateShipmentRequest

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
	shipment, err := h.uc.CreateShipment(c.Request.Context(), &req)
	if err != nil {
		// Verificar si es un error de duplicado
		if err == domain.ErrShipmentAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Envío ya existe",
				"error":   err.Error(),
			})
			return
		}

		if err == domain.ErrOrderIDRequired {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Order ID es requerido",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al crear envío",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Envío creado exitosamente",
		"data":    shipment,
	})
}


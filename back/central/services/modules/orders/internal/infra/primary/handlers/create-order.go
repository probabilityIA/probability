package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
)

// CreateOrder godoc
// @Summary      Crear orden
// @Description  Crea una nueva orden en el sistema
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        order  body      dtos.CreateOrderRequest  true  "Datos de la orden"
// @Security     BearerAuth
// @Success      201  {object}  dtos.OrderResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders [post]
func (h *Handlers) CreateOrder(c *gin.Context) {
	var req dtos.CreateOrderRequest

	// Validar el request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

	h.logger.Info(c.Request.Context()).
		Str("platform", req.Platform).
		Str("external_id", req.ExternalID).
		Uint("integration_id", req.IntegrationID).
		Str("customer_name", req.CustomerName).
		Str("customer_email", req.CustomerEmail).
		Str("customer_phone", req.CustomerPhone).
		Float64("total_amount", req.TotalAmount).
		Str("currency", req.Currency).
		Int("items_count", len(req.Items)).
		Interface("order_request", req).
		Msg("📥 CreateOrder request recibido")

	// Validaciones adicionales para prevenir órdenes mal formadas
	// Para órdenes manuales, el backend puede generar el external_id y usar una integración por defecto
	if req.Platform != "manual" {
		if req.ExternalID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "external_id es requerido",
			})
			return
		}
		if req.IntegrationID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "integration_id es requerido",
			})
			return
		}
	}

	// Llamar al caso de uso centralizado (pasa por MapAndSaveOrder para score, status mapping, etc.)
	order, err := h.createUC.CreateManualOrder(c.Request.Context(), &req)
	if err != nil {
		// Verificar si es un error de duplicado
		if err.Error() == "order with this external_id already exists for this integration" {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Orden ya existe",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al crear orden",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Orden creada exitosamente",
		"data":    order,
	})
}

package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/test/internal/domain"
)

// GenerateOrders godoc
// @Summary      Generar órdenes aleatorias
// @Description  Genera órdenes canónicas aleatorias y las publica a la cola de RabbitMQ para que el consumidor las procese
// @Tags         Test
// @Accept       json
// @Produce      json
// @Param        request  body      domain.GenerateOrderRequest  true  "Configuración para generar órdenes"
// @Security     BearerAuth
// @Success      200  {object}  domain.GenerateOrderResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /test/generate-orders [post]
func (h *Handlers) GenerateOrders(c *gin.Context) {
	var req domain.GenerateOrderRequest

	// Validar el request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

	// Forzar business_id=7 e integration_id según la plataforma
	if req.Platform == "shopify" {
		businessID := uint(7)
		req.BusinessID = &businessID
		req.IntegrationID = 1
	} else if req.Platform == "meli" || req.Platform == "mercado_libre" {
		businessID := uint(7)
		req.BusinessID = &businessID
		req.IntegrationID = 5
	} else if req.Platform == "woocommerce" || req.Platform == "woo" {
		businessID := uint(7)
		req.BusinessID = &businessID
		req.IntegrationID = 6
	}

	// Llamar al caso de uso
	response, err := h.uc.GenerateAndPublishOrders(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al generar y publicar órdenes",
			"error":   err.Error(),
			"data":    response, // Incluir datos parciales si hay errores
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Órdenes generadas y publicadas exitosamente: %d generadas, %d publicadas, %d fallidas", response.Generated, response.Published, response.Failed),
		"data":    response,
	})
}

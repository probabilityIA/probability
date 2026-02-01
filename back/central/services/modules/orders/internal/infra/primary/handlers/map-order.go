package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/request"
)

// MapAndSaveOrder godoc
// @Summary      Mapear y guardar orden canónica
// @Description  Recibe una orden en formato canónico (después de mapeo) y la guarda en todas las tablas relacionadas
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        order  body      request.MapOrder  true  "Orden en formato de lógica de negocio"
// @Security     BearerAuth
// @Success      201  {object}  response.Order
// @Failure      400  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders/map [post]
func (h *Handlers) MapAndSaveOrder(c *gin.Context) {
	var req request.MapOrder // ✅ DTO HTTP con tags + datatypes.JSON

	// Validar el request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

	// Validaciones adicionales para prevenir órdenes mal formadas
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
	if req.BusinessID == nil || *req.BusinessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "business_id es requerido",
		})
		return
	}

	// ✅ Convertir HTTP request → Domain DTO (datatypes.JSON → []byte)
	domainReq := mappers.MapOrderRequestToDomain(&req)

	// Llamar al caso de uso de mapeo con DTO de dominio (SIN tags)
	domainResp, err := h.orderMapping.MapAndSaveOrder(c.Request.Context(), domainReq)
	if err != nil {
		// Verificar si es un error de duplicado
		if errors.Is(err, domainerrors.ErrOrderAlreadyExists) {
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

	// ✅ Convertir Domain response → HTTP response ([]byte → datatypes.JSON)
	httpResp := mappers.OrderToResponse(domainResp)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Orden mapeada y guardada exitosamente",
		"data":    httpResp,
	})
}

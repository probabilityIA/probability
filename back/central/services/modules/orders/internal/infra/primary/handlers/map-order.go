package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/request"
)

func (h *Handlers) MapAndSaveOrder(c *gin.Context) {
	var req request.MapOrder

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

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

	tokenBusinessID, _ := c.Get("business_id")
	if requesterBusinessID, _ := tokenBusinessID.(uint); requesterBusinessID > 0 {
		req.BusinessID = &requesterBusinessID
	}
	if req.BusinessID == nil || *req.BusinessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "business_id es requerido",
		})
		return
	}

	domainReq := mappers.MapOrderRequestToDomain(&req)

	domainResp, err := h.createUC.MapAndSaveOrder(c.Request.Context(), domainReq)
	if err != nil {
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

	httpResp := mappers.OrderToResponse(domainResp)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Orden mapeada y guardada exitosamente",
		"data":    httpResp,
	})
}

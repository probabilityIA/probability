package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/infra/primary/handlers/request"
)

func (h *Handlers) Apply(c *gin.Context) {
	var body request.ApplyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "cuerpo invalido: " + err.Error()})
		return
	}

	businessID, ok := h.resolveBusinessID(c, nil)
	if !ok {
		return
	}

	result, err := h.useCase.Apply(c.Request.Context(), dtos.ApplyCommand{
		BusinessID:    businessID,
		IntegrationID: body.IntegrationID,
		ExternalIDs:   body.ExternalIDs,
	})
	if err != nil {
		h.logger.Error(c.Request.Context()).Err(err).
			Uint("integration_id", body.IntegrationID).
			Uint("business_id", businessID).
			Msg("No se pudieron importar las ordenes del canal")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

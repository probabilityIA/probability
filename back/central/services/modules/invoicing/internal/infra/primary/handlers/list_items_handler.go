package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
)

// listItemsRequest cuerpo opcional de la solicitud de comparación de ítems
type listItemsRequest struct {
	BusinessID *uint `json:"business_id,omitempty"` // solo super admin
}

// ListItems inicia una comparación asíncrona de ítems del proveedor vs productos del sistema.
// Retorna 202 con un correlation_id; el resultado llega por SSE con evento "invoice.list_items_ready".
func (h *handler) ListItems(c *gin.Context) {
	// Extraer businessID del JWT
	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "business context not found"})
		return
	}

	if businessID == 0 {
		// Super admin: business_id debe venir en el body
		var req listItemsRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.BusinessID == nil || *req.BusinessID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required for super admin"})
			return
		}
		businessID = *req.BusinessID
	}

	dto := &dtos.ListItemsRequestDTO{
		BusinessID: businessID,
	}

	correlationID, err := h.useCase.RequestListItems(c.Request.Context(), dto)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Msg("Failed to start items comparison")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"correlation_id": correlationID,
		"message":        "Comparación de ítems iniciada. Recibirás el resultado por SSE.",
	})
}

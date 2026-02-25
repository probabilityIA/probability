package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
)

// compareRequest cuerpo de la solicitud de comparación
type compareRequest struct {
	DateFrom   string `json:"date_from" binding:"required"`
	DateTo     string `json:"date_to" binding:"required"`
	BusinessID *uint  `json:"business_id,omitempty"` // solo super admin
}

// CompareInvoices inicia una comparación asíncrona de facturas entre el sistema y el proveedor.
// Retorna 202 con un correlation_id; el resultado llega por SSE con evento "invoice.compare_ready".
func (h *handler) CompareInvoices(c *gin.Context) {
	var req compareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date_from and date_to are required (YYYY-MM-DD format)"})
		return
	}

	// Extraer businessID del JWT
	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "business context not found"})
		return
	}

	if businessID == 0 {
		// Super admin: business_id debe venir en el body
		if req.BusinessID == nil || *req.BusinessID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required for super admin"})
			return
		}
		businessID = *req.BusinessID
	}

	dto := &dtos.CompareRequestDTO{
		DateFrom:   req.DateFrom,
		DateTo:     req.DateTo,
		BusinessID: businessID,
	}

	correlationID, err := h.useCase.RequestComparison(c.Request.Context(), dto)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Msg("Failed to start invoice comparison")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"correlation_id": correlationID,
		"message":        "Comparación iniciada. Recibirás el resultado por SSE.",
	})
}

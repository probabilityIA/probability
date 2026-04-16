package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
)

// listBankAccountsRequest cuerpo opcional de la solicitud de cuentas bancarias
type listBankAccountsRequest struct {
	BusinessID *uint `json:"business_id,omitempty"` // solo super admin
}

// ListBankAccounts inicia una solicitud asíncrona de cuentas bancarias del proveedor.
// Retorna 202 con un correlation_id; el resultado llega por SSE con evento "invoice.list_bank_accounts_ready".
func (h *handler) ListBankAccounts(c *gin.Context) {
	// Extraer businessID del JWT
	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "business context not found"})
		return
	}

	if businessID == 0 {
		// Super admin: business_id debe venir en el body
		var req listBankAccountsRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.BusinessID == nil || *req.BusinessID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required for super admin"})
			return
		}
		businessID = *req.BusinessID
	}

	dto := &dtos.ListBankAccountsRequestDTO{
		BusinessID: businessID,
	}

	correlationID, err := h.useCase.RequestListBankAccounts(c.Request.Context(), dto)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Msg("Failed to start bank accounts request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"correlation_id": correlationID,
		"message":        "Solicitud de cuentas bancarias iniciada. Recibirás el resultado por SSE.",
	})
}

// GetListBankAccountsResult retorna el resultado de cuentas bancarias almacenado en Redis.
// Retorna 404 si el resultado aún no está listo o expiró.
func (h *handler) GetListBankAccountsResult(c *gin.Context) {
	correlationID := c.Param("correlationId")
	if correlationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "correlationId is required"})
		return
	}

	result, err := h.useCase.GetListBankAccountsResult(c.Request.Context(), correlationID)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).
			Str("correlation_id", correlationID).
			Msg("Failed to get bank accounts result")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve bank accounts result"})
		return
	}

	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "bank accounts result not found",
			"message": "El resultado aún no está listo o ya expiró (TTL 5 min).",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

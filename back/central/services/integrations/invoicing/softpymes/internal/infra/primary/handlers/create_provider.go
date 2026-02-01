package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/handlers/response"
)

// CreateProvider crea un nuevo proveedor de facturación Softpymes
func (h *handler) CreateProvider(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.CreateProvider
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(ctx).Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	h.log.Info(ctx).
		Str("name", req.Name).
		Str("provider_type", req.ProviderTypeCode).
		Uint("business_id", req.BusinessID).
		Msg("Creating Softpymes provider")

	dto := mappers.CreateProviderRequestToDTO(&req)

	provider, err := h.useCase.CreateProvider(ctx, dto)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("Failed to create provider")
		c.JSON(http.StatusInternalServerError, response.Error{
			Error:   "provider_creation_failed",
			Message: err.Error(),
		})
		return
	}

	resp := mappers.ProviderToResponse(provider, req.ProviderTypeCode)

	h.log.Info(ctx).
		Uint("provider_id", provider.ID).
		Str("name", provider.Name).
		Msg("Provider created successfully")

	c.JSON(http.StatusCreated, resp)
}

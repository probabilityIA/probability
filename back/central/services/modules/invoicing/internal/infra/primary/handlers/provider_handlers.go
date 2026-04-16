package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// CreateProvider crea un nuevo proveedor de facturación
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
		Msg("Creating invoicing provider")

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

	resp := mappers.ProviderToResponse(provider)

	h.log.Info(ctx).
		Uint("provider_id", provider.ID).
		Str("name", provider.Name).
		Msg("Provider created successfully")

	c.JSON(http.StatusCreated, resp)
}

// ListProviders lista proveedores de facturación
func (h *handler) ListProviders(c *gin.Context) {
	ctx := c.Request.Context()

	var businessID uint
	if businessIDStr := c.Query("business_id"); businessIDStr != "" {
		id, _ := strconv.ParseUint(businessIDStr, 10, 32)
		businessID = uint(id)
	}

	page := 1
	if p := c.Query("page"); p != "" {
		page, _ = strconv.Atoi(p)
	}

	pageSize := 20
	if ps := c.Query("page_size"); ps != "" {
		pageSize, _ = strconv.Atoi(ps)
	}

	if pageSize > 100 {
		pageSize = 100
	}

	h.log.Debug(ctx).Uint("business_id", businessID).Msg("Listing providers")

	providers, err := h.useCase.ListProviders(ctx, businessID)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("Failed to list providers")
		c.JSON(http.StatusInternalServerError, response.Error{
			Error:   "list_providers_failed",
			Message: err.Error(),
		})
		return
	}

	resp := mappers.ProvidersToResponse(providers, int64(len(providers)), page, pageSize)
	c.JSON(http.StatusOK, resp)
}

// GetProvider obtiene un proveedor por ID
func (h *handler) GetProvider(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_id",
			Message: "Invalid provider ID",
		})
		return
	}

	h.log.Debug(ctx).Uint("provider_id", uint(id)).Msg("Getting provider")

	provider, err := h.useCase.GetProvider(ctx, uint(id))
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("provider_id", uint(id)).Msg("Failed to get provider")
		c.JSON(http.StatusNotFound, response.Error{
			Error:   "provider_not_found",
			Message: err.Error(),
		})
		return
	}

	resp := mappers.ProviderToResponse(provider)
	c.JSON(http.StatusOK, resp)
}

// UpdateProvider actualiza un proveedor
func (h *handler) UpdateProvider(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_id",
			Message: "Invalid provider ID",
		})
		return
	}

	var req request.UpdateProvider
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(ctx).Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	h.log.Info(ctx).Uint("provider_id", uint(id)).Msg("Updating provider")

	dto := mappers.UpdateProviderRequestToDTO(&req)

	err = h.useCase.UpdateProvider(ctx, uint(id), dto)
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("provider_id", uint(id)).Msg("Failed to update provider")
		c.JSON(http.StatusInternalServerError, response.Error{
			Error:   "provider_update_failed",
			Message: err.Error(),
		})
		return
	}

	// Get updated provider to return
	provider, err := h.useCase.GetProvider(ctx, uint(id))
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("provider_id", uint(id)).Msg("Failed to get updated provider")
		c.JSON(http.StatusInternalServerError, response.Error{
			Error:   "get_provider_failed",
			Message: err.Error(),
		})
		return
	}

	resp := mappers.ProviderToResponse(provider)

	h.log.Info(ctx).Uint("provider_id", provider.ID).Msg("Provider updated successfully")

	c.JSON(http.StatusOK, resp)
}

// TestProvider prueba la conexión con un proveedor
func (h *handler) TestProvider(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_id",
			Message: "Invalid provider ID",
		})
		return
	}

	h.log.Info(ctx).Uint("provider_id", uint(id)).Msg("Testing provider connection")

	err = h.useCase.TestProviderConnection(ctx, uint(id))
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("provider_id", uint(id)).Msg("Provider test failed")
		c.JSON(http.StatusOK, response.TestProviderResult{
			Success: false,
			Message: "Connection test failed",
			Error:   err.Error(),
		})
		return
	}

	h.log.Info(ctx).Uint("provider_id", uint(id)).Msg("Provider test successful")

	c.JSON(http.StatusOK, response.TestProviderResult{
		Success: true,
		Message: "Connection test successful",
	})
}

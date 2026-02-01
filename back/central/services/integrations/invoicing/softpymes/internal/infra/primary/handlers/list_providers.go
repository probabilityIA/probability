package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/handlers/response"
)

// ListProviders lista proveedores de facturación Softpymes
func (h *handler) ListProviders(c *gin.Context) {
	ctx := c.Request.Context()

	// Construir filtros desde query params
	filters := &dtos.ProviderFiltersDTO{
		Limit:  20,
		Offset: 0,
	}

	if businessIDStr := c.Query("business_id"); businessIDStr != "" {
		id, _ := strconv.ParseUint(businessIDStr, 10, 32)
		businessID := uint(id)
		filters.BusinessID = &businessID
	}

	if providerTypeCode := c.Query("provider_type_code"); providerTypeCode != "" {
		filters.ProviderTypeCode = &providerTypeCode
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filters.IsActive = &isActive
	}

	if isDefaultStr := c.Query("is_default"); isDefaultStr != "" {
		isDefault := isDefaultStr == "true"
		filters.IsDefault = &isDefault
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

	filters.Limit = pageSize
	filters.Offset = (page - 1) * pageSize

	h.log.Debug(ctx).Interface("filters", filters).Msg("Listing providers")

	providers, err := h.useCase.ListProviders(ctx, filters)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("Failed to list providers")
		c.JSON(http.StatusInternalServerError, response.Error{
			Error:   "list_providers_failed",
			Message: err.Error(),
		})
		return
	}

	// TODO: Obtener tipos de proveedores para el mapeo
	providerTypes := make(map[uint]string)
	for _, p := range providers {
		providerTypes[p.ProviderTypeID] = "softpymes"
	}

	resp := mappers.ProvidersToResponse(providers, providerTypes, int64(len(providers)), page, pageSize)
	c.JSON(http.StatusOK, resp)
}

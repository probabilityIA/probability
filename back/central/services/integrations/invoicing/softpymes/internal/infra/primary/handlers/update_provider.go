package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/handlers/response"
)

// UpdateProvider actualiza un proveedor Softpymes
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

	provider, err := h.useCase.UpdateProvider(ctx, uint(id), dto)
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("provider_id", uint(id)).Msg("Failed to update provider")
		c.JSON(http.StatusInternalServerError, response.Error{
			Error:   "provider_update_failed",
			Message: err.Error(),
		})
		return
	}

	resp := mappers.ProviderToResponse(provider, "softpymes")

	h.log.Info(ctx).Uint("provider_id", provider.ID).Msg("Provider updated successfully")

	c.JSON(http.StatusOK, resp)
}

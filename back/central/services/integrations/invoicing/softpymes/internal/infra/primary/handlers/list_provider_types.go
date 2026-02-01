package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/handlers/response"
)

// ListProviderTypes lista los tipos de proveedores disponibles
func (h *handler) ListProviderTypes(c *gin.Context) {
	ctx := c.Request.Context()

	h.log.Debug(ctx).Msg("Listing provider types")

	providerTypes, err := h.useCase.ListProviderTypes(ctx)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("Failed to list provider types")
		c.JSON(http.StatusInternalServerError, response.Error{
			Error:   "list_provider_types_failed",
			Message: err.Error(),
		})
		return
	}

	resp := mappers.ProviderTypesToResponse(providerTypes)
	c.JSON(http.StatusOK, resp)
}

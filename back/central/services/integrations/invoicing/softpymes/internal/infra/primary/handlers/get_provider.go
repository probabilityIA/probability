package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/handlers/response"
)

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

	// TODO: Obtener tipo de proveedor para incluir en response
	resp := mappers.ProviderToResponse(provider, "softpymes")
	c.JSON(http.StatusOK, resp)
}

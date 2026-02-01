package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/handlers/response"
)

// TestProvider prueba la conexión con Softpymes
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
		Message: "Connection test successful with Softpymes API",
	})
}

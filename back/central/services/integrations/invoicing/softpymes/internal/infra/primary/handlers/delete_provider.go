package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/handlers/response"
)

// DeleteProvider elimina un proveedor Softpymes (soft delete)
func (h *handler) DeleteProvider(c *gin.Context) {
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

	h.log.Info(ctx).Uint("provider_id", uint(id)).Msg("Deleting provider")

	err = h.useCase.DeleteProvider(ctx, uint(id))
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("provider_id", uint(id)).Msg("Failed to delete provider")
		c.JSON(http.StatusInternalServerError, response.Error{
			Error:   "provider_deletion_failed",
			Message: err.Error(),
		})
		return
	}

	h.log.Info(ctx).Uint("provider_id", uint(id)).Msg("Provider deleted successfully")

	c.JSON(http.StatusNoContent, nil)
}

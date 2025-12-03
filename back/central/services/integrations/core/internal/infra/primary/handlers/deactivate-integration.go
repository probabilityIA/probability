package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/response"
)

// DeactivateIntegrationHandler desactiva una integración
//
//	@Summary		Desactivar integración
//	@Description	Desactiva una integración del sistema
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int		true	"ID de la integración"
//	@Success		200	{object}	response.IntegrationMessageResponse
//	@Failure		400	{object}	response.IntegrationErrorResponse
//	@Failure		404	{object}	response.IntegrationErrorResponse
//	@Failure		500	{object}	response.IntegrationErrorResponse
//	@Router			/integrations/{id}/deactivate [put]
func (h *IntegrationHandler) DeactivateIntegrationHandler(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, response.IntegrationErrorResponse{
			Success: false,
			Message: "Solo los super usuarios pueden desactivar integraciones",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "ID inválido",
		})
		return
	}

	if err := h.usecase.DeactivateIntegration(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.IntegrationErrorResponse{
			Success: false,
			Message: "Error al desactivar integración",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.IntegrationMessageResponse{
		Success: true,
		Message: "Integración desactivada exitosamente",
	})
}

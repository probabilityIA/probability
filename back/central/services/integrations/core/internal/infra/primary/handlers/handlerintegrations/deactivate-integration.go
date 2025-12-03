package handlerintegrations

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
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
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/integrations/{id}/deactivate [put]
func (h *IntegrationHandler) DeactivateIntegrationHandler(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		h.logger.Error().Str("endpoint", "/integrations/:id/deactivate").Str("method", "PUT").Msg("Intento de desactivar integración sin permisos de super admin")
		c.JSON(http.StatusForbidden, response.IntegrationErrorResponse{
			Success: false,
			Message: "Solo los super usuarios pueden desactivar integraciones",
			Error:   "permisos insuficientes",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Str("endpoint", "/integrations/:id/deactivate").Str("method", "PUT").Msg("ID de integración inválido al intentar desactivar")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "ID inválido",
			Error:   "El ID debe ser un número válido",
		})
		return
	}

	if err := h.usecase.DeactivateIntegration(c.Request.Context(), uint(id)); err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := "Error al desactivar integración"

		if errors.Is(err, domain.ErrIntegrationNotFound) {
			statusCode = http.StatusNotFound
			errorMsg = "La integración especificada no existe"
		}

		h.logger.Error().
			Err(err).
			Uint64("integration_id", id).
			Int("status_code", statusCode).
			Msg("Error al desactivar integración en el usecase")
		c.JSON(statusCode, response.IntegrationErrorResponse{
			Success: false,
			Message: errorMsg,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.IntegrationMessageResponse{
		Success: true,
		Message: "Integración desactivada exitosamente",
	})
}

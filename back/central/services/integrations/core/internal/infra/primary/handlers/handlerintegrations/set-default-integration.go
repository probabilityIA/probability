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

// SetAsDefaultHandler marca una integración como default
//
//	@Summary		Marcar como default
//	@Description	Marca una integración como default para su tipo
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int		true	"ID de la integración"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
func (h *IntegrationHandler) SetAsDefaultHandler(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		h.logger.Error().Str("endpoint", "/integrations/:id/set-default").Str("method", "PUT").Msg("Intento de marcar integración como default sin permisos de super admin")
		c.JSON(http.StatusForbidden, response.IntegrationErrorResponse{
			Success: false,
			Message: "Solo los super usuarios pueden marcar integraciones como default",
			Error:   "permisos insuficientes",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Str("endpoint", "/integrations/:id/set-default").Str("method", "PUT").Msg("ID de integración inválido al intentar marcar como default")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "ID inválido",
			Error:   "El ID debe ser un número válido",
		})
		return
	}

	if err := h.usecase.SetAsDefault(c.Request.Context(), uint(id)); err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := "Error al marcar integración como default"

		if errors.Is(err, domain.ErrIntegrationNotFound) {
			statusCode = http.StatusNotFound
			errorMsg = "La integración especificada no existe"
		}

		h.logger.Error().
			Err(err).
			Uint64("integration_id", id).
			Int("status_code", statusCode).
			Msg("Error al marcar integración como default en el usecase")
		c.JSON(statusCode, response.IntegrationErrorResponse{
			Success: false,
			Message: errorMsg,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.IntegrationMessageResponse{
		Success: true,
		Message: "Integración marcada como default exitosamente",
	})
}

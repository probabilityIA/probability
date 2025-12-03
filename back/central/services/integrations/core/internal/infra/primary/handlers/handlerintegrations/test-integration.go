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

// TestIntegrationHandler prueba la conexión de una integración
//
//	@Summary		Probar integración
//	@Description	Prueba la conexión de una integración sin guardar cambios
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
//	@Router			/integrations/{id}/test [post]
func (h *IntegrationHandler) TestIntegrationHandler(c *gin.Context) {
	// Solo super admins pueden probar integraciones
	if !middleware.IsSuperAdmin(c) {
		h.logger.Error().Str("endpoint", "/integrations/:id/test").Str("method", "POST").Msg("Intento de probar integración sin permisos de super admin")
		c.JSON(http.StatusForbidden, response.IntegrationErrorResponse{
			Success: false,
			Message: "Solo los super usuarios pueden probar integraciones",
			Error:   "permisos insuficientes",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Str("endpoint", "/integrations/:id/test").Str("method", "POST").Msg("ID de integración inválido al intentar probar")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "ID inválido",
			Error:   "El ID debe ser un número válido",
		})
		return
	}

	if err := h.usecase.TestIntegration(c.Request.Context(), uint(id)); err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := "Error al probar conexión de la integración"

		if errors.Is(err, domain.ErrIntegrationNotFound) {
			statusCode = http.StatusNotFound
			errorMsg = "La integración especificada no existe"
		} else if errors.Is(err, domain.ErrIntegrationCredentialsDecrypt) {
			errorMsg = "Error al procesar credenciales de la integración"
		} else if errors.Is(err, domain.ErrIntegrationConfigDeserialize) {
			errorMsg = "Error al procesar configuración de la integración"
		} else if errors.Is(err, domain.ErrIntegrationTestFailed) {
			errorMsg = "La prueba de conexión falló. Verifique las credenciales y configuración"
		} else if errors.Is(err, domain.ErrIntegrationAccessTokenNotFound) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		h.logger.Error().
			Err(err).
			Uint64("integration_id", id).
			Int("status_code", statusCode).
			Msg("Error al probar integración en el usecase")
		c.JSON(statusCode, response.IntegrationErrorResponse{
			Success: false,
			Message: errorMsg,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.IntegrationMessageResponse{
		Success: true,
		Message: "Conexión probada exitosamente",
	})
}

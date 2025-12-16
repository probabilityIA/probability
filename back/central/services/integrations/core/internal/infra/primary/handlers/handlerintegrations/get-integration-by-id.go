package handlerintegrations

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
)

// GetIntegrationByIDHandler obtiene una integración por su ID
//
//	@Summary		Obtener integración por ID
//	@Description	Obtiene una integración específica por su ID
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
//	@Router			/integrations/{id} [get]
func (h *IntegrationHandler) GetIntegrationByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "ID inválido",
			Error:   "El ID debe ser un número válido",
		})
		return
	}

	// Si es super admin, obtener con credenciales desencriptadas
	if middleware.IsSuperAdmin(c) {
		integrationWithCreds, err := h.usecase.GetIntegrationByIDWithCredentials(c.Request.Context(), uint(id))
		if err != nil {
			statusCode := http.StatusNotFound
			errorMsg := "Integración no encontrada"

			if !errors.Is(err, domain.ErrIntegrationNotFound) {
				statusCode = http.StatusInternalServerError
				errorMsg = "Error interno del servidor al obtener la integración"
			}

			h.logger.Error().
				Err(err).
				Uint64("integration_id", id).
				Int("status_code", statusCode).
				Msg("Error al obtener integración por ID con credenciales en el usecase")
			c.JSON(statusCode, response.IntegrationErrorResponse{
				Success: false,
				Message: errorMsg,
				Error:   err.Error(),
			})
			return
		}

		// Convertir a respuesta con credenciales desencriptadas
		imageURLBase := h.getImageURLBase()
		integrationResp := mapper.ToIntegrationResponse(&integrationWithCreds.Integration, imageURLBase)
		integrationResp.Credentials = integrationWithCreds.DecryptedCredentials

		c.JSON(http.StatusOK, response.IntegrationSuccessResponse{
			Success: true,
			Message: "Integración obtenida exitosamente",
			Data:    integrationResp,
		})
		return
	}

	// Si no es super admin, obtener sin credenciales (comportamiento normal)
	integration, err := h.usecase.GetIntegrationByID(c.Request.Context(), uint(id))
	if err != nil {
		statusCode := http.StatusNotFound
		errorMsg := "Integración no encontrada"

		if !errors.Is(err, domain.ErrIntegrationNotFound) {
			statusCode = http.StatusInternalServerError
			errorMsg = "Error interno del servidor al obtener la integración"
		}

		h.logger.Error().
			Err(err).
			Uint64("integration_id", id).
			Int("status_code", statusCode).
			Msg("Error al obtener integración por ID en el usecase")
		c.JSON(statusCode, response.IntegrationErrorResponse{
			Success: false,
			Message: errorMsg,
			Error:   err.Error(),
		})
		return
	}

	imageURLBase := h.getImageURLBase()
	integrationResp := mapper.ToIntegrationResponse(integration, imageURLBase)
	c.JSON(http.StatusOK, response.IntegrationSuccessResponse{
		Success: true,
		Message: "Integración obtenida exitosamente",
		Data:    integrationResp,
	})
}

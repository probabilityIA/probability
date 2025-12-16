package handlerintegrations

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/request"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
)

// UpdateIntegrationHandler actualiza una integración
//
//	@Summary		Actualizar integración
//	@Description	Actualiza una integración existente
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int									true	"ID de la integración"
//	@Param			request	body		request.UpdateIntegrationRequest	true	"Datos a actualizar"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/integrations/{id} [put]
func (h *IntegrationHandler) UpdateIntegrationHandler(c *gin.Context) {
	// Solo super admins pueden actualizar integraciones
	if !middleware.IsSuperAdmin(c) {
		h.logger.Error().Str("endpoint", "/integrations/:id").Str("method", "PUT").Msg("Intento de actualizar integración sin permisos de super admin")
		c.JSON(http.StatusForbidden, response.IntegrationErrorResponse{
			Success: false,
			Message: "Solo los super usuarios pueden actualizar integraciones",
			Error:   "permisos insuficientes",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Str("endpoint", "/integrations/:id").Str("method", "PUT").Msg("ID de integración inválido al intentar actualizar")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "ID inválido",
			Error:   "El ID debe ser un número válido",
		})
		return
	}

	var req request.UpdateIntegrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Uint64("id", id).Str("endpoint", "/integrations/:id").Str("method", "PUT").Msg("Error al parsear datos JSON para actualizar integración")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Error:   err.Error(),
		})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		h.logger.Error().Uint64("integration_id", id).Str("endpoint", "/integrations/:id").Str("method", "PUT").Msg("Intento de actualizar integración sin usuario autenticado")
		c.JSON(http.StatusUnauthorized, response.IntegrationErrorResponse{
			Success: false,
			Message: "Usuario no autenticado",
			Error:   "token de autenticación inválido o ausente",
		})
		return
	}

	dto := mapper.ToUpdateIntegrationDTO(req, userID)
	integration, err := h.usecase.UpdateIntegration(c.Request.Context(), uint(id), dto)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := "Error al actualizar integración"

		if errors.Is(err, domain.ErrIntegrationNotFound) {
			statusCode = http.StatusNotFound
			errorMsg = "La integración especificada no existe"
		} else if errors.Is(err, domain.ErrIntegrationCodeExists) {
			statusCode = http.StatusConflict
			errorMsg = "Ya existe otra integración con el código proporcionado"
		}

		h.logger.Error().
			Err(err).
			Uint64("integration_id", id).
			Uint("user_id", userID).
			Int("status_code", statusCode).
			Msg("Error al actualizar integración en el usecase")
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
		Message: "Integración actualizada exitosamente",
		Data:    integrationResp,
	})
}

package handlerintegrations

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/request"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
)

// CreateIntegrationHandler crea una nueva integración
//
//	@Summary		Crear integración
//	@Description	Crea una nueva integración en el sistema
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		request.CreateIntegrationRequest	true	"Datos de la integración"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		409		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/integrations [post]
func (h *IntegrationHandler) CreateIntegrationHandler(c *gin.Context) {
	// Solo super admins pueden crear integraciones
	if !middleware.IsSuperAdmin(c) {
		h.logger.Error().Str("endpoint", "/integrations").Str("method", "POST").Msg("Intento de crear integración sin permisos de super admin")
		c.JSON(http.StatusForbidden, response.IntegrationErrorResponse{
			Success: false,
			Message: "Solo los super usuarios pueden crear integraciones",
			Error:   "permisos insuficientes",
		})
		return
	}

	var req request.CreateIntegrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Str("endpoint", "/integrations").Str("method", "POST").Msg("Error al validar datos de entrada para crear integración")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Error:   err.Error(),
		})
		return
	}

	// Obtener ID del usuario autenticado
	userID, exists := middleware.GetUserID(c)
	if !exists {
		h.logger.Error().Str("endpoint", "/integrations").Str("method", "POST").Msg("Intento de crear integración sin usuario autenticado")
		c.JSON(http.StatusUnauthorized, response.IntegrationErrorResponse{
			Success: false,
			Message: "Usuario no autenticado",
			Error:   "token de autenticación inválido o ausente",
		})
		return
	}

	dto := mapper.ToCreateIntegrationDTO(req, userID)
	integration, err := h.usecase.CreateIntegration(c.Request.Context(), dto)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := "Error al crear integración"

		if errors.Is(err, domain.ErrIntegrationCodeExists) {
			statusCode = http.StatusConflict
			errorMsg = "Ya existe una integración con el código proporcionado"
		} else if errors.Is(err, domain.ErrIntegrationNameRequired) ||
			errors.Is(err, domain.ErrIntegrationCodeRequired) ||
			errors.Is(err, domain.ErrIntegrationTypeRequired) ||
			errors.Is(err, domain.ErrIntegrationCategoryInvalid) {
			statusCode = http.StatusBadRequest
			errorMsg = err.Error()
		}

		h.logger.Error().
			Err(err).
			Uint("user_id", userID).
			Str("integration_code", req.Code).
			Uint("integration_type_id", req.IntegrationTypeID).
			Int("status_code", statusCode).
			Msg("Error al crear integración en el usecase")
		c.JSON(statusCode, response.IntegrationErrorResponse{
			Success: false,
			Message: errorMsg,
			Error:   err.Error(),
		})
		return
	}

	imageURLBase := h.getImageURLBase()
	integrationResp := mapper.ToIntegrationResponse(integration, imageURLBase)
	c.JSON(http.StatusCreated, response.IntegrationSuccessResponse{
		Success: true,
		Message: "Integración creada exitosamente",
		Data:    integrationResp,
	})
}

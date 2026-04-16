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
	// Verificar permisos: Super Admin O rol de Administrador (ID 4)
	isSuper := middleware.IsSuperAdmin(c)
	roleID, _ := middleware.GetRoleID(c)

	// TODO: Implementar validación granular de permisos (ej: "integrations:create")
	// Por ahora permitimos explícitamente al rol 4 (Administrador)
	if !isSuper && roleID != 4 {
		h.logger.Error().
			Str("endpoint", "/integrations").
			Str("method", "POST").
			Uint("role_id", roleID).
			Msg("Intento de crear integración sin permisos suficientes")
		c.JSON(http.StatusForbidden, response.IntegrationErrorResponse{
			Success: false,
			Message: "No tienes permisos para crear integraciones",
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

	businessID := c.GetUint("business_id")

	dto := mapper.ToCreateIntegrationDTO(req, userID, businessID)

	// Validar que se haya asignado un BusinessID (crítico para la sincronización)
	if dto.BusinessID == nil || *dto.BusinessID == 0 {
		h.logger.Error().
			Uint("user_id", userID).
			Msg("Intento de crear integración sin BusinessID asignado (Super Admin debe especificar business_id)")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "BusinessID es requerido",
			Error:   "Como super admin, debes especificar el business_id en el cuerpo de la solicitud",
		})
		return
	}

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

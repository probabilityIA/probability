package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/response"
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
//	@Success		201		{object}	response.IntegrationSuccessResponse
//	@Failure		400		{object}	response.IntegrationErrorResponse
//	@Failure		401		{object}	response.IntegrationErrorResponse
//	@Failure		409		{object}	response.IntegrationErrorResponse
//	@Failure		500		{object}	response.IntegrationErrorResponse
//	@Router			/integrations [post]
func (h *IntegrationHandler) CreateIntegrationHandler(c *gin.Context) {
	// Solo super admins pueden crear integraciones
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, response.IntegrationErrorResponse{
			Success: false,
			Message: "Solo los super usuarios pueden crear integraciones",
			Error:   "permisos insuficientes",
		})
		return
	}

	var req request.CreateIntegrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Error al validar datos de entrada")
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
		c.JSON(http.StatusUnauthorized, response.IntegrationErrorResponse{
			Success: false,
			Message: "Usuario no autenticado",
		})
		return
	}

	dto := mapper.ToCreateIntegrationDTO(req, userID)
	integration, err := h.usecase.CreateIntegration(c.Request.Context(), dto)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al crear integración")
		statusCode := http.StatusInternalServerError
		if err.Error() == "ya existe una integración con el código" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, response.IntegrationErrorResponse{
			Success: false,
			Message: "Error al crear integración",
			Error:   err.Error(),
		})
		return
	}

	integrationResp := mapper.ToIntegrationResponse(integration)
	c.JSON(http.StatusCreated, response.IntegrationSuccessResponse{
		Success: true,
		Message: "Integración creada exitosamente",
		Data:    integrationResp,
	})
}

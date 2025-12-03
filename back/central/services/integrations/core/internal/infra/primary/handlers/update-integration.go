package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/response"
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
//	@Success		200		{object}	response.IntegrationSuccessResponse
//	@Failure		400		{object}	response.IntegrationErrorResponse
//	@Failure		404		{object}	response.IntegrationErrorResponse
//	@Failure		500		{object}	response.IntegrationErrorResponse
//	@Router			/integrations/{id} [put]
func (h *IntegrationHandler) UpdateIntegrationHandler(c *gin.Context) {
	// Solo super admins pueden actualizar integraciones
	if !middleware.IsSuperAdmin(c) {
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
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "ID inválido",
		})
		return
	}

	var req request.UpdateIntegrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Error:   err.Error(),
		})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, response.IntegrationErrorResponse{
			Success: false,
			Message: "Usuario no autenticado",
		})
		return
	}

	dto := mapper.ToUpdateIntegrationDTO(req, userID)
	integration, err := h.usecase.UpdateIntegration(c.Request.Context(), uint(id), dto)
	if err != nil {
		h.logger.Error().Err(err).Uint64("id", id).Msg("Error al actualizar integración")
		statusCode := http.StatusInternalServerError
		if err.Error() == "integración no encontrada" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, response.IntegrationErrorResponse{
			Success: false,
			Message: "Error al actualizar integración",
			Error:   err.Error(),
		})
		return
	}

	integrationResp := mapper.ToIntegrationResponse(integration)
	c.JSON(http.StatusOK, response.IntegrationSuccessResponse{
		Success: true,
		Message: "Integración actualizada exitosamente",
		Data:    integrationResp,
	})
}

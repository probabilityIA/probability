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

	if !middleware.IsSuperAdmin(c) {
		businessID, hasBusinessID := middleware.GetBusinessID(c)
		if !hasBusinessID || businessID == 0 {
			c.JSON(http.StatusForbidden, response.IntegrationErrorResponse{
				Success: false,
				Message: "No tienes permisos para ver esta integración",
				Error:   "permisos insuficientes",
			})
			return
		}
		existing, err := h.usecase.GetIntegrationByID(c.Request.Context(), uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, response.IntegrationErrorResponse{
				Success: false,
				Message: "Integración no encontrada",
				Error:   err.Error(),
			})
			return
		}
		if existing.BusinessID == nil || *existing.BusinessID != businessID {
			h.logger.Error().
				Uint64("integration_id", id).
				Uint("business_id", businessID).
				Str("endpoint", "/integrations/:id").
				Str("method", "GET").
				Msg("Intento de ver integración de otro negocio")
			c.JSON(http.StatusForbidden, response.IntegrationErrorResponse{
				Success: false,
				Message: "No tienes permisos para ver esta integración",
				Error:   "la integración no pertenece a tu negocio",
			})
			return
		}
	}

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
			Msg("Error al obtener integración por ID con credenciales")
		c.JSON(statusCode, response.IntegrationErrorResponse{
			Success: false,
			Message: errorMsg,
			Error:   err.Error(),
		})
		return
	}

	imageURLBase := h.getImageURLBase()
	integrationResp := mapper.ToIntegrationResponse(&integrationWithCreds.Integration, imageURLBase)
	integrationResp.Credentials = integrationWithCreds.DecryptedCredentials

	c.JSON(http.StatusOK, response.IntegrationSuccessResponse{
		Success: true,
		Message: "Integración obtenida exitosamente",
		Data:    integrationResp,
	})
}

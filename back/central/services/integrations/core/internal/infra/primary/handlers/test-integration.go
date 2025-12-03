package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/response"
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
//	@Success		200	{object}	response.IntegrationMessageResponse
//	@Failure		400	{object}	response.IntegrationErrorResponse
//	@Failure		404	{object}	response.IntegrationErrorResponse
//	@Failure		500	{object}	response.IntegrationErrorResponse
//	@Router			/integrations/{id}/test [post]
func (h *IntegrationHandler) TestIntegrationHandler(c *gin.Context) {
	// Solo super admins pueden probar integraciones
	if !middleware.IsSuperAdmin(c) {
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
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "ID inválido",
		})
		return
	}

	if err := h.usecase.TestIntegration(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error().Err(err).Uint64("id", id).Msg("Error al probar integración")
		c.JSON(http.StatusInternalServerError, response.IntegrationErrorResponse{
			Success: false,
			Message: "Error al probar integración",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.IntegrationMessageResponse{
		Success: true,
		Message: "Conexión probada exitosamente",
	})
}

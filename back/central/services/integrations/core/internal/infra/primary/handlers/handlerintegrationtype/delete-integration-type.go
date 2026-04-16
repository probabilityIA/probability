package handlerintegrationtype

// IntegrationTypeHandler está definido en integration-type-handler-constructor.go

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/response"
)

// DeleteIntegrationTypeHandler elimina un tipo de integración
//
//	@Summary		Eliminar tipo de integración
//	@Description	Elimina un tipo de integración del sistema
//	@Tags			IntegrationTypes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"ID del tipo de integración"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/integration-types/{id} [delete]
func (h *IntegrationTypeHandler) DeleteIntegrationTypeHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Str("endpoint", "/integration-types/:id").Str("method", "DELETE").Msg("ID de tipo de integración inválido al intentar eliminar")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "ID inválido",
			Error:   "El ID debe ser un número válido",
		})
		return
	}

	if err := h.usecase.DeleteIntegrationType(c.Request.Context(), uint(id)); err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := "Error al eliminar tipo de integración"

		if errors.Is(err, domain.ErrIntegrationTypeNotFound) {
			statusCode = http.StatusNotFound
			errorMsg = "El tipo de integración especificado no existe"
		} else if errors.Is(err, domain.ErrIntegrationTypeHasIntegrations) {
			statusCode = http.StatusConflict
			errorMsg = "No se puede eliminar un tipo de integración que tiene integraciones asociadas"
		}

		h.logger.Error().
			Err(err).
			Uint64("integration_type_id", id).
			Int("status_code", statusCode).
			Msg("Error al eliminar tipo de integración en el usecase")
		c.JSON(statusCode, response.IntegrationErrorResponse{
			Success: false,
			Message: errorMsg,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.IntegrationMessageResponse{
		Success: true,
		Message: "Tipo de integración eliminado exitosamente",
	})
}

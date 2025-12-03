package handlerintegrationtype

// IntegrationTypeHandler está definido en integration-type-handler-constructor.go

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/request"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/response"
)

// UpdateIntegrationTypeHandler actualiza un tipo de integración
//
//	@Summary		Actualizar tipo de integración
//	@Description	Actualiza un tipo de integración existente
//	@Tags			IntegrationTypes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id					path		int									true	"ID del tipo de integración"
//	@Param			integrationType		body		request.UpdateIntegrationTypeRequest	true	"Datos actualizados del tipo de integración"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/integration-types/{id} [put]
func (h *IntegrationTypeHandler) UpdateIntegrationTypeHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Str("endpoint", "/integration-types/:id").Str("method", "PUT").Msg("ID de tipo de integración inválido al intentar actualizar")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "ID inválido",
			Error:   "El ID debe ser un número válido",
		})
		return
	}

	var req request.UpdateIntegrationTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Uint64("id", id).Str("endpoint", "/integration-types/:id").Str("method", "PUT").Msg("Error al parsear datos JSON para actualizar tipo de integración")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Error:   err.Error(),
		})
		return
	}

	dto := mapper.ToUpdateIntegrationTypeDTO(req)

	integrationType, err := h.usecase.UpdateIntegrationType(c.Request.Context(), uint(id), dto)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := "Error al actualizar tipo de integración"

		if errors.Is(err, domain.ErrIntegrationTypeNotFound) {
			statusCode = http.StatusNotFound
			errorMsg = "El tipo de integración especificado no existe"
		} else if errors.Is(err, domain.ErrIntegrationTypeCodeExists) {
			statusCode = http.StatusConflict
			errorMsg = "Ya existe otro tipo de integración con el código proporcionado"
		} else if errors.Is(err, domain.ErrIntegrationTypeNameExists) {
			statusCode = http.StatusConflict
			errorMsg = "Ya existe otro tipo de integración con el nombre proporcionado"
		}

		h.logger.Error().
			Err(err).
			Uint64("integration_type_id", id).
			Int("status_code", statusCode).
			Msg("Error al actualizar tipo de integración en el usecase")
		c.JSON(statusCode, response.IntegrationErrorResponse{
			Success: false,
			Message: errorMsg,
			Error:   err.Error(),
		})
		return
	}

	integrationTypeResp := mapper.ToIntegrationTypeResponse(integrationType)
	c.JSON(http.StatusOK, response.IntegrationTypeDetailResponse{
		Success: true,
		Message: "Tipo de integración actualizado exitosamente",
		Data:    integrationTypeResp,
	})
}

package handlerintegrationtype

// IntegrationTypeHandler está definido en integration-type-handler-constructor.go

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/request"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/response"
)

// CreateIntegrationTypeHandler crea un nuevo tipo de integración
//
//	@Summary		Crear tipo de integración
//	@Description	Crea un nuevo tipo de integración en el sistema
//	@Tags			IntegrationTypes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			integrationType	body		request.CreateIntegrationTypeRequest	true	"Datos del tipo de integración"
//	@Success		201				{object}	map[string]interface{}
//	@Failure		400				{object}	map[string]interface{}
//	@Failure		401				{object}	map[string]interface{}
//	@Failure		500				{object}	map[string]interface{}
//	@Router			/integration-types [post]
func (h *IntegrationTypeHandler) CreateIntegrationTypeHandler(c *gin.Context) {
	var req request.CreateIntegrationTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Str("endpoint", "/integration-types").Str("method", "POST").Msg("Error al parsear request JSON para crear tipo de integración")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Error:   err.Error(),
		})
		return
	}

	dto := mapper.ToCreateIntegrationTypeDTO(req)

	integrationType, err := h.usecase.CreateIntegrationType(c.Request.Context(), dto)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := "Error al crear tipo de integración"

		if errors.Is(err, domain.ErrIntegrationTypeCodeExists) {
			statusCode = http.StatusConflict
			errorMsg = "Ya existe un tipo de integración con el código proporcionado"
		} else if errors.Is(err, domain.ErrIntegrationTypeNameExists) {
			statusCode = http.StatusConflict
			errorMsg = "Ya existe un tipo de integración con el nombre proporcionado"
		}

		h.logger.Error().
			Err(err).
			Str("code", req.Code).
			Str("name", req.Name).
			Int("status_code", statusCode).
			Msg("Error al crear tipo de integración en el usecase")
		c.JSON(statusCode, response.IntegrationErrorResponse{
			Success: false,
			Message: errorMsg,
			Error:   err.Error(),
		})
		return
	}

	integrationTypeResp := mapper.ToIntegrationTypeResponse(integrationType)
	c.JSON(http.StatusCreated, response.IntegrationTypeDetailResponse{
		Success: true,
		Message: "Tipo de integración creado exitosamente",
		Data:    integrationTypeResp,
	})
}

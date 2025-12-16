package handlerintegrationtype

// IntegrationTypeHandler está definido en integration-type-handler-constructor.go

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/response"
)

// GetIntegrationTypeByIDHandler obtiene un tipo de integración por ID
//
//	@Summary		Obtener tipo de integración por ID
//	@Description	Obtiene un tipo de integración específico por su ID
//	@Tags			IntegrationTypes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"ID del tipo de integración"	example(1)
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/integration-types/{id} [get]
func (h *IntegrationTypeHandler) GetIntegrationTypeByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Str("endpoint", "/integration-types/:id").Str("method", "GET").Msg("ID de tipo de integración inválido al intentar obtener")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "ID inválido",
			Error:   "El ID debe ser un número válido",
		})
		return
	}

	integrationType, err := h.usecase.GetIntegrationTypeByID(c.Request.Context(), uint(id))
	if err != nil {
		statusCode := http.StatusNotFound
		errorMsg := "Tipo de integración no encontrado"

		if !errors.Is(err, domain.ErrIntegrationTypeNotFound) {
			statusCode = http.StatusInternalServerError
			errorMsg = "Error interno del servidor al obtener el tipo de integración"
		}

		h.logger.Error().
			Err(err).
			Uint64("integration_type_id", id).
			Int("status_code", statusCode).
			Msg("Error al obtener tipo de integración por ID en el usecase")
		c.JSON(statusCode, response.IntegrationErrorResponse{
			Success: false,
			Message: errorMsg,
			Error:   err.Error(),
		})
		return
	}

	// Obtener URL base de S3 para construir URLs completas
	imageURLBase := h.getImageURLBase()
	integrationTypeResp := mapper.ToIntegrationTypeResponse(integrationType, imageURLBase)
	c.JSON(http.StatusOK, response.IntegrationTypeDetailResponse{
		Success: true,
		Message: "Tipo de integración obtenido exitosamente",
		Data:    integrationTypeResp,
	})
}

// GetIntegrationTypeByCodeHandler obtiene un tipo de integración por código
//
//	@Summary		Obtener tipo de integración por código
//	@Description	Obtiene un tipo de integración específico por su código
//	@Tags			IntegrationTypes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			code	path		string	true	"Código del tipo de integración"	example(whatsapp)
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/integration-types/code/{code} [get]
func (h *IntegrationTypeHandler) GetIntegrationTypeByCodeHandler(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		h.logger.Error().Str("endpoint", "/integration-types/code/:code").Str("method", "GET").Msg("Código de tipo de integración vacío en la petición")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "Código requerido",
			Error:   "El parámetro 'code' no puede estar vacío",
		})
		return
	}

	integrationType, err := h.usecase.GetIntegrationTypeByCode(c.Request.Context(), code)
	if err != nil {
		statusCode := http.StatusNotFound
		errorMsg := "Tipo de integración no encontrado"

		if !errors.Is(err, domain.ErrIntegrationTypeNotFound) {
			statusCode = http.StatusInternalServerError
			errorMsg = "Error interno del servidor al obtener el tipo de integración"
		}

		h.logger.Error().
			Err(err).
			Str("code", code).
			Int("status_code", statusCode).
			Msg("Error al obtener tipo de integración por código en el usecase")
		c.JSON(statusCode, response.IntegrationErrorResponse{
			Success: false,
			Message: errorMsg,
			Error:   err.Error(),
		})
		return
	}

	// Obtener URL base de S3 para construir URLs completas
	imageURLBase := h.getImageURLBase()
	integrationTypeResp := mapper.ToIntegrationTypeResponse(integrationType, imageURLBase)
	c.JSON(http.StatusOK, response.IntegrationTypeDetailResponse{
		Success: true,
		Message: "Tipo de integración obtenido exitosamente",
		Data:    integrationTypeResp,
	})
}

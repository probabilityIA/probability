package handlerintegrationtype

// IntegrationTypeHandler está definido en integration-type-handler-constructor.go

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/response"
)

// ListIntegrationTypesHandler obtiene todos los tipos de integración
//
//	@Summary		Listar tipos de integración
//	@Description	Obtiene todos los tipos de integración disponibles
//	@Tags			IntegrationTypes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/integration-types [get]
func (h *IntegrationTypeHandler) ListIntegrationTypesHandler(c *gin.Context) {
	integrationTypes, err := h.usecase.ListIntegrationTypes(c.Request.Context())
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("endpoint", "/integration-types").
			Str("method", "GET").
			Int("status_code", http.StatusInternalServerError).
			Msg("Error al listar tipos de integración en el usecase")
		c.JSON(http.StatusInternalServerError, response.IntegrationErrorResponse{
			Success: false,
			Message: "Error interno del servidor al obtener la lista de tipos de integración",
			Error:   err.Error(),
		})
		return
	}

	// Obtener URL base de S3 para construir URLs completas
	imageURLBase := h.getImageURLBase()
	responses := make([]response.IntegrationTypeResponse, len(integrationTypes))
	for i, it := range integrationTypes {
		responses[i] = mapper.ToIntegrationTypeResponse(it, imageURLBase)
	}

	c.JSON(http.StatusOK, response.IntegrationTypeListResponse{
		Success: true,
		Message: "Tipos de integración obtenidos exitosamente",
		Data:    responses,
	})
}

// ListActiveIntegrationTypesHandler obtiene solo los tipos de integración activos
//
//	@Summary		Listar tipos de integración activos
//	@Description	Obtiene solo los tipos de integración que están activos
//	@Tags			IntegrationTypes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/integration-types/active [get]
func (h *IntegrationTypeHandler) ListActiveIntegrationTypesHandler(c *gin.Context) {
	integrationTypes, err := h.usecase.ListActiveIntegrationTypes(c.Request.Context())
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("endpoint", "/integration-types/active").
			Str("method", "GET").
			Int("status_code", http.StatusInternalServerError).
			Msg("Error al listar tipos de integración activos en el usecase")
		c.JSON(http.StatusInternalServerError, response.IntegrationErrorResponse{
			Success: false,
			Message: "Error interno del servidor al obtener la lista de tipos de integración activos",
			Error:   err.Error(),
		})
		return
	}

	// Obtener URL base de S3 para construir URLs completas
	imageURLBase := h.getImageURLBase()
	responses := make([]response.IntegrationTypeResponse, len(integrationTypes))
	for i, it := range integrationTypes {
		responses[i] = mapper.ToIntegrationTypeResponse(it, imageURLBase)
	}

	c.JSON(http.StatusOK, response.IntegrationTypeListResponse{
		Success: true,
		Message: "Tipos de integración activos obtenidos exitosamente",
		Data:    responses,
	})
}

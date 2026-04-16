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
//	@Description	Crea un nuevo tipo de integración en el sistema. Soporta JSON y multipart/form-data
//	@Tags			IntegrationTypes
//	@Accept			json,mpfd
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name			formData	string					false	"Nombre del tipo de integración"
//	@Param			code			formData	string					false	"Código del tipo de integración"
//	@Param			description		formData	string					false	"Descripción"
//	@Param			icon			formData	string					false	"Icono"
//	@Param			category		formData	string					false	"Categoría (internal/external)"
//	@Param			is_active		formData	boolean					false	"¿Activo?"
//	@Param			credentials_schema	formData	string					false	"JSON schema para credenciales"
//	@Param			image_file		formData	file					false	"Imagen del logo (sube a S3)"
//	@Param			integrationType	body		request.CreateIntegrationTypeRequest	true	"Datos del tipo de integración (JSON)"
//	@Success		201				{object}	map[string]interface{}
//	@Failure		400				{object}	map[string]interface{}
//	@Failure		401				{object}	map[string]interface{}
//	@Failure		500				{object}	map[string]interface{}
//	@Router			/integration-types [post]
func (h *IntegrationTypeHandler) CreateIntegrationTypeHandler(c *gin.Context) {
	var req request.CreateIntegrationTypeRequest
	// Seleccionar binder según Content-Type
	if c.ContentType() == "application/json" {
		if err := c.ShouldBindJSON(&req); err != nil {
			h.logger.Error().Err(err).Str("endpoint", "/integration-types").Str("method", "POST").Msg("Error al parsear request JSON para crear tipo de integración")
			c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
				Success: false,
				Message: "Datos de entrada inválidos",
				Error:   err.Error(),
			})
			return
		}
	} else if err := c.ShouldBind(&req); err != nil {
		h.logger.Error().Err(err).Str("endpoint", "/integration-types").Str("method", "POST").Msg("Error al parsear request para crear tipo de integración")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Error:   err.Error(),
		})
		return
	}

	// Log para verificar si el archivo está llegando
	if req.ImageFile != nil {
		h.logger.Info().Str("name", req.Name).Str("filename", req.ImageFile.Filename).Int64("size", req.ImageFile.Size).Msg("Archivo de imagen recibido en creación")
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

	// Obtener URL base de S3 para construir URLs completas
	imageURLBase := h.getImageURLBase()
	integrationTypeResp := mapper.ToIntegrationTypeResponse(integrationType, imageURLBase)
	c.JSON(http.StatusCreated, response.IntegrationTypeDetailResponse{
		Success: true,
		Message: "Tipo de integración creado exitosamente",
		Data:    integrationTypeResp,
	})
}

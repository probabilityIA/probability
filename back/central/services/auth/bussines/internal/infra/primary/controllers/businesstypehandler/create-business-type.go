package businesstypehandler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesstypehandler/mapper"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesstypehandler/request"
)

// CreateBusinessType godoc
//
//	@Summary		Crear un nuevo tipo de negocio
//	@Description	Crea un nuevo tipo de negocio en el sistema
//	@Tags			business-types
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			businessType	body		request.BusinessTypeRequest	true	"Datos del tipo de negocio a crear"
//	@Success		201				{object}	map[string]interface{}		"Tipo de negocio creado exitosamente"
//	@Failure		400				{object}	map[string]interface{}		"Solicitud inválida"
//	@Failure		401				{object}	map[string]interface{}		"Token de acceso requerido"
//	@Failure		500				{object}	map[string]interface{}		"Error interno del servidor"
//	@Router			/business-types [post]
func (h *BusinessTypeHandler) CreateBusinessTypeHandler(c *gin.Context) {
	var createRequest request.BusinessTypeRequest

	// Validar y parsear el request
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		h.logger.Error().Err(err).Msg("Error al validar los datos de entrada para crear tipo de negocio")
		c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("invalid_request", fmt.Sprintf("Los datos proporcionados son inválidos: %s", err.Error())))
		return
	}

	// Validar campos requeridos (solo Name es obligatorio)
	if createRequest.Name == "" {
		h.logger.Warn().Msg("Intento de crear tipo de negocio sin nombre")
		c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("missing_fields", "El nombre del tipo de negocio es obligatorio"))
		return
	}

	// Ejecutar caso de uso
	businessTypeRequest := mapper.RequestToDTO(createRequest)
	businessType, err := h.usecase.CreateBusinessType(c.Request.Context(), businessTypeRequest)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBusinessTypeNameAlreadyExists):
			h.logger.Warn().
				Str("name", createRequest.Name).
				Msg("Intento de crear tipo de negocio con nombre duplicado")
			c.JSON(http.StatusConflict, mapper.BuildErrorResponse("name_already_exists", "El nombre del tipo de negocio ya está en uso. Por favor, proporciona un nombre diferente."))
			return
		default:
			h.logger.Error().
				Err(err).
				Str("name", createRequest.Name).
				Str("code", createRequest.Code).
				Msg("Error inesperado al crear tipo de negocio")
			// Intentar extraer un mensaje de error más descriptivo
			errorMessage := "Error al crear el tipo de negocio. Por favor, verifica los datos e intenta nuevamente."
			if errMsg := err.Error(); errMsg != "" {
				// Si el error contiene información útil, incluirla
				if strings.Contains(errMsg, "error al") || strings.Contains(errMsg, "Error al") {
					errorMessage = errMsg
				}
			}
			c.JSON(http.StatusInternalServerError, mapper.BuildErrorResponse("internal_error", errorMessage))
			return
		}
	}

	// Construir respuesta exitosa usando el DTO retornado
	response := mapper.BuildCreateBusinessTypeResponseFromDTO(businessType, "Tipo de negocio creado exitosamente")
	c.JSON(http.StatusCreated, response)
}

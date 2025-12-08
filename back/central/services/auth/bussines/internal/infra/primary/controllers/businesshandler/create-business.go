package businesshandler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesshandler/mapper"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesshandler/request"
)

// CreateBusiness godoc
//
//	@Summary		Crear un nuevo negocio
//	@Description	Crea un nuevo negocio en el sistema
//	@Tags			businesses
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name				formData	string					true	"Nombre del negocio"
//	@Param			code				formData	string					false	"Código del negocio"
//	@Param			business_type_id	formData	int						false	"ID del tipo de negocio"
//	@Param			timezone			formData	string					false	"Zona horaria"
//	@Param			address				formData	string					false	"Dirección"
//	@Param			description			formData	string					false	"Descripción"
//	@Param			logo_file			formData	file					false	"Logo del negocio (sube a S3)"
//	@Param			primary_color		formData	string					false	"Color primario"
//	@Param			secondary_color		formData	string					false	"Color secundario"
//	@Param			tertiary_color		formData	string					false	"Color terciario"
//	@Param			quaternary_color	formData	string					false	"Color cuaternario"
//	@Param			navbar_image_file	formData	file					false	"Imagen de navbar (sube a S3)"
//	@Param			custom_domain		formData	string					false	"Dominio personalizado"
//	@Param			is_active			formData	boolean					false	"¿Activo?"
//	@Param			enable_delivery		formData	boolean					false	"Habilitar delivery"
//	@Param			enable_pickup		formData	boolean					false	"Habilitar pickup"
//	@Param			enable_reservations	formData	boolean					false	"Habilitar reservas"
//	@Success		201					{object}	map[string]interface{}	"Negocio creado exitosamente"
//	@Failure		400					{object}	map[string]interface{}	"Solicitud inválida"
//	@Failure		401					{object}	map[string]interface{}	"Token de acceso requerido"
//	@Failure		500					{object}	map[string]interface{}	"Error interno del servidor"
//	@Router			/businesses [post]
func (h *BusinessHandler) CreateBusinessHandler(c *gin.Context) {
	var createRequest request.BusinessRequest

	// Validar y parsear el request
	if err := c.ShouldBind(&createRequest); err != nil {
		h.logger.Error().Err(err).Msg("Error al validar los datos de entrada para crear negocio")
		c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("invalid_request", fmt.Sprintf("Los datos proporcionados son inválidos: %s", err.Error())))
		return
	}

	// Validar campos requeridos (solo Name es obligatorio)
	if createRequest.Name == "" {
		h.logger.Warn().Msg("Intento de crear negocio sin nombre")
		c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("missing_fields", "El nombre del negocio es obligatorio"))
		return
	}

	// Ejecutar caso de uso
	businessRequest := mapper.RequestToDTO(createRequest)
	business, err := h.usecase.CreateBusiness(c.Request.Context(), businessRequest)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBusinessTypeIDRequired):
			h.logger.Warn().
				Str("name", createRequest.Name).
				Uint("business_type_id", createRequest.BusinessTypeID).
				Msg("Intento de crear negocio sin tipo de negocio")
			c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("business_type_required", "El tipo de negocio es obligatorio. Por favor, proporciona un tipo de negocio válido."))
			return
		case errors.Is(err, domain.ErrBusinessTypeIDInvalid):
			h.logger.Warn().
				Str("name", createRequest.Name).
				Uint("business_type_id", createRequest.BusinessTypeID).
				Msg("Intento de crear negocio con tipo de negocio inválido")
			c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("business_type_invalid", "El tipo de negocio especificado no existe o no es válido. Por favor, verifica el ID del tipo de negocio."))
			return
		case errors.Is(err, domain.ErrBusinessCodeAlreadyExists):
			h.logger.Warn().
				Str("name", createRequest.Name).
				Str("code", createRequest.Code).
				Msg("Intento de crear negocio con código duplicado")
			c.JSON(http.StatusConflict, mapper.BuildErrorResponse("code_already_exists", "El código del negocio ya está en uso. Por favor, proporciona un código diferente o deja que se genere automáticamente."))
			return
		case errors.Is(err, domain.ErrBusinessDomainAlreadyExists):
			h.logger.Warn().
				Str("name", createRequest.Name).
				Str("domain", createRequest.CustomDomain).
				Msg("Intento de crear negocio con dominio personalizado duplicado")
			c.JSON(http.StatusConflict, mapper.BuildErrorResponse("domain_already_exists", "El dominio personalizado ya está en uso. Por favor, proporciona un dominio diferente."))
			return
		default:
			h.logger.Error().
				Err(err).
				Str("name", createRequest.Name).
				Str("code", createRequest.Code).
				Uint("business_type_id", createRequest.BusinessTypeID).
				Msg("Error inesperado al crear negocio")
			// Detectar errores de foreign key constraint
			errMsg := err.Error()
			if strings.Contains(errMsg, "foreign key constraint") || strings.Contains(errMsg, "SQLSTATE 23503") {
				errorMessage := "El tipo de negocio especificado no existe o no es válido. Por favor, verifica el ID del tipo de negocio."
				c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("foreign_key_constraint", errorMessage))
				return
			}
			// Intentar extraer un mensaje de error más descriptivo
			errorMessage := "Error al crear el negocio. Por favor, verifica los datos e intenta nuevamente."
			if errMsg != "" {
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
	response := mapper.BuildCreateBusinessResponseFromDTO(business, "Negocio creado exitosamente")
	c.JSON(http.StatusCreated, response)
}

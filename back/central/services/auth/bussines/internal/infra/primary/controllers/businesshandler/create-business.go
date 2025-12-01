package businesshandler

import (
	"errors"
	"fmt"
	"net/http"

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
//	@Param			code				formData	string					true	"Código del negocio"
//	@Param			business_type_id	formData	int						true	"ID del tipo de negocio"
//	@Param			timezone			formData	string					false	"Zona horaria"
//	@Param			address				formData	string					false	"Dirección"
//	@Param			description			formData	string					false	"Descripción"
//	@Param			logo_url			formData	file					false	"Logo del negocio"
//	@Param			primary_color		formData	string					false	"Color primario"
//	@Param			secondary_color		formData	string					false	"Color secundario"
//	@Param			custom_domain		formData	string					false	"Dominio personalizado"
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
		c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("invalid_request", fmt.Sprintf("Datos de entrada inválidos: %s", err.Error())))
		return
	}

	// Validar campos requeridos
	if createRequest.Name == "" || createRequest.Code == "" || createRequest.BusinessTypeID == 0 {
		c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("missing_fields", "Nombre, código y tipo de negocio son requeridos"))
		return
	}

	// Ejecutar caso de uso
	businessRequest := mapper.RequestToDTO(createRequest)
	business, err := h.usecase.CreateBusiness(c.Request.Context(), businessRequest)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBusinessCodeAlreadyExists):
			c.JSON(http.StatusConflict, mapper.BuildErrorResponse("code_already_exists", "El código del negocio ya está en uso"))
			return
		case errors.Is(err, domain.ErrBusinessDomainAlreadyExists):
			c.JSON(http.StatusConflict, mapper.BuildErrorResponse("domain_already_exists", "El dominio personalizado ya está en uso"))
			return
		default:
			h.logger.Error().Err(err).Msg("Error al crear negocio")
			c.JSON(http.StatusInternalServerError, mapper.BuildErrorResponse("internal_error", "Error interno del servidor"))
			return
		}
	}

	// Construir respuesta exitosa usando el DTO retornado
	response := mapper.BuildCreateBusinessResponseFromDTO(business, "Negocio creado exitosamente")
	c.JSON(http.StatusCreated, response)
}

package businesstypehandler

import (
	"net/http"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesstypehandler/mapper"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesstypehandler/request"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("invalid_request", "Datos de entrada inválidos"))
		return
	}

	// Validar campos requeridos
	if createRequest.Name == "" || createRequest.Code == "" || createRequest.Description == "" {
		c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("missing_fields", "Nombre, código y descripción son requeridos"))
		return
	}

	// Ejecutar caso de uso
	businessTypeRequest := mapper.RequestToDTO(createRequest)
	businessType, err := h.usecase.CreateBusinessType(c.Request.Context(), businessTypeRequest)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al crear tipo de negocio")
		c.JSON(http.StatusInternalServerError, mapper.BuildErrorResponse("internal_error", "Error interno del servidor"))
		return
	}

	// Construir respuesta exitosa usando el DTO retornado
	response := mapper.BuildCreateBusinessTypeResponseFromDTO(businessType, "Tipo de negocio creado exitosamente")
	c.JSON(http.StatusCreated, response)
}

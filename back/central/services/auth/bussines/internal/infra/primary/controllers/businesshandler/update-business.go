package businesshandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesshandler/mapper"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesshandler/request"
)

// UpdateBusiness godoc
//
//	@Summary		Actualizar negocio
//	@Description	Actualiza un negocio existente
//	@Tags			businesses
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id					path		int						true	"ID del negocio"
//	@Param			name				formData	string					false	"Nombre del negocio"
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
//	@Success		201					{object}	map[string]interface{}	"Negocio actualizado exitosamente"
//	@Failure		400					{object}	map[string]interface{}	"Solicitud inválida"
//	@Failure		401					{object}	map[string]interface{}	"Token de acceso requerido"
//	@Failure		500					{object}	map[string]interface{}	"Error interno del servidor"
//	@Router			/businesses/{id} [put]
func (h *BusinessHandler) UpdateBusinessHandler(c *gin.Context) {
	// Obtener ID del path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("invalid_id", "ID de negocio inválido"))
		return
	}

	var updateRequest request.UpdateBusinessRequest

	// Validar y parsear el request (ShouldBind maneja tanto JSON como FormData)
	if err := c.ShouldBind(&updateRequest); err != nil {
		h.logger.Error().Err(err).Str("content_type", c.GetHeader("Content-Type")).Msg("Error al parsear request de actualización")

		// Log adicional para debugging
		body, _ := c.GetRawData()
		h.logger.Error().Str("body", string(body)).Msg("Contenido del request")

		c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("invalid_request", "Datos de entrada inválidos"))
		return
	}

	// Ejecutar caso de uso
	businessRequest := mapper.UpdateRequestToUpdateDTO(updateRequest)
	business, err := h.usecase.UpdateBusiness(c.Request.Context(), uint(id), businessRequest)
	if err != nil {
		if err.Error() == "negocio no encontrado" {
			c.JSON(http.StatusNotFound, mapper.BuildErrorResponse("not_found", "Negocio no encontrado"))
			return
		}
		h.logger.Error().Err(err).Uint("id", uint(id)).Msg("Error al actualizar negocio")
		c.JSON(http.StatusInternalServerError, mapper.BuildErrorResponse("internal_error", "Error interno del servidor"))
		return
	}

	// Construir respuesta exitosa
	response := mapper.BuildUpdateBusinessResponseFromDTO(business, "Negocio actualizado exitosamente")
	c.JSON(http.StatusOK, response)
}

package businesstypehandler

import (
	"net/http"
	"strconv"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesstypehandler/mapper"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesstypehandler/request"

	"github.com/gin-gonic/gin"
)

// UpdateBusinessType godoc
//
//	@Summary		Actualizar tipo de negocio
//	@Description	Actualiza un tipo de negocio existente
//	@Tags			business-types
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id				path		int							true	"ID del tipo de negocio"
//	@Param			businessType	body		request.BusinessTypeRequest	true	"Datos del tipo de negocio a actualizar"
//	@Success		201				{object}	map[string]interface{}		"Tipo de negocio actualizado exitosamente"
//	@Failure		400				{object}	map[string]interface{}		"Solicitud inválida"
//	@Failure		401				{object}	map[string]interface{}		"Token de acceso requerido"
//	@Failure		500				{object}	map[string]interface{}		"Error interno del servidor"
//	@Router			/business-types/{id} [put]
func (h *BusinessTypeHandler) UpdateBusinessTypeHandler(c *gin.Context) {
	// Obtener ID del path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("invalid_id", "ID de tipo de negocio inválido"))
		return
	}

	var updateRequest request.BusinessTypeRequest

	// Validar y parsear el request
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("invalid_request", "Datos de entrada inválidos"))
		return
	}

	// Ejecutar caso de uso
	businessTypeRequest := mapper.RequestToDTO(updateRequest)
	businessType, err := h.usecase.UpdateBusinessType(c.Request.Context(), uint(id), businessTypeRequest)
	if err != nil {
		if err.Error() == "tipo de negocio no encontrado" {
			c.JSON(http.StatusNotFound, mapper.BuildErrorResponse("not_found", "Tipo de negocio no encontrado"))
			return
		}
		h.logger.Error().Err(err).Uint("id", uint(id)).Msg("Error al actualizar tipo de negocio")
		c.JSON(http.StatusInternalServerError, mapper.BuildErrorResponse("internal_error", "Error interno del servidor"))
		return
	}

	// Construir respuesta exitosa
	response := mapper.BuildUpdateBusinessTypeResponseFromDTO(businessType, "Tipo de negocio actualizado exitosamente")
	c.JSON(http.StatusOK, response)
}

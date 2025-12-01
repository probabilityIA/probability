package businesstypehandler

import (
	"net/http"
	"strconv"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesstypehandler/mapper"

	"github.com/gin-gonic/gin"
)

// GetBusinessTypeByID godoc
//
//	@Summary		Obtener tipo de negocio por ID
//	@Description	Obtiene un tipo de negocio específico por su ID
//	@Tags			business-types
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int						true	"ID del tipo de negocio"
//	@Success		201	{object}	map[string]interface{}	"Tipo de negocio obtenido exitosamente"
//	@Failure		400	{object}	map[string]interface{}	"Solicitud inválida"
//	@Failure		401	{object}	map[string]interface{}	"Token de acceso requerido"
//	@Failure		500	{object}	map[string]interface{}	"Error interno del servidor"
//	@Router			/business-types/{id} [get]
func (h *BusinessTypeHandler) GetBusinessTypeByIDHandler(c *gin.Context) {
	// Obtener ID del path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("invalid_id", "ID de tipo de negocio inválido"))
		return
	}

	// Ejecutar caso de uso
	businessType, err := h.usecase.GetBusinessTypeByID(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "tipo de negocio no encontrado" {
			c.JSON(http.StatusNotFound, mapper.BuildErrorResponse("not_found", "Tipo de negocio no encontrado"))
			return
		}
		h.logger.Error().Err(err).Uint("id", uint(id)).Msg("Error al obtener tipo de negocio")
		c.JSON(http.StatusInternalServerError, mapper.BuildErrorResponse("internal_error", "Error interno del servidor"))
		return
	}

	// Construir respuesta exitosa
	response := mapper.BuildGetBusinessTypeResponseFromDTO(businessType, "Tipo de negocio obtenido exitosamente")
	c.JSON(http.StatusOK, response)
}

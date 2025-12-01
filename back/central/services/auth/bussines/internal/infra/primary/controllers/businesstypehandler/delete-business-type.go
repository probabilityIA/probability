package businesstypehandler

import (
	"net/http"
	"strconv"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesstypehandler/mapper"

	"github.com/gin-gonic/gin"
)

// DeleteBusinessType godoc
//
//	@Summary		Eliminar tipo de negocio
//	@Description	Elimina un tipo de negocio del sistema
//	@Tags			business-types
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int						true	"ID del tipo de negocio"
//	@Success		201	{object}	map[string]interface{}	"Tipo de negocio eliminado exitosamente"
//	@Failure		400	{object}	map[string]interface{}	"Solicitud inválida"
//	@Failure		401	{object}	map[string]interface{}	"Token de acceso requerido"
//	@Failure		500	{object}	map[string]interface{}	"Error interno del servidor"
//	@Router			/business-types/{id} [delete]
func (h *BusinessTypeHandler) DeleteBusinessTypeHandler(c *gin.Context) {
	// Obtener ID del path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("invalid_id", "ID de tipo de negocio inválido"))
		return
	}

	// Ejecutar caso de uso
	err = h.usecase.DeleteBusinessType(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "tipo de negocio no encontrado" {
			c.JSON(http.StatusNotFound, mapper.BuildErrorResponse("not_found", "Tipo de negocio no encontrado"))
			return
		}
		h.logger.Error().Err(err).Uint("id", uint(id)).Msg("Error al eliminar tipo de negocio")
		c.JSON(http.StatusInternalServerError, mapper.BuildErrorResponse("internal_error", "Error interno del servidor"))
		return
	}

	// Construir respuesta exitosa
	response := mapper.BuildDeleteBusinessTypeResponse("Tipo de negocio eliminado exitosamente")
	c.JSON(http.StatusOK, response)
}

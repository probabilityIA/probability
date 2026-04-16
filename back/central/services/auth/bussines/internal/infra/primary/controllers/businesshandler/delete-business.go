package businesshandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesshandler/mapper"
)

// DeleteBusiness godoc
//
//	@Summary		Eliminar negocio
//	@Description	Elimina un negocio del sistema
//	@Tags			businesses
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int						true	"ID del negocio"
//	@Success		201	{object}	map[string]interface{}	"Negocio eliminado exitosamente"
//	@Failure		400	{object}	map[string]interface{}	"Solicitud inválida"
//	@Failure		401	{object}	map[string]interface{}	"Token de acceso requerido"
//	@Failure		500	{object}	map[string]interface{}	"Error interno del servidor"
//	@Router			/businesses/{id} [delete]
func (h *BusinessHandler) DeleteBusinessHandler(c *gin.Context) {
	// Obtener ID del path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("invalid_id", "ID de negocio inválido"))
		return
	}

	// Ejecutar caso de uso
	err = h.usecase.DeleteBusiness(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "negocio no encontrado" {
			c.JSON(http.StatusNotFound, mapper.BuildErrorResponse("not_found", "Negocio no encontrado"))
			return
		}
		h.logger.Error().Err(err).Uint("id", uint(id)).Msg("Error al eliminar negocio")
		c.JSON(http.StatusInternalServerError, mapper.BuildErrorResponse("internal_error", "Error interno del servidor"))
		return
	}

	// Construir respuesta exitosa
	response := mapper.BuildDeleteBusinessResponse("Negocio eliminado exitosamente")
	c.JSON(http.StatusOK, response)
}

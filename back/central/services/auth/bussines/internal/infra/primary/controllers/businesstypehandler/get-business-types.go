package businesstypehandler

import (
	"net/http"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesstypehandler/mapper"

	"github.com/gin-gonic/gin"
)

// GetBusinessTypes godoc
//
//	@Summary		Obtener lista de tipos de negocio
//	@Description	Obtiene una lista de todos los tipos de negocio del sistema
//	@Tags			business-types
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		201	{object}	map[string]interface{}	"Tipos de negocio obtenidos exitosamente"
//	@Failure		400	{object}	map[string]interface{}	"Solicitud inválida"
//	@Failure		401	{object}	map[string]interface{}	"Token de acceso requerido"
//	@Failure		500	{object}	map[string]interface{}	"Error interno del servidor"
//	@Router			/business-types [get]
func (h *BusinessTypeHandler) GetBusinessTypesHandler(c *gin.Context) {
	// Ejecutar caso de uso
	businessTypes, err := h.usecase.GetBusinessTypes(c.Request.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al obtener tipos de negocio")
		c.JSON(http.StatusInternalServerError, mapper.BuildErrorResponse("internal_error", "Error interno del servidor"))
		return
	}

	// Construir respuesta exitosa
	response := mapper.BuildGetBusinessTypesResponseFromDTOs(businessTypes, "Tipos de negocio obtenidos exitosamente")
	c.JSON(http.StatusOK, response)
}

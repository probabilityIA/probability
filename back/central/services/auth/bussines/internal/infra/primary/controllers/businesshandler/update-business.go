package businesshandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesshandler/mapper"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/infra/primary/controllers/businesshandler/request"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *BusinessHandler) UpdateBusinessHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("invalid_id", "ID de negocio inválido"))
		return
	}

	if !middleware.IsSuperAdmin(c) {
		tokenBusinessID, hasBusinessID := middleware.GetBusinessID(c)
		if !hasBusinessID || tokenBusinessID == 0 || tokenBusinessID != uint(id) {
			h.logger.Error().
				Uint("token_business_id", tokenBusinessID).
				Uint("target_business_id", uint(id)).
				Str("endpoint", "/businesses/:id").
				Str("method", "PUT").
				Msg("Intento de actualizar un negocio ajeno")
			c.JSON(http.StatusForbidden, mapper.BuildErrorResponse("forbidden", "No tienes permisos para actualizar este negocio"))
			return
		}
	}

	var updateRequest request.UpdateBusinessRequest

	if err := c.ShouldBind(&updateRequest); err != nil {
		h.logger.Error().Err(err).Str("content_type", c.GetHeader("Content-Type")).Msg("Error al parsear request de actualización")

		body, _ := c.GetRawData()
		h.logger.Error().Str("body", string(body)).Msg("Contenido del request")

		c.JSON(http.StatusBadRequest, mapper.BuildErrorResponse("invalid_request", "Datos de entrada inválidos"))
		return
	}

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

	response := mapper.BuildUpdateBusinessResponseFromDTO(business, "Negocio actualizado exitosamente")
	c.JSON(http.StatusOK, response)
}

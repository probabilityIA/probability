package rolehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/response"
)

// GetRolesByLevelHandler maneja la solicitud de obtener roles por nivel
//
//	@Summary		Obtener roles por nivel
//	@Description	Obtiene todos los roles de un nivel específico con información del scope
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			level	query		int							true	"Nivel del rol"	minimum(1)	maximum(10)
//	@Success		200		{object}	response.RoleListResponse	"Roles obtenidos exitosamente"
//	@Failure		400		{object}	response.RoleErrorResponse	"Nivel inválido"
//	@Failure		401		{object}	response.RoleErrorResponse	"Token de acceso requerido"
//	@Failure		500		{object}	response.RoleErrorResponse	"Error interno del servidor"
//	@Router			/roles/by-level [get]
func (h *RoleHandler) GetRolesByLevelHandler(c *gin.Context) {
	var req request.GetRolesByLevelRequest

	// Binding automático con validaciones para parámetros de query
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error().Err(err).Msg("Error al validar nivel del rol")
		c.JSON(http.StatusBadRequest, response.RoleErrorResponse{
			Error: "Nivel inválido: " + err.Error(),
		})
		return
	}

	h.logger.Info().Int("level", req.Level).Msg("Iniciando solicitud para obtener roles por nivel")

	filters := mapper.ToRoleFilters(req)
	roles, err := h.usecase.GetRolesByLevel(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error().Err(err).Int("level", req.Level).Msg("Error al obtener roles por nivel desde el caso de uso")
		c.JSON(http.StatusInternalServerError, response.RoleErrorResponse{
			Error: "Error interno del servidor",
		})
		return
	}

	rolesResponse := mapper.ToRoleListResponse(roles)

	h.logger.Info().Int("level", req.Level).Int("count", len(roles)).Msg("Roles obtenidos exitosamente")
	c.JSON(http.StatusOK, rolesResponse)
}

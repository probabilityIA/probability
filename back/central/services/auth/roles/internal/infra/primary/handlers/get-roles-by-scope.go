package rolehandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/response"
)

// GetRolesByScopeHandler maneja la solicitud de obtener roles por scope
//
//	@Summary		Obtener roles por scope
//	@Description	Obtiene todos los roles de un scope específico
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			scope_id	path		int							true	"ID del scope"	minimum(1)
//	@Success		200			{object}	response.RoleListResponse	"Roles por scope obtenidos exitosamente"
//	@Failure		400			{object}	response.RoleErrorResponse	"Scope ID inválido"
//	@Failure		401			{object}	response.RoleErrorResponse	"Token de acceso requerido"
//	@Failure		500			{object}	response.RoleErrorResponse	"Error interno del servidor"
//	@Router			/roles/scope/{scope_id} [get]
func (h *RoleHandler) GetRolesByScopeHandler(c *gin.Context) {
	scopeIDStr := c.Param("scope_id")
	scopeID, err := strconv.ParseUint(scopeIDStr, 10, 32)
	if err != nil {
		h.logger.Error().Str("scope_id", scopeIDStr).Err(err).Msg("Scope ID de rol inválido")
		c.JSON(http.StatusBadRequest, response.RoleErrorResponse{
			Error: "Scope ID de rol inválido",
		})
		return
	}

	h.logger.Info().Uint64("scope_id", scopeID).Msg("Iniciando solicitud para obtener roles por scope")

	roles, err := h.usecase.GetRolesByScopeID(c.Request.Context(), uint(scopeID))
	if err != nil {
		h.logger.Error().Uint64("scope_id", scopeID).Err(err).Msg("Error al obtener roles por scope desde el caso de uso")
		c.JSON(http.StatusInternalServerError, response.RoleErrorResponse{
			Error: "Error interno del servidor",
		})
		return
	}

	response := mapper.ToRoleListResponse(roles)

	h.logger.Info().Uint64("scope_id", scopeID).Int("count", len(roles)).Msg("Roles por scope obtenidos exitosamente")
	c.JSON(http.StatusOK, response)
}

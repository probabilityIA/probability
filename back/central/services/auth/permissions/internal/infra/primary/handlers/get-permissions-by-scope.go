package permissionhandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers/response"
)

// GetPermissionsByScopeHandler maneja la solicitud de obtener permisos por scope
//
//	@Summary		Obtener permisos por scope
//	@Description	Obtiene todos los permisos de un scope específico
//	@Tags			Permissions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			scope_id	path		int									true	"ID del scope"	minimum(1)
//	@Success		200			{object}	response.PermissionListResponse		"Permisos por scope obtenidos exitosamente"
//	@Failure		400			{object}	response.PermissionErrorResponse	"Scope ID inválido"
//	@Failure		401			{object}	response.PermissionErrorResponse	"Token de acceso requerido"
//	@Failure		500			{object}	response.PermissionErrorResponse	"Error interno del servidor"
//	@Router			/permissions/scope/{scope_id} [get]
func (h *PermissionHandler) GetPermissionsByScopeHandler(c *gin.Context) {
	scopeIDStr := c.Param("scope_id")
	scopeID, err := strconv.ParseUint(scopeIDStr, 10, 32)
	if err != nil {
		h.logger.Error().Str("scope_id", scopeIDStr).Err(err).Msg("Scope ID de permiso inválido")
		c.JSON(http.StatusBadRequest, response.PermissionErrorResponse{
			Error: "Scope ID de permiso inválido",
		})
		return
	}

	h.logger.Info().Uint64("scope_id", scopeID).Msg("Iniciando solicitud para obtener permisos por scope")

	permissions, err := h.usecase.GetPermissionsByScopeID(c.Request.Context(), uint(scopeID))
	if err != nil {
		h.logger.Error().Uint64("scope_id", scopeID).Err(err).Msg("Error al obtener permisos por scope desde el caso de uso")
		c.JSON(http.StatusInternalServerError, response.PermissionErrorResponse{
			Error: "Error interno del servidor",
		})
		return
	}

	response := mapper.ToPermissionListResponse(permissions)

	h.logger.Info().Uint64("scope_id", scopeID).Int("count", len(permissions)).Msg("Permisos por scope obtenidos exitosamente")
	c.JSON(http.StatusOK, response)
}

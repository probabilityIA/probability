package rolehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/response"
)

// GetSystemRolesHandler maneja la solicitud de obtener roles del sistema
//
//	@Summary		Obtener roles del sistema
//	@Description	Obtiene solo los roles del sistema (is_system = true)
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.RoleListResponse	"Roles del sistema obtenidos exitosamente"
//	@Failure		401	{object}	response.RoleErrorResponse	"Token de acceso requerido"
//	@Failure		500	{object}	response.RoleErrorResponse	"Error interno del servidor"
//	@Router			/roles/system [get]
func (h *RoleHandler) GetSystemRolesHandler(c *gin.Context) {
	h.logger.Info().Msg("Iniciando solicitud para obtener roles del sistema")

	roles, err := h.usecase.GetSystemRoles(c.Request.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al obtener roles del sistema desde el caso de uso")
		c.JSON(http.StatusInternalServerError, response.RoleErrorResponse{
			Error: "Error interno del servidor",
		})
		return
	}

	response := mapper.ToRoleListResponse(roles)

	h.logger.Info().Int("count", len(roles)).Msg("Roles del sistema obtenidos exitosamente")
	c.JSON(http.StatusOK, response)
}

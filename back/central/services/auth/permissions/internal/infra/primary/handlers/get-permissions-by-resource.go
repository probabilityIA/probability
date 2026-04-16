package permissionhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/permissions/internal/infra/primary/handlers/response"
)

// GetPermissionsByResourceHandler maneja la solicitud de obtener permisos por recurso
//
//	@Summary		Obtener permisos por recurso
//	@Description	Obtiene todos los permisos de un recurso específico
//	@Tags			Permissions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			resource	path		string								true	"Recurso de los permisos"	example("users")
//	@Success		200			{object}	response.PermissionListResponse		"Permisos por recurso obtenidos exitosamente"
//	@Failure		400			{object}	response.PermissionErrorResponse	"Recurso inválido"
//	@Failure		401			{object}	response.PermissionErrorResponse	"Token de acceso requerido"
//	@Failure		500			{object}	response.PermissionErrorResponse	"Error interno del servidor"
//	@Router			/permissions/resource/{resource} [get]
func (h *PermissionHandler) GetPermissionsByResourceHandler(c *gin.Context) {
	resource := c.Param("resource")
	if resource == "" {
		h.logger.Error().Msg("Recurso de permiso inválido")
		c.JSON(http.StatusBadRequest, response.PermissionErrorResponse{
			Error: "Recurso de permiso inválido",
		})
		return
	}

	h.logger.Info().Str("resource", resource).Msg("Iniciando solicitud para obtener permisos por recurso")

	permissions, err := h.usecase.GetPermissionsByResource(c.Request.Context(), resource)
	if err != nil {
		h.logger.Error().Str("resource", resource).Err(err).Msg("Error al obtener permisos por recurso desde el caso de uso")
		c.JSON(http.StatusInternalServerError, response.PermissionErrorResponse{
			Error: "Error interno del servidor",
		})
		return
	}

	response := mapper.ToPermissionListResponse(permissions)

	h.logger.Info().Str("resource", resource).Int("count", len(permissions)).Msg("Permisos por recurso obtenidos exitosamente")
	c.JSON(http.StatusOK, response)
}

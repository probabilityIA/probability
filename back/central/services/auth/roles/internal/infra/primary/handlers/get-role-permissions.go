package rolehandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/response"
)

// GetRolePermissions obtiene los permisos de un rol
//
// @Summary Obtener permisos de un rol
// @Description Obtiene la lista de permisos asignados a un rol
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path int true "ID del rol"
// @Success 200 {object} response.GetRolePermissionsResponse
// @Failure 400 {object} response.RoleErrorResponse "ID de rol inválido"
// @Failure 404 {object} response.RoleErrorResponse "Rol no encontrado"
// @Failure 500 {object} response.RoleErrorResponse "Error interno del servidor"
// @Router /roles/{id}/permissions [get]
//
// @Security BearerAuth
func (h *RoleHandler) GetRolePermissions(c *gin.Context) {
	// Parsear ID del rol desde la URL
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.RoleErrorResponse{
			Error: "ID de rol inválido",
		})
		return
	}

	// Obtener permisos usando el caso de uso
	permissions, err := h.usecase.GetRolePermissions(c.Request.Context(), uint(roleID))
	if err != nil {
		h.logger.Error().
			Err(err).
			Uint("role_id", uint(roleID)).
			Msg("Error al obtener permisos del rol")

		// Detectar si es error de validación
		if err.Error() == "rol no encontrado" {
			c.JSON(http.StatusNotFound, response.RoleErrorResponse{
				Error: "Rol no encontrado",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.RoleErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Convertir permisos a DTOs
	permissionDTOs := make([]mapper.PermissionDTO, len(permissions))
	for i, perm := range permissions {
		permissionDTOs[i] = mapper.PermissionToDTO(perm)
	}

	// Construir respuesta exitosa
	response := mapper.ToGetRolePermissionsResponse(uint(roleID), permissionDTOs)

	h.logger.Info().
		Uint("role_id", uint(roleID)).
		Int("permission_count", len(permissions)).
		Msg("Permisos del rol obtenidos exitosamente")

	c.JSON(http.StatusOK, response)
}

package rolehandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/response"
)

// AssignPermissionsToRole asigna permisos a un rol
//
// @Summary Asignar permisos a un rol
// @Description Asigna permisos a un rol. Solo se pueden asignar permisos que pertenezcan al mismo business_type que el rol
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path int true "ID del rol"
// @Param request body request.AssignPermissionsToRoleRequest true "IDs de permisos a asignar"
// @Success 200 {object} response.AssignPermissionsToRoleResponse
// @Failure 400 {object} response.RoleErrorResponse "Datos de entrada inválidos"
// @Failure 404 {object} response.RoleErrorResponse "Rol no encontrado"
// @Failure 500 {object} response.RoleErrorResponse "Error interno del servidor"
// @Router /roles/{id}/permissions [post]
//
// @Security BearerAuth
func (h *RoleHandler) AssignPermissionsToRole(c *gin.Context) {
	var req request.AssignPermissionsToRoleRequest

	// Parsear ID del rol desde la URL
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.RoleErrorResponse{
			Error: "ID de rol inválido",
		})
		return
	}

	// Validar request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.RoleErrorResponse{
			Error: "Datos de entrada inválidos",
		})
		return
	}

	// Asignar permisos usando el caso de uso
	err = h.usecase.AssignPermissionsToRole(c.Request.Context(), uint(roleID), req.PermissionIDs)
	if err != nil {
		h.logger.Error().
			Err(err).
			Uint("role_id", uint(roleID)).
			Msg("Error al asignar permisos al rol")

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

	// Construir respuesta exitosa
	response := mapper.ToAssignPermissionsToRoleResponse(uint(roleID), req.PermissionIDs)

	h.logger.Info().
		Uint("role_id", uint(roleID)).
		Int("permission_count", len(req.PermissionIDs)).
		Msg("Permisos asignados exitosamente al rol")

	c.JSON(http.StatusOK, response)
}

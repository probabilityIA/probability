package rolehandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/response"
)

// RemovePermissionFromRole elimina un permiso de un rol
//
// @Summary Eliminar permiso de un rol
// @Description Elimina un permiso específico de un rol
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path int true "ID del rol"
// @Param permission_id path int true "ID del permiso"
// @Success 200 {object} response.RoleSuccessResponse
// @Failure 400 {object} response.RoleErrorResponse "Datos de entrada inválidos"
// @Failure 404 {object} response.RoleErrorResponse "Rol o permiso no encontrado"
// @Failure 500 {object} response.RoleErrorResponse "Error interno del servidor"
// @Router /roles/{id}/permissions/{permission_id} [delete]
//
// @Security BearerAuth
func (h *RoleHandler) RemovePermissionFromRole(c *gin.Context) {
	// Parsear IDs desde la URL
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.RoleErrorResponse{
			Error: "ID de rol inválido",
		})
		return
	}

	permissionIDStr := c.Param("permission_id")
	permissionID, err := strconv.ParseUint(permissionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.RoleErrorResponse{
			Error: "ID de permiso inválido",
		})
		return
	}

	// Eliminar permiso usando el caso de uso
	err = h.usecase.RemovePermissionFromRole(c.Request.Context(), uint(roleID), uint(permissionID))
	if err != nil {
		h.logger.Error().
			Err(err).
			Uint("role_id", uint(roleID)).
			Uint("permission_id", uint(permissionID)).
			Msg("Error al eliminar permiso del rol")

		// Detectar si es error de validación
		if err.Error() == "rol no encontrado" || err.Error() == "permiso no encontrado" {
			c.JSON(http.StatusNotFound, response.RoleErrorResponse{
				Error: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.RoleErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Construir respuesta exitosa
	h.logger.Info().
		Uint("role_id", uint(roleID)).
		Uint("permission_id", uint(permissionID)).
		Msg("Permiso eliminado exitosamente del rol")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Permiso eliminado exitosamente del rol",
	})
}

package rolehandler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/request"
)

// UpdateRole actualiza un rol existente
// @Summary Actualizar un rol
// @Description Actualiza un rol existente en el sistema (actualización parcial)
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path int true "ID del rol a actualizar"
// @Param role body request.UpdateRoleRequest true "Datos del rol a actualizar"
// @Success 200 {object} response.UpdateRoleResponse
// @Failure 400 {object} map[string]interface{} "Datos de entrada inválidos"
// @Failure 404 {object} map[string]interface{} "Rol no encontrado"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor"
// @Router /roles/{id} [put]
// @Security BearerAuth
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	// Obtener ID del path parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID inválido",
			"error":   "El ID debe ser un número válido",
		})
		return
	}

	var req request.UpdateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

	// Convertir request a DTO
	roleDTO := mapper.ToUpdateRoleDTO(req)

	// Actualizar el rol
	role, err := h.usecase.UpdateRole(c.Request.Context(), uint(id), roleDTO)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Error al actualizar el rol"
		errMsg := err.Error()

		// Detectar errores de validación de negocio (deben loggearse como WARN)
		if errMsg == "rol no encontrado" {
			statusCode = http.StatusNotFound
			message = "Rol no encontrado"
		} else if strings.Contains(errMsg, "ya existe un rol con el nombre") {
			// Detectar error de nombre duplicado (validación previa)
			statusCode = http.StatusConflict
			message = errMsg
		} else if strings.Contains(errMsg, "duplicate key") && strings.Contains(errMsg, "uni_role_name") {
			// Fallback: Detectar error de nombre duplicado desde la BD (no debería ocurrir)
			statusCode = http.StatusConflict
			message = "Ya existe un rol con este nombre. Por favor, use un nombre diferente."
		} else if strings.Contains(errMsg, "duplicate key") && strings.Contains(errMsg, "SQLSTATE 23505") {
			// Fallback: Detectar cualquier otro error de clave duplicada desde la BD
			statusCode = http.StatusConflict
			message = "Ya existe un rol con estos datos. Por favor, verifique los valores ingresados."
		}

		// Log apropiado según el tipo de error
		name := ""
		if req.Name != nil {
			name = *req.Name
		}
		if statusCode == http.StatusConflict {
			h.logger.Warn().Str("name", name).Uint("id", uint(id)).Msg("Intento de actualizar rol con nombre duplicado")
		} else if statusCode == http.StatusNotFound {
			h.logger.Warn().Uint("id", uint(id)).Msg("Intento de actualizar rol no encontrado")
		} else {
			h.logger.Error().Err(err).Uint("id", uint(id)).Msg("Error al actualizar rol")
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": message,
			"error":   errMsg,
		})
		return
	}

	// Convertir a response
	response := mapper.ToUpdateRoleResponse(role)

	c.JSON(http.StatusOK, response)
}

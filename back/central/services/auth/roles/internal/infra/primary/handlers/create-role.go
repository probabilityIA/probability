package rolehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/request"
)

// CreateRole crea un nuevo rol
// @Summary Crear un nuevo rol
// @Description Crea un nuevo rol en el sistema con todos los campos obligatorios
// @Tags Roles
// @Accept json
// @Produce json
// @Param role body request.CreateRoleRequest true "Datos del rol a crear"
// @Success 201 {object} response.CreateRoleResponse
// @Failure 400 {object} map[string]interface{} "Datos de entrada inválidos"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor"
// @Router /roles [post]
//
//	@Security		BearerAuth
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req request.CreateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

	// Convertir request a DTO
	roleDTO := mapper.ToCreateRoleDTO(req)

	// Crear el rol
	role, err := h.usecase.CreateRole(c.Request.Context(), roleDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al crear el rol",
			"error":   err.Error(),
		})
		return
	}

	// Convertir a response
	response := mapper.ToCreateRoleResponse(role)

	c.JSON(http.StatusCreated, response)
}

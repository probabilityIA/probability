package rolehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/roles/internal/infra/primary/handlers/response"
)

// GetRoleByIDHandler maneja la solicitud de obtener un rol por ID
//
//	@Summary		Obtener rol por ID
//	@Description	Obtiene un rol específico por su ID con información del scope
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int								true	"ID del rol"	minimum(1)
//	@Success		200	{object}	response.RoleSuccessResponse	"Rol obtenido exitosamente"
//	@Failure		400	{object}	response.RoleErrorResponse		"ID inválido"
//	@Failure		401	{object}	response.RoleErrorResponse		"Token de acceso requerido"
//	@Failure		404	{object}	response.RoleErrorResponse		"Rol no encontrado"
//	@Failure		500	{object}	response.RoleErrorResponse		"Error interno del servidor"
//	@Router			/roles/{id} [get]
func (h *RoleHandler) GetRoleByIDHandler(c *gin.Context) {
	var req request.GetRoleByIDRequest

	// Binding automático con validaciones para parámetros de URL
	if err := c.ShouldBindUri(&req); err != nil {
		h.logger.Error().Err(err).Msg("Error al validar ID del rol")
		c.JSON(http.StatusBadRequest, response.RoleErrorResponse{
			Error: "ID inválido: " + err.Error(),
		})
		return
	}

	h.logger.Info().Uint("id", req.ID).Msg("Iniciando solicitud para obtener rol por ID")

	role, err := h.usecase.GetRoleByID(c.Request.Context(), req.ID)
	if err != nil {
		h.logger.Error().Err(err).Uint("id", req.ID).Msg("Error al obtener rol por ID desde el caso de uso")

		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "rol no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = "Rol no encontrado"
		}

		c.JSON(statusCode, response.RoleErrorResponse{
			Error: errorMessage,
		})
		return
	}

	roleResponse := mapper.ToRoleResponse(*role)

	h.logger.Info().Uint("id", req.ID).Msg("Rol obtenido exitosamente")
	c.JSON(http.StatusOK, response.RoleSuccessResponse{
		Success: true,
		Data:    roleResponse,
	})
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/response"
)

// Deletehandlers maneja la solicitud de eliminar un usuario
//
//	@Summary		Eliminar usuario
//	@Description	Elimina un usuario del sistema
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int								true	"ID del usuario"	minimum(1)
//	@Success		200	{object}	response.UserMessageResponse	"Usuario eliminado exitosamente"
//	@Failure		400	{object}	response.UserErrorResponse		"ID inválido"
//	@Failure		401	{object}	response.UserErrorResponse		"Token de acceso requerido"
//	@Failure		404	{object}	response.UserErrorResponse		"Usuario no encontrado"
//	@Failure		500	{object}	response.UserErrorResponse		"Error interno del servidor"
//	@Router			/users/{id} [delete]
func (h *handlers) Deletehandlers(c *gin.Context) {
	var req request.DeleteUserRequest

	// Binding automático con validaciones para parámetros de URL
	if err := c.ShouldBindUri(&req); err != nil {
		h.logger.Error().Err(err).Msg("Error al validar ID del usuario")
		c.JSON(http.StatusBadRequest, response.UserErrorResponse{
			Error: "ID inválido: " + err.Error(),
		})
		return
	}

	h.logger.Info().Uint("id", req.ID).Msg("Iniciando solicitud para eliminar usuario")

	message, err := h.usecase.DeleteUser(c.Request.Context(), req.ID)
	if err != nil {
		h.logger.Error().Err(err).Uint("id", req.ID).Msg("Error al eliminar usuario desde el caso de uso")

		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "usuario no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = "Usuario no encontrado"
		}

		c.JSON(statusCode, response.UserErrorResponse{
			Error: errorMessage,
		})
		return
	}

	h.logger.Info().Uint("id", req.ID).Msg("Usuario eliminado exitosamente")
	c.JSON(http.StatusOK, response.UserMessageResponse{
		Success: true,
		Message: message,
	})
}

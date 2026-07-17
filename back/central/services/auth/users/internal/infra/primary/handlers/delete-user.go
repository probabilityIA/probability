package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/users/internal/infra/primary/handlers/response"
)

func (h *handlers) Deletehandlers(c *gin.Context) {
	var req request.DeleteUserRequest

	if err := c.ShouldBindUri(&req); err != nil {
		h.logger.Error().Err(err).Msg("Error al validar ID del usuario")
		c.JSON(http.StatusBadRequest, response.UserErrorResponse{
			Error: "ID inválido: " + err.Error(),
		})
		return
	}

	if !middleware.IsSuperAdmin(c) {
		tokenBusinessID, ok := middleware.GetBusinessID(c)
		if !ok || tokenBusinessID == 0 {
			c.JSON(http.StatusForbidden, response.UserErrorResponse{
				Error: "No tienes permisos para eliminar este usuario",
			})
			return
		}
		businesses, err := h.usecase.GetUserBusinesses(c.Request.Context(), req.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, response.UserErrorResponse{
				Error: "Usuario no encontrado",
			})
			return
		}
		belongs := false
		for _, b := range businesses {
			if b.ID == tokenBusinessID {
				belongs = true
				break
			}
		}
		if !belongs {
			h.logger.Error().
				Uint("requested_user_id", req.ID).
				Uint("token_business_id", tokenBusinessID).
				Str("endpoint", "/users/:id").
				Str("method", "DELETE").
				Msg("Intento de eliminar un usuario de otro negocio")
			c.JSON(http.StatusForbidden, response.UserErrorResponse{
				Error: "No tienes permisos para eliminar este usuario",
			})
			return
		}
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

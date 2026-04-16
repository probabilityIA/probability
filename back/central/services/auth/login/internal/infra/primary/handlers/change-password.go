package authhandler

import (
	"net/http"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// ChangePasswordHandler maneja la solicitud de cambio de contraseña
//
//	@Summary		Cambiar contraseña
//	@Description	Cambia la contraseña del usuario autenticado
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			password	body		request.ChangePasswordRequest	true	"Datos para cambiar contraseña"
//	@Success		200			{object}	response.ChangePasswordResponse	"Contraseña cambiada exitosamente"
//	@Failure		400			{object}	response.LoginErrorResponse		"Datos inválidos"
//	@Failure		401			{object}	response.LoginErrorResponse		"Token de acceso requerido"
//	@Failure		403			{object}	response.LoginErrorResponse		"Contraseña actual incorrecta"
//	@Failure		500			{object}	response.LoginErrorResponse		"Error interno del servidor"
//	@Router			/auth/change-password [post]
func (h *AuthHandler) ChangePasswordHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "ChangePasswordHandler")

	var changePasswordRequest request.ChangePasswordRequest

	// Validar y bindear el request
	if err := c.ShouldBindJSON(&changePasswordRequest); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al validar request de cambio de contraseña")
		c.JSON(http.StatusBadRequest, response.LoginErrorResponse{
			Error: "Datos de entrada inválidos: " + err.Error(),
		})
		return
	}

	// Obtener userID del token JWT (asumiendo que está en el contexto del middleware de auth)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error(ctx).Msg("UserID no encontrado en el contexto")
		c.JSON(http.StatusUnauthorized, response.LoginErrorResponse{
			Error: "Token de acceso inválido",
		})
		return
	}

	// Convertir userID a uint
	userIDUint, ok := userID.(uint)
	if !ok {
		h.logger.Error(ctx).Msg("UserID en contexto no es del tipo correcto")
		c.JSON(http.StatusUnauthorized, response.LoginErrorResponse{
			Error: "Token de acceso inválido",
		})
		return
	}

	h.logger.Info(ctx).Uint("user_id", userIDUint).Msg("Iniciando cambio de contraseña")

	// Convertir request a dominio
	domainRequest := domain.ChangePasswordRequest{
		UserID:          userIDUint,
		CurrentPassword: changePasswordRequest.CurrentPassword,
		NewPassword:     changePasswordRequest.NewPassword,
	}

	// Ejecutar caso de uso
	domainResponse, err := h.usecase.ChangePassword(ctx, domainRequest)
	if err != nil {
		h.logger.Error(ctx).Err(err).Uint("user_id", userIDUint).Msg("Error en proceso de cambio de contraseña")

		// Determinar el código de estado HTTP apropiado
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "contraseña actual incorrecta" {
			statusCode = http.StatusForbidden
			errorMessage = "Contraseña actual incorrecta"
		} else if err.Error() == "usuario no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = "Usuario no encontrado"
		} else if err.Error() == "usuario inactivo" {
			statusCode = http.StatusForbidden
			errorMessage = "Usuario inactivo"
		} else if err.Error() == "la nueva contraseña debe ser diferente a la actual" {
			statusCode = http.StatusBadRequest
			errorMessage = "La nueva contraseña debe ser diferente a la actual"
		}

		c.JSON(statusCode, response.LoginErrorResponse{
			Error: errorMessage,
		})
		return
	}

	h.logger.Info(ctx).Uint("user_id", userIDUint).Msg("Contraseña cambiada exitosamente")

	// Retornar respuesta exitosa
	c.JSON(http.StatusOK, response.ChangePasswordResponse{
		Success: domainResponse.Success,
		Message: domainResponse.Message,
	})
}

package authhandler

import (
	"net/http"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) ResetPasswordHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "ResetPasswordHandler")

	var req request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al validar request de restablecer contrasena")
		c.JSON(http.StatusBadRequest, response.LoginErrorResponse{
			Error: "Datos de entrada invalidos: " + err.Error(),
		})
		return
	}

	domainResponse, err := h.usecase.ResetPassword(ctx, domain.ResetPasswordRequest{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error en proceso de restablecer contrasena")
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"
		switch err.Error() {
		case "token invalido", "token invalido o expirado":
			statusCode = http.StatusBadRequest
			errorMessage = "El enlace es invalido o ha expirado. Solicita uno nuevo."
		case "la contrasena debe tener al menos 6 caracteres":
			statusCode = http.StatusBadRequest
			errorMessage = "La contrasena debe tener al menos 6 caracteres"
		}
		c.JSON(statusCode, response.LoginErrorResponse{Error: errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": domainResponse.Success,
		"message": domainResponse.Message,
	})
}

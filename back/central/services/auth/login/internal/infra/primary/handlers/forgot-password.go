package authhandler

import (
	"net/http"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) ForgotPasswordHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "ForgotPasswordHandler")

	var req request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al validar request de recuperacion de contrasena")
		c.JSON(http.StatusBadRequest, response.LoginErrorResponse{
			Error: "Datos de entrada invalidos: " + err.Error(),
		})
		return
	}

	domainResponse, err := h.usecase.ForgotPassword(ctx, domain.ForgotPasswordRequest{Email: req.Email, Channel: req.Channel})
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error en proceso de recuperacion de contrasena")
		c.JSON(http.StatusInternalServerError, response.LoginErrorResponse{
			Error: "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": domainResponse.Success,
		"message": domainResponse.Message,
	})
}

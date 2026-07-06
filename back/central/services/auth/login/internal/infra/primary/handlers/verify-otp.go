package authhandler

import (
	"net/http"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) VerifyOTPHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "VerifyOTPHandler")

	var req request.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al validar request de verificacion OTP")
		c.JSON(http.StatusBadRequest, response.LoginErrorResponse{
			Error: "Datos de entrada invalidos: " + err.Error(),
		})
		return
	}

	domainResponse, err := h.usecase.VerifyOTP(ctx, domain.VerifyOTPRequest{Email: req.Email, Code: req.Code})
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error en verificacion OTP")
		c.JSON(http.StatusInternalServerError, response.LoginErrorResponse{
			Error: "Error interno del servidor",
		})
		return
	}

	if !domainResponse.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": domainResponse.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": domainResponse.Success,
		"message": domainResponse.Message,
		"token":   domainResponse.Token,
	})
}

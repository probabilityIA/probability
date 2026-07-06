package authhandler

import (
	"net/http"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/auth/login/internal/infra/primary/handlers/response"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) RecoveryChannelsHandler(c *gin.Context) {
	ctx := log.WithFunctionCtx(c.Request.Context(), "RecoveryChannelsHandler")

	var req request.RecoveryChannelsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error al validar request de canales de recuperacion")
		c.JSON(http.StatusBadRequest, response.LoginErrorResponse{
			Error: "Datos de entrada invalidos: " + err.Error(),
		})
		return
	}

	domainResponse, err := h.usecase.RecoveryChannels(ctx, domain.RecoveryChannelsRequest{Email: req.Email})
	if err != nil {
		h.logger.Error(ctx).Err(err).Msg("Error obteniendo canales de recuperacion")
		c.JSON(http.StatusInternalServerError, response.LoginErrorResponse{
			Error: "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email": domainResponse.Email,
		"whatsapp": gin.H{
			"available":    domainResponse.WhatsApp.Available,
			"masked_phone": domainResponse.WhatsApp.MaskedPhone,
		},
	})
}

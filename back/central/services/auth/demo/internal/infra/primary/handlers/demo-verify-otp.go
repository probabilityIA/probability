package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/demo/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/demo/internal/infra/primary/handlers/request"
)

func (h *Handler) DemoVerifyOTPHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.DemoVerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Datos invalidos en verificacion OTP demo")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada invalidos: " + err.Error()})
		return
	}

	resp, err := h.usecase.DemoVerifyOTP(ctx, domain.DemoVerifyOTPRequest{Email: req.Email, Code: req.Code})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": resp.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": resp.Success, "message": resp.Message})
}

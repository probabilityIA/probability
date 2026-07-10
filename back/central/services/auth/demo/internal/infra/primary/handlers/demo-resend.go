package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/demo/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/demo/internal/infra/primary/handlers/request"
)

func (h *Handler) DemoResendHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.DemoResendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Datos invalidos en reenvio de verificacion demo")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada invalidos: " + err.Error()})
		return
	}

	resp, err := h.usecase.DemoResend(ctx, domain.DemoResendRequest{
		Email:   req.Email,
		Channel: req.Channel,
		Phone:   req.Phone,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "el telefono es obligatorio para verificar por WhatsApp" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": resp.Success, "message": resp.Message})
}

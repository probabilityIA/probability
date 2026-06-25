package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/demo/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/demo/internal/infra/primary/handlers/request"
)

func (h *Handler) VerifyEmailHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Datos invalidos en verificacion de email")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada invalidos: " + err.Error()})
		return
	}

	resp, err := h.usecase.VerifyEmail(ctx, domain.VerifyEmailRequest{Token: req.Token})
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "token invalido", "token invalido o expirado":
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": "El enlace es invalido o ha expirado. Solicita uno nuevo."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": resp.Success, "message": resp.Message})
}

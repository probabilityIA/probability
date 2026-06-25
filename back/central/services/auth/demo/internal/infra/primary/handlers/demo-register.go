package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/demo/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/demo/internal/infra/primary/handlers/request"
)

func (h *Handler) DemoRegisterHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.DemoRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx).Err(err).Msg("Datos invalidos en registro demo")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada invalidos: " + err.Error()})
		return
	}

	resp, err := h.usecase.DemoRegister(ctx, domain.DemoRegisterRequest{
		FullName:     req.FullName,
		BusinessName: req.BusinessName,
		Email:        req.Email,
		Password:     req.Password,
	})
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "el correo ya esta registrado":
			status = http.StatusConflict
		case "la contrasena debe tener al menos 6 caracteres", "nombre, negocio y correo son obligatorios":
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": resp.Success, "message": resp.Message})
}

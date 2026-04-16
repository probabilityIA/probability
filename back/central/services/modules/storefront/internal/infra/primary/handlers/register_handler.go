package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/infra/primary/handlers/request"
)

func (h *Handlers) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto := mappers.RequestToRegisterDTO(&req)

	err := h.uc.Register(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, domainerrors.ErrBusinessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domainerrors.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domainerrors.ErrRoleNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "configuracion del sistema incompleta"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registro exitoso"})
}

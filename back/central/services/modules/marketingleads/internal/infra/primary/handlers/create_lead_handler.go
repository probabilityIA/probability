package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/infra/primary/handlers/request"
)

func (h *Handlers) CreateLead(c *gin.Context) {
	var req request.CreateLeadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos invalidos: " + err.Error()})
		return
	}

	if err := h.uc.CreateLead(c.Request.Context(), req.ToDTO()); err != nil {
		if errors.Is(err, domainerrors.ErrInvalidLead) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/response"
)

func (h *Handlers) CreateTax(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	var req request.CreateTaxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	tax, err := h.uc.CreateTax(c.Request.Context(), dtos.CreateTaxDTO{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		RatePercent: req.RatePercent,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrDuplicateCode):
			c.JSON(http.StatusConflict, gin.H{"success": false, "error": err.Error()})
		case errors.Is(err, domainerrors.ErrInvalidRate), errors.Is(err, domainerrors.ErrCodeRequired), errors.Is(err, domainerrors.ErrNameRequired):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": response.FromTax(tax)})
}

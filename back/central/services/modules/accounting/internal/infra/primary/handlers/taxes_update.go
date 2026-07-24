package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/response"
)

func (h *Handlers) UpdateTax(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}
	var req request.UpdateTaxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	tax, err := h.uc.UpdateTax(c.Request.Context(), dtos.UpdateTaxDTO{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
		RatePercent: req.RatePercent,
		IsActive:    isActive,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrTaxNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		case errors.Is(err, domainerrors.ErrInvalidRate), errors.Is(err, domainerrors.ErrNameRequired):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": response.FromTax(tax)})
}

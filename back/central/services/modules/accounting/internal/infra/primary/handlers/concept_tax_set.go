package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/request"
)

func (h *Handlers) SetConceptTax(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	conceptID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || conceptID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "id invalido"})
		return
	}
	taxID, err := strconv.ParseUint(c.Param("taxId"), 10, 64)
	if err != nil || taxID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "taxId invalido"})
		return
	}
	var req request.SetConceptTaxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	err = h.uc.SetConceptTax(c.Request.Context(), dtos.SetConceptTaxDTO{
		ConceptID: uint(conceptID),
		TaxID:     uint(taxID),
		IsActive:  req.IsActive,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrConceptNotFound), errors.Is(err, domainerrors.ErrTaxNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		case errors.Is(err, domainerrors.ErrConceptTaxNotAllowed):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/infra/primary/handlers/response"
)

func (h *Handlers) CreateEntry(c *gin.Context) {
	if !h.requireSuperAdmin(c) {
		return
	}
	var req request.CreateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	entryDate, err := time.Parse("2006-01-02", req.EntryDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "entry_date debe tener formato YYYY-MM-DD"})
		return
	}
	entry, err := h.uc.CreateManualEntry(c.Request.Context(), dtos.CreateEntryDTO{
		ConceptID:   req.ConceptID,
		BusinessID:  req.BusinessID,
		EntryDate:   entryDate,
		Amount:      req.Amount,
		Description: req.Description,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrConceptNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		case errors.Is(err, domainerrors.ErrInvalidAmount), errors.Is(err, domainerrors.ErrConceptInactive):
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": response.FromEntry(entry)})
}

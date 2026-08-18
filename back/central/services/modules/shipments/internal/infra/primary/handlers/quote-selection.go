package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/app/usecasequotes"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type quoteSelectionRequest struct {
	SelectedCarrier string `json:"selected_carrier"`
	SelectedIDRate  *int64 `json:"selected_id_rate"`
}

func (h *Handlers) SetSavedQuoteSelection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}

	var req quoteSelectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	businessID, ok := h.resolveBusinessIDParam(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se pudo identificar la empresa"})
		return
	}

	a := h.resolveActor(c)

	resp, err := h.uc.Quotes.SetSelection(c.Request.Context(), domain.QuoteSelectionInput{
		QuoteID:         uint(id),
		BusinessID:      businessID,
		UpdatedBy:       a.ID,
		UpdatedByName:   a.Name,
		SelectedCarrier: req.SelectedCarrier,
		SelectedIDRate:  req.SelectedIDRate,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecasequotes.ErrQuoteNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, usecasequotes.ErrQuoteBusinessMatch):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

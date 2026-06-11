package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/app/usecasequotes"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type associateQuoteRequest struct {
	OrderUUID       string `json:"order_uuid" binding:"required"`
	SelectedCarrier string `json:"selected_carrier"`
	SelectedIDRate  *int64 `json:"selected_id_rate"`
	GuideRequested  bool   `json:"guide_requested"`
}

func (h *Handlers) AssociateSavedQuote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}

	var req associateQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	businessID, ok := h.resolveBusinessIDParam(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se pudo identificar la empresa"})
		return
	}

	resp, err := h.uc.Quotes.Associate(c.Request.Context(), domain.AssociateQuoteInput{
		QuoteID:         uint(id),
		BusinessID:      businessID,
		OrderUUID:       req.OrderUUID,
		SelectedCarrier: req.SelectedCarrier,
		SelectedIDRate:  req.SelectedIDRate,
		GuideRequested:  req.GuideRequested,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecasequotes.ErrQuoteNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, usecasequotes.ErrQuoteBusinessMatch), errors.Is(err, usecasequotes.ErrQuoteAlreadyLinked):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
)

type detailRequest struct {
	SKU   string `json:"sku"`
	Label string `json:"label"`
	Tone  string `json:"tone"`
	Group string `json:"group"`
}

type recordRequest struct {
	IntegrationID uint   `json:"integration_id" binding:"required"`
	BusinessID    *uint  `json:"business_id"`
	Kind          string `json:"kind" binding:"required"`
	CorrelationID string `json:"correlation_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`

	Total     int `json:"total"`
	Updated   int `json:"updated"`
	Unchanged int `json:"unchanged"`
	Skipped   int `json:"skipped"`
	Failed    int `json:"failed"`

	Matched           int `json:"matched"`
	NotAssociated     int `json:"not_associated"`
	OnlyInProbability int `json:"only_in_probability"`
	OnlyInChannel     int `json:"only_in_channel"`

	Detail []detailRequest `json:"detail"`
}

func (h *handler) Record(c *gin.Context) {
	var req recordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "integration_id y kind son requeridos"})
		return
	}

	if req.Kind != domain.KindInventory && req.Kind != domain.KindProducts {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": domain.ErrInvalidKind.Error()})
		return
	}

	businessID, ok := h.resolveBusinessID(c, req.BusinessID)
	if !ok {
		return
	}

	now := time.Now()
	run := domain.SyncRun{
		BusinessID:        businessID,
		IntegrationID:     req.IntegrationID,
		Kind:              req.Kind,
		CorrelationID:     req.CorrelationID,
		StartedAt:         now,
		FinishedAt:        &now,
		Total:             req.Total,
		Updated:           req.Updated,
		Unchanged:         req.Unchanged,
		Skipped:           req.Skipped,
		Failed:            req.Failed,
		Matched:           req.Matched,
		NotAssociated:     req.NotAssociated,
		OnlyInProbability: req.OnlyInProbability,
		OnlyInChannel:     req.OnlyInChannel,
		Status:            req.Status,
		Message:           req.Message,
	}
	for _, d := range req.Detail {
		run.Detail = append(run.Detail, domain.DetailItem{SKU: d.SKU, Label: d.Label, Tone: d.Tone, Group: d.Group})
	}

	if err := h.useCase.Record(c.Request.Context(), &run); err != nil {
		if errors.Is(err, domain.ErrIntegrationNotOwned) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "no se pudo guardar el resultado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"id": run.ID}})
}

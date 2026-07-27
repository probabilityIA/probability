package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
)

type detailResponse struct {
	SKU   string `json:"sku"`
	Label string `json:"label"`
	Tone  string `json:"tone"`
}

type runResponse struct {
	IntegrationID uint   `json:"integration_id"`
	Kind          string `json:"kind"`
	Status        string `json:"status"`
	Message       string `json:"message,omitempty"`
	FinishedAt    string `json:"finished_at,omitempty"`

	Total     int `json:"total"`
	Updated   int `json:"updated"`
	Unchanged int `json:"unchanged"`
	Skipped   int `json:"skipped"`
	Failed    int `json:"failed"`

	Matched           int `json:"matched"`
	NotAssociated     int `json:"not_associated"`
	OnlyInProbability int `json:"only_in_probability"`
	OnlyInChannel     int `json:"only_in_channel"`

	Detail []detailResponse `json:"detail"`
}

func (h *handler) ListLast(c *gin.Context) {
	var provided *uint
	if raw := c.Query("business_id"); raw != "" {
		if parsed, err := strconv.ParseUint(raw, 10, 64); err == nil {
			value := uint(parsed)
			provided = &value
		}
	}

	businessID, ok := h.resolveBusinessID(c, provided)
	if !ok {
		return
	}

	runs, err := h.useCase.ListLast(c.Request.Context(), businessID)
	if err != nil {
		h.logger.Error(c.Request.Context()).Err(err).Msg("Error al listar los resultados de sincronizacion")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "no se pudieron obtener los resultados"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": toRunResponses(runs)})
}

func toRunResponses(runs []domain.SyncRun) []runResponse {
	out := make([]runResponse, 0, len(runs))
	for _, run := range runs {
		item := runResponse{
			IntegrationID:     run.IntegrationID,
			Kind:              run.Kind,
			Status:            run.Status,
			Message:           run.Message,
			Total:             run.Total,
			Updated:           run.Updated,
			Unchanged:         run.Unchanged,
			Skipped:           run.Skipped,
			Failed:            run.Failed,
			Matched:           run.Matched,
			NotAssociated:     run.NotAssociated,
			OnlyInProbability: run.OnlyInProbability,
			OnlyInChannel:     run.OnlyInChannel,
			Detail:            make([]detailResponse, 0, len(run.Detail)),
		}
		if run.FinishedAt != nil {
			item.FinishedAt = run.FinishedAt.Format("2006-01-02T15:04:05Z07:00")
		}
		for _, d := range run.Detail {
			item.Detail = append(item.Detail, detailResponse{SKU: d.SKU, Label: d.Label, Tone: d.Tone})
		}
		out = append(out, item)
	}
	return out
}

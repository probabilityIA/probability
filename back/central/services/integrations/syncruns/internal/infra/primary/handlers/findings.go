package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type findingResponse struct {
	Code     string   `json:"code"`
	Severity string   `json:"severity"`
	Title    string   `json:"title"`
	Detail   string   `json:"detail"`
	Count    int      `json:"count"`
	Channels []string `json:"channels,omitempty"`
}

type channelSummaryResponse struct {
	IntegrationID   uint   `json:"integration_id"`
	IntegrationName string `json:"integration_name"`
	ChannelCode     string `json:"channel_code"`
	Matched         int    `json:"matched"`
	NotAssociated   int    `json:"not_associated"`
	OnlyInChannel   int    `json:"only_in_channel"`
	ChannelNoSKU    int    `json:"channel_no_sku"`
	SKUChanged      int    `json:"sku_changed"`
	SKUTypo         int    `json:"sku_typo"`
	SKUSpacing      int    `json:"sku_spacing"`
	ComparedAt      string `json:"compared_at,omitempty"`
}

type findingsResponse struct {
	Total    int                      `json:"total"`
	Findings []findingResponse        `json:"findings"`
	Channels []channelSummaryResponse `json:"channels"`
}

func (h *handler) Findings(c *gin.Context) {
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

	report, err := h.useCase.Findings(c.Request.Context(), businessID)
	if err != nil {
		h.logger.Error(c.Request.Context()).Err(err).
			Uint("business_id", businessID).
			Msg("Error al calcular los hallazgos de integraciones")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "no se pudieron calcular los hallazgos"})
		return
	}

	resp := findingsResponse{
		Total:    report.Total,
		Findings: make([]findingResponse, 0, len(report.Findings)),
		Channels: make([]channelSummaryResponse, 0, len(report.Channels)),
	}
	for _, f := range report.Findings {
		resp.Findings = append(resp.Findings, findingResponse{
			Code: f.Code, Severity: f.Severity, Title: f.Title,
			Detail: f.Detail, Count: f.Count, Channels: f.Channels,
		})
	}
	for _, ch := range report.Channels {
		resp.Channels = append(resp.Channels, channelSummaryResponse{
			IntegrationID: ch.IntegrationID, IntegrationName: ch.IntegrationName,
			ChannelCode: ch.ChannelCode, Matched: ch.Matched, NotAssociated: ch.NotAssociated,
			OnlyInChannel: ch.OnlyInChannel, ChannelNoSKU: ch.ChannelNoSKU,
			SKUChanged: ch.SKUChanged, SKUTypo: ch.SKUTypo, SKUSpacing: ch.SKUSpacing,
			ComparedAt: ch.ComparedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

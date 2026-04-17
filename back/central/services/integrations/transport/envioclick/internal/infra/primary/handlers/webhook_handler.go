package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

func (h *Handlers) Webhook(c *gin.Context) {
	ctx := c.Request.Context()

	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("Failed to read webhook body")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Failed to read request body",
		})
		return
	}

	var payload domain.WebhookPayload
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		h.log.Warn(ctx).Err(err).Str("body_preview", previewBody(rawBody)).Msg("Invalid webhook payload")
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Payload invalido, ignorado",
		})
		return
	}

	if payload.TrackingCode == "" {
		h.log.Warn(ctx).Msg("Webhook without trackingCode")
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "trackingCode ausente, ignorado",
		})
		return
	}

	result, err := h.uc.Process(ctx, app.WebhookRequest{
		URL:      c.Request.URL.String(),
		Method:   c.Request.Method,
		RemoteIP: c.ClientIP(),
		Headers:  c.Request.Header,
		RawBody:  rawBody,
		Payload:  payload,
	})
	if err != nil {
		h.log.Error(ctx).Err(err).Str("tracking_code", payload.TrackingCode).Msg("Webhook processing failed")
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   "Error procesando webhook",
		})
		return
	}

	if result.IsIgnored {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": result.IgnoredReason,
		})
		return
	}

	h.log.Info(ctx).
		Str("tracking_number", result.TrackingNumber).
		Str("probability_status", result.ProbabilityStatus.String()).
		Str("raw_status", result.RawStatusStep).
		Bool("is_unknown", result.IsUnknownStatus).
		Str("correlation_id", result.CorrelationID).
		Msg("Webhook processed")

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"message":        "Webhook procesado correctamente",
		"tracking":       result.TrackingNumber,
		"status":         result.ProbabilityStatus.String(),
		"correlation_id": result.CorrelationID,
	})
}

func previewBody(body []byte) string {
	const maxLen = 200
	if len(body) <= maxLen {
		return string(body)
	}
	return string(body[:maxLen]) + "..."
}

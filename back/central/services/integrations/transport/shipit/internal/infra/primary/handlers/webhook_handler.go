package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
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

	if payload.TrackingNumber == "" {
		h.log.Warn(ctx).Str("reference", payload.Reference).Msg("Webhook without tracking_number")
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "tracking_number ausente, ignorado",
		})
		return
	}

	if h.rabbit != nil {
		msg := consumer.WebhookQueueMessage{
			URL:      c.Request.URL.String(),
			Method:   c.Request.Method,
			RemoteIP: c.ClientIP(),
			Headers:  c.Request.Header,
			RawBody:  rawBody,
		}
		if data, err := json.Marshal(msg); err == nil {
			if err := h.rabbit.Publish(context.Background(), rabbitmq.QueueWebhooksShipitReceived, data); err == nil {
				c.JSON(http.StatusOK, gin.H{
					"success":  true,
					"message":  "Webhook recibido",
					"tracking": payload.TrackingNumber,
				})
				return
			}
			h.log.Warn(ctx).Str("tracking_number", payload.TrackingNumber).
				Msg("No se pudo encolar el webhook, procesando inline como fallback")
		}
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
		h.log.Error(ctx).Err(err).Str("tracking_number", payload.TrackingNumber).Msg("Webhook processing failed")
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

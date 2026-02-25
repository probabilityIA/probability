package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/infra/primary/handlers/request"
)

// WebhookGrafana maneja el webhook de alertas de Grafana Cloud
func (h *handler) WebhookGrafana(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. Leer el body raw (necesario para validación HMAC)
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("[Monitoring] Error leyendo body del webhook")
		c.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body"})
		return
	}

	// 2. Validar firma HMAC-SHA256
	secret := h.env.Get("GRAFANA_WEBHOOK_SECRET")
	if secret != "" {
		sig := c.GetHeader("X-Grafana-Signature-V2")
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(rawBody)
		expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(sig), []byte(expected)) {
			h.log.Warn(ctx).
				Str("received_sig", sig).
				Msg("[Monitoring] Firma HMAC inválida")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
			return
		}
	} else {
		h.log.Warn(ctx).Msg("[Monitoring] GRAFANA_WEBHOOK_SECRET vacío - omitiendo validación de firma (dev mode)")
	}

	// 3. Deserializar payload
	var req request.GrafanaWebhookRequest
	if err := json.Unmarshal(rawBody, &req); err != nil {
		h.log.Error(ctx).Err(err).Msg("[Monitoring] Error parseando payload de Grafana")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// 4. Mapear a DTO de dominio
	webhookDTO := dtos.GrafanaWebhookDTO{
		Status: req.Status,
		Title:  req.Title,
	}
	for _, a := range req.Alerts {
		webhookDTO.Alerts = append(webhookDTO.Alerts, dtos.GrafanaAlertDTO{
			Status:      a.Status,
			Labels:      a.Labels,
			Annotations: a.Annotations,
			StartsAt:    a.StartsAt,
			ValueString: a.ValueString,
		})
	}

	// 5. Procesar en el use case
	if err := h.useCase.ProcessGrafanaAlert(ctx, webhookDTO); err != nil {
		// Log del error pero igual retornar 200 para evitar reenvíos de Grafana
		h.log.Error(ctx).Err(err).Msg("[Monitoring] Error procesando alerta de Grafana")
	}

	// 6. Siempre retornar 200 para que Grafana no reintente en loop
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

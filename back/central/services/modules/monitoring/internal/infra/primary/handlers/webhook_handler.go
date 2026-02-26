package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/infra/primary/handlers/request"
)

// WebhookGrafana maneja el webhook de alertas de Grafana Cloud
func (h *handler) WebhookGrafana(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. Validar Bearer token
	// Grafana envía: Authorization: Bearer <GRAFANA_WEBHOOK_SECRET>
	secret := h.env.Get("GRAFANA_WEBHOOK_SECRET")
	if secret != "" {
		authHeader := c.GetHeader("Authorization")
		expectedAuth := "Bearer " + secret
		if !strings.EqualFold(strings.TrimSpace(authHeader), strings.TrimSpace(expectedAuth)) {
			h.log.Warn(ctx).Msg("[Monitoring] Authorization header inválido")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
	} else {
		h.log.Warn(ctx).Msg("[Monitoring] GRAFANA_WEBHOOK_SECRET vacío - sin validación (dev mode)")
	}

	// 2. Leer y parsear body
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("[Monitoring] Error leyendo body del webhook")
		c.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body"})
		return
	}

	var req request.GrafanaWebhookRequest
	if err := json.Unmarshal(rawBody, &req); err != nil {
		h.log.Error(ctx).Err(err).Msg("[Monitoring] Error parseando payload de Grafana")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// 3. Mapear a DTO de dominio
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

	// 4. Procesar en el use case
	if err := h.useCase.ProcessGrafanaAlert(ctx, webhookDTO); err != nil {
		// Log del error pero retornar 200 para evitar reintentos de Grafana
		h.log.Error(ctx).Err(err).Msg("[Monitoring] Error procesando alerta de Grafana")
	}

	// 5. Siempre 200 para que Grafana no reintente en loop
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

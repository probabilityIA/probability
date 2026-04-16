package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/domain/entities"
)

// alertTypeMap mapea alertnames de Grafana a tipos legibles
var alertTypeMap = map[string]string{
	"RAM_ALTA":   "RAM",
	"CPU_ALTO":   "CPU",
	"DISCO_ALTO": "Disco",
}

// ProcessGrafanaAlert procesa el payload del webhook de Grafana Cloud
func (uc *useCase) ProcessGrafanaAlert(ctx context.Context, dto dtos.GrafanaWebhookDTO) error {
	uc.log.Info(ctx).
		Str("status", dto.Status).
		Str("title", dto.Title).
		Int("alerts_count", len(dto.Alerts)).
		Msg("[Monitoring] Processing Grafana alert webhook")

	uc.log.Debug(ctx).
		Interface("payload", dto).
		Msg("📋 [Monitoring] Complete Grafana webhook payload")

	for _, alert := range dto.Alerts {
		if alert.Status != "firing" {
			continue
		}

		alertName := alert.Labels["alertname"]
		alertType := alertTypeMap[alertName]
		if alertType == "" {
			alertType = alertName
		}

		summary := alert.Annotations["summary"]
		if summary == "" {
			summary = alert.ValueString
		}

		event := entities.AlertEvent{
			AlertType: alertType,
			Summary:   summary,
			Status:    alert.Status,
			FiredAt:   alert.StartsAt,
		}

		uc.log.Info(ctx).
			Str("alert_name", alertName).
			Str("alert_type", alertType).
			Str("summary", summary).
			Msg("[Monitoring] Publishing alert event")

		uc.log.Debug(ctx).
			Interface("event", event).
			Msg("📤 [Monitoring] Alert event to be published")

		if err := uc.publisher.Publish(ctx, event); err != nil {
			uc.log.Error(ctx).
				Err(err).
				Str("alert_name", alertName).
				Msg("[Monitoring] Error publishing alert event")
			// No retornamos error para evitar que Grafana reintente en loop
		}
	}

	return nil
}

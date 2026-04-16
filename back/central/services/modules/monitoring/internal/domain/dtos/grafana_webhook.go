package dtos

import "time"

// GrafanaWebhookDTO representa el payload del webhook de Grafana Cloud
type GrafanaWebhookDTO struct {
	Status string           `json:"status"`
	Title  string           `json:"title"`
	Alerts []GrafanaAlertDTO `json:"alerts"`
}

// GrafanaAlertDTO representa una alerta individual en el payload de Grafana
type GrafanaAlertDTO struct {
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    time.Time         `json:"startsAt"`
	ValueString string            `json:"valueString"`
}

package request

import "time"

// GrafanaWebhookRequest es el struct para binding del body del webhook de Grafana
type GrafanaWebhookRequest struct {
	Status string                `json:"status"`
	Title  string                `json:"title"`
	Alerts []GrafanaAlertRequest `json:"alerts"`
}

// GrafanaAlertRequest representa una alerta individual en el payload
type GrafanaAlertRequest struct {
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    time.Time         `json:"startsAt"`
	ValueString string            `json:"valueString"`
}

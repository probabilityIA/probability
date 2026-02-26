package entities

import "time"

// GrafanaAlert representa una alerta individual de Grafana Cloud
type GrafanaAlert struct {
	Status      string
	AlertName   string
	Summary     string
	Description string
	FiredAt     time.Time
}

// AlertEvent es el evento publicado en la cola monitoring.alerts
type AlertEvent struct {
	AlertType string    `json:"alert_type"`
	Summary   string    `json:"summary"`
	Status    string    `json:"status"`
	FiredAt   time.Time `json:"fired_at"`
}

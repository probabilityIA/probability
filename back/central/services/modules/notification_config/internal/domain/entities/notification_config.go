package entities

import (
	"encoding/json"
	"time"
)

// IntegrationNotificationConfig representa la configuración de notificaciones para una integración
type IntegrationNotificationConfig struct {
	ID               uint
	IntegrationID    uint
	NotificationType string // "whatsapp", "email", "sms"
	IsActive         bool
	Conditions       NotificationConditions
	Config           NotificationConfig
	Description      string
	Priority         int // Mayor prioridad se evalúa primero
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NotificationConditions define las condiciones para disparar una notificación
type NotificationConditions struct {
	Trigger        string   // "order.created", "order.updated", "order.status_changed"
	Statuses       []string // ["pending", "processing"] - vacío = todos
	PaymentMethods []uint   // [1, 3, 5] - vacío = todos
}

// NotificationConfig contiene la configuración específica de la notificación
type NotificationConfig struct {
	TemplateName  string // "confirmacion_pedido_contraentrega"
	RecipientType string // "customer", "business"
	Language      string // "es", "en"
}

// MarshalJSON convierte NotificationConditions a JSON
func (nc NotificationConditions) MarshalJSON() ([]byte, error) {
	type Alias NotificationConditions
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&nc),
	})
}

// UnmarshalJSON convierte JSON a NotificationConditions
func (nc *NotificationConditions) UnmarshalJSON(data []byte) error {
	type Alias NotificationConditions
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(nc),
	}
	return json.Unmarshal(data, &aux)
}

// MarshalJSON convierte NotificationConfig a JSON
func (ncfg NotificationConfig) MarshalJSON() ([]byte, error) {
	type Alias NotificationConfig
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&ncfg),
	})
}

// UnmarshalJSON convierte JSON a NotificationConfig
func (ncfg *NotificationConfig) UnmarshalJSON(data []byte) error {
	type Alias NotificationConfig
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(ncfg),
	}
	return json.Unmarshal(data, &aux)
}

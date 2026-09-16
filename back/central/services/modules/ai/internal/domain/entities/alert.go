package entities

import "time"

const (
	AlertSeverityInfo     = "info"
	AlertSeverityWarning  = "warning"
	AlertSeverityCritical = "critical"
)

type Alert struct {
	ID               string
	BusinessID       uint
	EventID          string
	EventType        string
	Severity         string
	Title            string
	Body             string
	DestinationKey   string
	DestinationRoute string
	ReferenceType    string
	ReferenceID      string
	CreatedAt        time.Time
	Unread           bool
}

type AlertsUnread struct {
	Count  int64
	Latest *Alert
}

type ChatIdentity struct {
	UserName     string
	BusinessName string
	IsSuperAdmin bool
}

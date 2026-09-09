package entities

import "time"

type DeviceToken struct {
	ID         uint
	UserID     uint
	BusinessID *uint
	Token      string
	Platform   string
	AppVersion string
	DeviceName string
	IsActive   bool
	LastSeenAt *time.Time
	CreatedAt  time.Time
}

type PushMessage struct {
	Title string
	Body  string
	Data  map[string]string
}

type SendResult struct {
	Sent           int
	InvalidTokens  []string
	TransientError error
}

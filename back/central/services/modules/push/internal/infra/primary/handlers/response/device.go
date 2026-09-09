package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/push/internal/domain/entities"
)

type Device struct {
	ID         uint       `json:"id"`
	Platform   string     `json:"platform"`
	DeviceName string     `json:"device_name"`
	AppVersion string     `json:"app_version"`
	IsActive   bool       `json:"is_active"`
	LastSeenAt *time.Time `json:"last_seen_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

func FromEntities(items []entities.DeviceToken) []Device {
	out := make([]Device, 0, len(items))
	for _, item := range items {
		out = append(out, Device{
			ID:         item.ID,
			Platform:   item.Platform,
			DeviceName: item.DeviceName,
			AppVersion: item.AppVersion,
			IsActive:   item.IsActive,
			LastSeenAt: item.LastSeenAt,
			CreatedAt:  item.CreatedAt,
		})
	}
	return out
}

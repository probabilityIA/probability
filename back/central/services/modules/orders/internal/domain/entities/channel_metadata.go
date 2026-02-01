package entities

import "time"

// ProbabilityOrderChannelMetadata representa metadata del canal que se guarda en la base de datos
// ✅ ENTIDAD PURA - SIN TAGS
type ProbabilityOrderChannelMetadata struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	OrderID string

	ChannelSource string
	IntegrationID uint

	RawData []byte

	Version     string
	ReceivedAt  time.Time
	ProcessedAt *time.Time
	IsLatest    bool

	LastSyncedAt *time.Time
	SyncStatus   string
}

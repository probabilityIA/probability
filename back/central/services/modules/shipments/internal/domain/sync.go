package domain

import "time"

const (
	SyncProviderEnvioclick = "envioclick"
	SyncBatchSize          = 40
)

type SyncShipmentsFilter struct {
	BusinessID uint
	Provider   string
	DateFrom   *time.Time
	DateTo     *time.Time
	Statuses   []string
}

type SyncShipmentRow struct {
	ShipmentID        uint
	TrackingNumber    string
	Carrier           string
	EnvioclickIDOrder *int64
}

type SyncShipmentsResult struct {
	CorrelationID            string
	TotalShipments           int
	Batches                  int
	BatchSize                int
	EstimatedDurationSeconds int
}

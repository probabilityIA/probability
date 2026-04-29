package domain

import "time"

type SyncMeta struct {
	URL            string
	Method         string
	RequestBody    []byte
	ResponseStatus int
	ResponseBody   []byte
	StartedAt      time.Time
	CompletedAt    time.Time
	DurationMs     int
}

type SyncLog struct {
	ID             uint
	ShipmentID     *uint
	OperationType  string
	Provider       string
	Status         string
	RequestURL     string
	RequestMethod  string
	RequestPayload []byte
	ResponseStatus int
	ResponseBody   []byte
	ErrorMessage   *string
	ErrorCode      *string
	CorrelationID  string
	TriggeredBy    string
	UserID         *uint
	StartedAt      time.Time
	CompletedAt    *time.Time
	Duration       *int
}

package domain

type SyncBatchItem struct {
	ShipmentID        uint
	TrackingNumber    string
	EnvioclickIDOrder *int64
}

type SyncBatchRequest struct {
	BusinessID    uint
	CorrelationID string
	BaseURL       string
	APIKey        string
	URL           string
	RemoteIP      string
	Items         []SyncBatchItem
}

type SyncBatchResult struct {
	Total     int
	Processed int
	Failed    int
	Unknown   int
	NotFound  int
}

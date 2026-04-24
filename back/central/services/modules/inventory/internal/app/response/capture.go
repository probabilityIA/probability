package response

import "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"

type ScanResponse struct {
	Resolved   bool                       `json:"resolved"`
	Resolution *entities.ScanResolution   `json:"resolution,omitempty"`
	Event      *entities.ScanEvent        `json:"event,omitempty"`
}

type InboundSyncResult struct {
	Log       *entities.InventorySyncLog `json:"log"`
	Duplicate bool                       `json:"duplicate"`
}

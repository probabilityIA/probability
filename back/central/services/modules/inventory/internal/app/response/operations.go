package response

import "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"

type PutawaySuggestResult struct {
	Suggestions []entities.PutawaySuggestion `json:"suggestions"`
	UnresolvedItems []string                 `json:"unresolved_items"`
}

type SlottingRunResult struct {
	BusinessID  uint                       `json:"business_id"`
	WarehouseID uint                       `json:"warehouse_id"`
	Period      string                     `json:"period"`
	TotalScanned int                       `json:"total_scanned"`
	Velocities  []entities.ProductVelocity `json:"velocities"`
}

type ReplenishmentDetectResult struct {
	Created int                          `json:"created"`
	Tasks   []entities.ReplenishmentTask `json:"tasks"`
}

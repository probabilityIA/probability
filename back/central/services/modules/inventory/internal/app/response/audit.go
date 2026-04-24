package response

import "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"

type GenerateCountTaskResult struct {
	Task  entities.CycleCountTask  `json:"task"`
	Lines []entities.CycleCountLine `json:"lines"`
}

type SubmitCountLineResult struct {
	Line        entities.CycleCountLine         `json:"line"`
	Discrepancy *entities.InventoryDiscrepancy  `json:"discrepancy,omitempty"`
}

type KardexExportResult struct {
	BusinessID  uint                    `json:"business_id"`
	ProductID   string                  `json:"product_id"`
	WarehouseID uint                    `json:"warehouse_id"`
	Entries     []entities.KardexEntry  `json:"entries"`
	TotalIn     int                     `json:"total_in"`
	TotalOut    int                     `json:"total_out"`
	FinalBalance int                    `json:"final_balance"`
}

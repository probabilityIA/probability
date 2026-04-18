package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
)

// RepositoryMock implementa ports.IRepository para tests
type RepositoryMock struct {
	GetProductInventoryFn                              func(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error)
	ListWarehouseInventoryFn                           func(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error)
	GetOrCreateLevelFn                                 func(ctx context.Context, productID string, warehouseID uint, locationID *uint, businessID uint) (*entities.InventoryLevel, error)
	UpdateLevelFn                                      func(ctx context.Context, level *entities.InventoryLevel) error
	CreateMovementFn                                   func(ctx context.Context, movement *entities.StockMovement) (*entities.StockMovement, error)
	ListMovementsFn                                    func(ctx context.Context, params dtos.ListMovementsParams) ([]entities.StockMovement, int64, error)
	ListMovementTypesFn                                func(ctx context.Context, params dtos.ListStockMovementTypesParams) ([]entities.StockMovementType, int64, error)
	GetMovementTypeByIDFn                              func(ctx context.Context, id uint) (*entities.StockMovementType, error)
	GetMovementTypeIDByCodeFn                          func(ctx context.Context, code string) (uint, error)
	CreateMovementTypeFn                               func(ctx context.Context, movType *entities.StockMovementType) (*entities.StockMovementType, error)
	UpdateMovementTypeFn                               func(ctx context.Context, movType *entities.StockMovementType) error
	DeleteMovementTypeFn                               func(ctx context.Context, id uint) error
	GetProductByIDFn                                   func(ctx context.Context, productID string, businessID uint) (string, string, bool, error)
	GetProductBySKUFn                                  func(ctx context.Context, sku string, businessID uint) (string, string, bool, error)
	EnableProductTrackInventoryFn                      func(ctx context.Context, productID string) error
	UpdateProductStockQuantityFn                       func(ctx context.Context, productID string, totalQuantity int) error
	WarehouseExistsFn                                  func(ctx context.Context, warehouseID uint, businessID uint) (bool, error)
	GetProductIntegrationsFn                           func(ctx context.Context, productID string, businessID uint) ([]ports.ProductIntegrationInfo, error)
	AdjustStockTxFn                                    func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error)
	TransferStockTxFn                                  func(ctx context.Context, params dtos.TransferStockTxParams) (*dtos.TransferStockTxResult, error)
	ReserveStockTxFn                                   func(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error)
	ConfirmSaleTxFn                                    func(ctx context.Context, params dtos.ConfirmSaleTxParams) error
	ReleaseStockTxFn                                   func(ctx context.Context, params dtos.ReleaseTxParams) error
	ReturnStockTxFn                                    func(ctx context.Context, params dtos.ReturnStockTxParams) error
	GetDefaultWarehouseIDFn                            func(ctx context.Context, businessID uint) (uint, error)
	GetLocationCapacityFn                              func(ctx context.Context, locationID uint) (*ports.LocationCapacityInfo, error)
	GetProductDimensionsFn                             func(ctx context.Context, productID string, businessID uint) (*ports.ProductDimensions, error)
	GetLocationOccupiedQtyFn                           func(ctx context.Context, locationID uint) (int, error)

	CreateLotFn         func(ctx context.Context, lot *entities.InventoryLot) (*entities.InventoryLot, error)
	GetLotByIDFn        func(ctx context.Context, businessID, lotID uint) (*entities.InventoryLot, error)
	ListLotsFn          func(ctx context.Context, params dtos.ListLotsParams) ([]entities.InventoryLot, int64, error)
	UpdateLotFn         func(ctx context.Context, lot *entities.InventoryLot) (*entities.InventoryLot, error)
	DeleteLotFn         func(ctx context.Context, businessID, lotID uint) error
	LotExistsByCodeFn   func(ctx context.Context, businessID uint, productID, code string, excludeID *uint) (bool, error)

	CreateSerialFn    func(ctx context.Context, serial *entities.InventorySerial) (*entities.InventorySerial, error)
	GetSerialByIDFn   func(ctx context.Context, businessID, serialID uint) (*entities.InventorySerial, error)
	ListSerialsFn     func(ctx context.Context, params dtos.ListSerialsParams) ([]entities.InventorySerial, int64, error)
	UpdateSerialFn    func(ctx context.Context, serial *entities.InventorySerial) (*entities.InventorySerial, error)
	SerialExistsFn    func(ctx context.Context, businessID uint, productID, serial string, excludeID *uint) (bool, error)

	ListInventoryStatesFn     func(ctx context.Context) ([]entities.InventoryState, error)
	GetInventoryStateByCodeFn func(ctx context.Context, code string) (*entities.InventoryState, error)

	ListUoMsFn    func(ctx context.Context) ([]entities.UnitOfMeasure, error)
	GetUoMByCodeFn func(ctx context.Context, code string) (*entities.UnitOfMeasure, error)
	GetUoMByIDFn   func(ctx context.Context, uomID uint) (*entities.UnitOfMeasure, error)

	CreateProductUoMFn  func(ctx context.Context, pu *entities.ProductUoM) (*entities.ProductUoM, error)
	ListProductUoMsFn   func(ctx context.Context, params dtos.ListProductUoMParams) ([]entities.ProductUoM, error)
	DeleteProductUoMFn  func(ctx context.Context, businessID, id uint) error
	GetBaseProductUoMFn func(ctx context.Context, businessID uint, productID string) (*entities.ProductUoM, error)

	ChangeStateTxFn      func(ctx context.Context, params dtos.ChangeInventoryStateTxParams) (*entities.StockMovement, error)
	ListLotsForReserveFn func(ctx context.Context, productID string, warehouseID, businessID uint, strategy string) ([]entities.InventoryLot, error)

	CreatePutawayRuleFn  func(ctx context.Context, rule *entities.PutawayRule) (*entities.PutawayRule, error)
	ListPutawayRulesFn   func(ctx context.Context, params dtos.ListPutawayRulesParams) ([]entities.PutawayRule, int64, error)
	GetPutawayRuleByIDFn func(ctx context.Context, businessID, ruleID uint) (*entities.PutawayRule, error)
	UpdatePutawayRuleFn  func(ctx context.Context, rule *entities.PutawayRule) (*entities.PutawayRule, error)
	DeletePutawayRuleFn  func(ctx context.Context, businessID, ruleID uint) error
	FindApplicableRuleFn func(ctx context.Context, businessID uint, productID string) (*entities.PutawayRule, error)
	PickLocationInZoneFn func(ctx context.Context, zoneID uint) (uint, error)

	CreatePutawaySuggestionFn  func(ctx context.Context, s *entities.PutawaySuggestion) (*entities.PutawaySuggestion, error)
	GetPutawaySuggestionByIDFn func(ctx context.Context, businessID, id uint) (*entities.PutawaySuggestion, error)
	ListPutawaySuggestionsFn   func(ctx context.Context, params dtos.ListPutawaySuggestionsParams) ([]entities.PutawaySuggestion, int64, error)
	UpdatePutawaySuggestionFn  func(ctx context.Context, s *entities.PutawaySuggestion) (*entities.PutawaySuggestion, error)

	CreateReplenishmentTaskFn       func(ctx context.Context, t *entities.ReplenishmentTask) (*entities.ReplenishmentTask, error)
	GetReplenishmentTaskByIDFn      func(ctx context.Context, businessID, id uint) (*entities.ReplenishmentTask, error)
	ListReplenishmentTasksFn        func(ctx context.Context, params dtos.ListReplenishmentTasksParams) ([]entities.ReplenishmentTask, int64, error)
	UpdateReplenishmentTaskFn       func(ctx context.Context, t *entities.ReplenishmentTask) (*entities.ReplenishmentTask, error)
	DetectReplenishmentCandidatesFn func(ctx context.Context, businessID uint) ([]entities.ReplenishmentTask, error)

	CreateCrossDockLinkFn  func(ctx context.Context, l *entities.CrossDockLink) (*entities.CrossDockLink, error)
	GetCrossDockLinkByIDFn func(ctx context.Context, businessID, id uint) (*entities.CrossDockLink, error)
	ListCrossDockLinksFn   func(ctx context.Context, params dtos.ListCrossDockLinksParams) ([]entities.CrossDockLink, int64, error)
	UpdateCrossDockLinkFn  func(ctx context.Context, l *entities.CrossDockLink) (*entities.CrossDockLink, error)

	ComputeVelocitiesFn func(ctx context.Context, businessID, warehouseID uint, period string) error
	ListVelocitiesFn    func(ctx context.Context, params dtos.ListVelocityParams) ([]entities.ProductVelocity, error)

	CreateCountPlanFn  func(ctx context.Context, p *entities.CycleCountPlan) (*entities.CycleCountPlan, error)
	GetCountPlanByIDFn func(ctx context.Context, businessID, id uint) (*entities.CycleCountPlan, error)
	ListCountPlansFn   func(ctx context.Context, params dtos.ListCycleCountPlansParams) ([]entities.CycleCountPlan, int64, error)
	UpdateCountPlanFn  func(ctx context.Context, p *entities.CycleCountPlan) (*entities.CycleCountPlan, error)
	DeleteCountPlanFn  func(ctx context.Context, businessID, id uint) error

	CreateCountTaskFn           func(ctx context.Context, t *entities.CycleCountTask) (*entities.CycleCountTask, error)
	GetCountTaskByIDFn          func(ctx context.Context, businessID, id uint) (*entities.CycleCountTask, error)
	ListCountTasksFn            func(ctx context.Context, params dtos.ListCycleCountTasksParams) ([]entities.CycleCountTask, int64, error)
	UpdateCountTaskFn           func(ctx context.Context, t *entities.CycleCountTask) (*entities.CycleCountTask, error)
	GenerateCountLinesForTaskFn func(ctx context.Context, task *entities.CycleCountTask, strategy string) ([]entities.CycleCountLine, error)

	CreateCountLineFn  func(ctx context.Context, line *entities.CycleCountLine) (*entities.CycleCountLine, error)
	GetCountLineByIDFn func(ctx context.Context, businessID, id uint) (*entities.CycleCountLine, error)
	ListCountLinesFn   func(ctx context.Context, params dtos.ListCycleCountLinesParams) ([]entities.CycleCountLine, int64, error)
	UpdateCountLineFn  func(ctx context.Context, line *entities.CycleCountLine) (*entities.CycleCountLine, error)

	CreateDiscrepancyFn    func(ctx context.Context, d *entities.InventoryDiscrepancy) (*entities.InventoryDiscrepancy, error)
	GetDiscrepancyByIDFn   func(ctx context.Context, businessID, id uint) (*entities.InventoryDiscrepancy, error)
	ListDiscrepanciesFn    func(ctx context.Context, params dtos.ListDiscrepanciesParams) ([]entities.InventoryDiscrepancy, int64, error)
	UpdateDiscrepancyFn    func(ctx context.Context, d *entities.InventoryDiscrepancy) (*entities.InventoryDiscrepancy, error)
	ApproveDiscrepancyTxFn func(ctx context.Context, params dtos.ApproveDiscrepancyTxParams) (*entities.InventoryDiscrepancy, error)

	GetKardexFn func(ctx context.Context, params dtos.KardexQueryParams) ([]entities.KardexEntry, error)

	CreateLPNFn       func(ctx context.Context, lpn *entities.LicensePlate) (*entities.LicensePlate, error)
	GetLPNByIDFn      func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error)
	GetLPNByCodeFn    func(ctx context.Context, businessID uint, code string) (*entities.LicensePlate, error)
	ListLPNsFn        func(ctx context.Context, params dtos.ListLPNParams) ([]entities.LicensePlate, int64, error)
	UpdateLPNFn       func(ctx context.Context, lpn *entities.LicensePlate) (*entities.LicensePlate, error)
	DeleteLPNFn       func(ctx context.Context, businessID, id uint) error
	LPNExistsByCodeFn func(ctx context.Context, businessID uint, code string, excludeID *uint) (bool, error)
	AddLPNLineFn      func(ctx context.Context, line *entities.LicensePlateLine) (*entities.LicensePlateLine, error)
	ListLPNLinesFn    func(ctx context.Context, lpnID uint) ([]entities.LicensePlateLine, error)
	DissolveLPNFn     func(ctx context.Context, businessID, id uint) error

	RecordScanEventFn func(ctx context.Context, event *entities.ScanEvent) (*entities.ScanEvent, error)
	ResolveScanCodeFn func(ctx context.Context, businessID uint, code string) (*entities.ScanResolution, error)

	CreateSyncLogFn       func(ctx context.Context, log *entities.InventorySyncLog) (*entities.InventorySyncLog, error)
	GetSyncLogByHashFn    func(ctx context.Context, businessID uint, direction, hash string) (*entities.InventorySyncLog, error)
	UpdateSyncLogStatusFn func(ctx context.Context, id uint, status, errorMsg string) error
	ListSyncLogsFn        func(ctx context.Context, params dtos.ListSyncLogsParams) ([]entities.InventorySyncLog, int64, error)
}

func (m *RepositoryMock) GetLocationCapacity(ctx context.Context, locationID uint) (*ports.LocationCapacityInfo, error) {
	if m.GetLocationCapacityFn != nil {
		return m.GetLocationCapacityFn(ctx, locationID)
	}
	return &ports.LocationCapacityInfo{ID: locationID}, nil
}

func (m *RepositoryMock) GetProductDimensions(ctx context.Context, productID string, businessID uint) (*ports.ProductDimensions, error) {
	if m.GetProductDimensionsFn != nil {
		return m.GetProductDimensionsFn(ctx, productID, businessID)
	}
	return &ports.ProductDimensions{ID: productID}, nil
}

func (m *RepositoryMock) GetLocationOccupiedQty(ctx context.Context, locationID uint) (int, error) {
	if m.GetLocationOccupiedQtyFn != nil {
		return m.GetLocationOccupiedQtyFn(ctx, locationID)
	}
	return 0, nil
}

func (m *RepositoryMock) GetProductInventory(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
	if m.GetProductInventoryFn != nil {
		return m.GetProductInventoryFn(ctx, params)
	}
	return nil, nil
}

func (m *RepositoryMock) ListWarehouseInventory(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error) {
	if m.ListWarehouseInventoryFn != nil {
		return m.ListWarehouseInventoryFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) GetOrCreateLevel(ctx context.Context, productID string, warehouseID uint, locationID *uint, businessID uint) (*entities.InventoryLevel, error) {
	if m.GetOrCreateLevelFn != nil {
		return m.GetOrCreateLevelFn(ctx, productID, warehouseID, locationID, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) UpdateLevel(ctx context.Context, level *entities.InventoryLevel) error {
	if m.UpdateLevelFn != nil {
		return m.UpdateLevelFn(ctx, level)
	}
	return nil
}

func (m *RepositoryMock) CreateMovement(ctx context.Context, movement *entities.StockMovement) (*entities.StockMovement, error) {
	if m.CreateMovementFn != nil {
		return m.CreateMovementFn(ctx, movement)
	}
	return movement, nil
}

func (m *RepositoryMock) ListMovements(ctx context.Context, params dtos.ListMovementsParams) ([]entities.StockMovement, int64, error) {
	if m.ListMovementsFn != nil {
		return m.ListMovementsFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) ListMovementTypes(ctx context.Context, params dtos.ListStockMovementTypesParams) ([]entities.StockMovementType, int64, error) {
	if m.ListMovementTypesFn != nil {
		return m.ListMovementTypesFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) GetMovementTypeByID(ctx context.Context, id uint) (*entities.StockMovementType, error) {
	if m.GetMovementTypeByIDFn != nil {
		return m.GetMovementTypeByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *RepositoryMock) GetMovementTypeIDByCode(ctx context.Context, code string) (uint, error) {
	if m.GetMovementTypeIDByCodeFn != nil {
		return m.GetMovementTypeIDByCodeFn(ctx, code)
	}
	return 1, nil
}

func (m *RepositoryMock) CreateMovementType(ctx context.Context, movType *entities.StockMovementType) (*entities.StockMovementType, error) {
	if m.CreateMovementTypeFn != nil {
		return m.CreateMovementTypeFn(ctx, movType)
	}
	return movType, nil
}

func (m *RepositoryMock) UpdateMovementType(ctx context.Context, movType *entities.StockMovementType) error {
	if m.UpdateMovementTypeFn != nil {
		return m.UpdateMovementTypeFn(ctx, movType)
	}
	return nil
}

func (m *RepositoryMock) DeleteMovementType(ctx context.Context, id uint) error {
	if m.DeleteMovementTypeFn != nil {
		return m.DeleteMovementTypeFn(ctx, id)
	}
	return nil
}

func (m *RepositoryMock) GetProductByID(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
	if m.GetProductByIDFn != nil {
		return m.GetProductByIDFn(ctx, productID, businessID)
	}
	return "Producto Test", "SKU-001", true, nil
}

func (m *RepositoryMock) GetProductBySKU(ctx context.Context, sku string, businessID uint) (string, string, bool, error) {
	if m.GetProductBySKUFn != nil {
		return m.GetProductBySKUFn(ctx, sku, businessID)
	}
	return "product-uuid", "Producto Test", true, nil
}

func (m *RepositoryMock) EnableProductTrackInventory(ctx context.Context, productID string) error {
	if m.EnableProductTrackInventoryFn != nil {
		return m.EnableProductTrackInventoryFn(ctx, productID)
	}
	return nil
}

func (m *RepositoryMock) UpdateProductStockQuantity(ctx context.Context, productID string, totalQuantity int) error {
	if m.UpdateProductStockQuantityFn != nil {
		return m.UpdateProductStockQuantityFn(ctx, productID, totalQuantity)
	}
	return nil
}

func (m *RepositoryMock) WarehouseExists(ctx context.Context, warehouseID uint, businessID uint) (bool, error) {
	if m.WarehouseExistsFn != nil {
		return m.WarehouseExistsFn(ctx, warehouseID, businessID)
	}
	return true, nil
}

func (m *RepositoryMock) GetProductIntegrations(ctx context.Context, productID string, businessID uint) ([]ports.ProductIntegrationInfo, error) {
	if m.GetProductIntegrationsFn != nil {
		return m.GetProductIntegrationsFn(ctx, productID, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) AdjustStockTx(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
	if m.AdjustStockTxFn != nil {
		return m.AdjustStockTxFn(ctx, params)
	}
	return &dtos.AdjustStockTxResult{
		Movement:    &entities.StockMovement{ID: 1},
		NewQuantity: params.Quantity,
	}, nil
}

func (m *RepositoryMock) TransferStockTx(ctx context.Context, params dtos.TransferStockTxParams) (*dtos.TransferStockTxResult, error) {
	if m.TransferStockTxFn != nil {
		return m.TransferStockTxFn(ctx, params)
	}
	return &dtos.TransferStockTxResult{
		FromNewQty: 0,
		ToNewQty:   params.Quantity,
	}, nil
}

func (m *RepositoryMock) ReserveStockTx(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error) {
	if m.ReserveStockTxFn != nil {
		return m.ReserveStockTxFn(ctx, params)
	}
	return &dtos.ReserveStockTxResult{
		Reserved:   params.Quantity,
		Sufficient: true,
	}, nil
}

func (m *RepositoryMock) ConfirmSaleTx(ctx context.Context, params dtos.ConfirmSaleTxParams) error {
	if m.ConfirmSaleTxFn != nil {
		return m.ConfirmSaleTxFn(ctx, params)
	}
	return nil
}

func (m *RepositoryMock) ReleaseStockTx(ctx context.Context, params dtos.ReleaseTxParams) error {
	if m.ReleaseStockTxFn != nil {
		return m.ReleaseStockTxFn(ctx, params)
	}
	return nil
}

func (m *RepositoryMock) ReturnStockTx(ctx context.Context, params dtos.ReturnStockTxParams) error {
	if m.ReturnStockTxFn != nil {
		return m.ReturnStockTxFn(ctx, params)
	}
	return nil
}

func (m *RepositoryMock) GetDefaultWarehouseID(ctx context.Context, businessID uint) (uint, error) {
	if m.GetDefaultWarehouseIDFn != nil {
		return m.GetDefaultWarehouseIDFn(ctx, businessID)
	}
	return 1, nil
}

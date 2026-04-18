package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

type UseCaseMock struct {
	GetProductInventoryFn    func(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error)
	ListWarehouseInventoryFn func(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error)
	AdjustStockFn            func(ctx context.Context, dto request.AdjustStockDTO) (*entities.StockMovement, error)
	TransferStockFn          func(ctx context.Context, dto request.TransferStockDTO) error
	ListMovementsFn          func(ctx context.Context, params dtos.ListMovementsParams) ([]entities.StockMovement, int64, error)
	ListMovementTypesFn      func(ctx context.Context, params dtos.ListStockMovementTypesParams) ([]entities.StockMovementType, int64, error)
	CreateMovementTypeFn     func(ctx context.Context, dto request.CreateStockMovementTypeDTO) (*entities.StockMovementType, error)
	UpdateMovementTypeFn     func(ctx context.Context, dto request.UpdateStockMovementTypeDTO) (*entities.StockMovementType, error)
	DeleteMovementTypeFn     func(ctx context.Context, id uint) error
	BulkLoadInventoryFn      func(ctx context.Context, dto request.BulkLoadDTO) (*response.BulkLoadResult, error)
	ReserveStockForOrderFn   func(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error)
	ConfirmSaleForOrderFn    func(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error)
	ReleaseStockForOrderFn   func(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error)
	ReturnStockForOrderFn    func(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error)
	ValidateCubingFn         func(ctx context.Context, dto request.ValidateCubingDTO) (*response.CubingCheckResult, error)

	CreateLotFn func(ctx context.Context, dto request.CreateLotDTO) (*entities.InventoryLot, error)
	GetLotFn    func(ctx context.Context, businessID, lotID uint) (*entities.InventoryLot, error)
	ListLotsFn  func(ctx context.Context, params dtos.ListLotsParams) ([]entities.InventoryLot, int64, error)
	UpdateLotFn func(ctx context.Context, dto request.UpdateLotDTO) (*entities.InventoryLot, error)
	DeleteLotFn func(ctx context.Context, businessID, lotID uint) error

	CreateSerialFn func(ctx context.Context, dto request.CreateSerialDTO) (*entities.InventorySerial, error)
	GetSerialFn    func(ctx context.Context, businessID, serialID uint) (*entities.InventorySerial, error)
	ListSerialsFn  func(ctx context.Context, params dtos.ListSerialsParams) ([]entities.InventorySerial, int64, error)
	UpdateSerialFn func(ctx context.Context, dto request.UpdateSerialDTO) (*entities.InventorySerial, error)

	ListInventoryStatesFn  func(ctx context.Context) ([]entities.InventoryState, error)
	ChangeInventoryStateFn func(ctx context.Context, dto request.ChangeInventoryStateDTO) (*entities.StockMovement, error)

	ListUoMsFn         func(ctx context.Context) ([]entities.UnitOfMeasure, error)
	ListProductUoMsFn  func(ctx context.Context, businessID uint, productID string) ([]entities.ProductUoM, error)
	CreateProductUoMFn func(ctx context.Context, dto request.CreateProductUoMDTO) (*entities.ProductUoM, error)
	DeleteProductUoMFn func(ctx context.Context, businessID, id uint) error
	ConvertUoMFn       func(ctx context.Context, dto request.ConvertUoMDTO) (*response.ConvertUoMResult, error)

	CreatePutawayRuleFn      func(ctx context.Context, dto request.CreatePutawayRuleDTO) (*entities.PutawayRule, error)
	ListPutawayRulesFn       func(ctx context.Context, params dtos.ListPutawayRulesParams) ([]entities.PutawayRule, int64, error)
	UpdatePutawayRuleFn      func(ctx context.Context, dto request.UpdatePutawayRuleDTO) (*entities.PutawayRule, error)
	DeletePutawayRuleFn      func(ctx context.Context, businessID, ruleID uint) error
	SuggestPutawayFn         func(ctx context.Context, dto request.PutawaySuggestDTO) (*response.PutawaySuggestResult, error)
	ConfirmPutawayFn         func(ctx context.Context, dto request.ConfirmPutawayDTO) (*entities.PutawaySuggestion, error)
	ListPutawaySuggestionsFn func(ctx context.Context, params dtos.ListPutawaySuggestionsParams) ([]entities.PutawaySuggestion, int64, error)

	CreateReplenishmentTaskFn  func(ctx context.Context, dto request.CreateReplenishmentTaskDTO) (*entities.ReplenishmentTask, error)
	ListReplenishmentTasksFn   func(ctx context.Context, params dtos.ListReplenishmentTasksParams) ([]entities.ReplenishmentTask, int64, error)
	AssignReplenishmentFn      func(ctx context.Context, dto request.AssignReplenishmentDTO) (*entities.ReplenishmentTask, error)
	CompleteReplenishmentFn    func(ctx context.Context, dto request.CompleteReplenishmentDTO) (*entities.ReplenishmentTask, error)
	CancelReplenishmentFn      func(ctx context.Context, businessID, taskID uint, reason string) (*entities.ReplenishmentTask, error)
	DetectReplenishmentNeedsFn func(ctx context.Context, businessID uint) (*response.ReplenishmentDetectResult, error)

	CreateCrossDockLinkFn func(ctx context.Context, dto request.CreateCrossDockLinkDTO) (*entities.CrossDockLink, error)
	ListCrossDockLinksFn  func(ctx context.Context, params dtos.ListCrossDockLinksParams) ([]entities.CrossDockLink, int64, error)
	ExecuteCrossDockFn    func(ctx context.Context, dto request.ExecuteCrossDockDTO) (*entities.CrossDockLink, error)

	RunSlottingFn    func(ctx context.Context, dto request.RunSlottingDTO) (*response.SlottingRunResult, error)
	ListVelocitiesFn func(ctx context.Context, params dtos.ListVelocityParams) ([]entities.ProductVelocity, error)

	CreateCountPlanFn func(ctx context.Context, dto request.CreateCountPlanDTO) (*entities.CycleCountPlan, error)
	GetCountPlanFn    func(ctx context.Context, businessID, id uint) (*entities.CycleCountPlan, error)
	ListCountPlansFn  func(ctx context.Context, params dtos.ListCycleCountPlansParams) ([]entities.CycleCountPlan, int64, error)
	UpdateCountPlanFn func(ctx context.Context, dto request.UpdateCountPlanDTO) (*entities.CycleCountPlan, error)
	DeleteCountPlanFn func(ctx context.Context, businessID, id uint) error

	GenerateCountTaskFn func(ctx context.Context, dto request.GenerateCountTaskDTO) (*response.GenerateCountTaskResult, error)
	ListCountTasksFn    func(ctx context.Context, params dtos.ListCycleCountTasksParams) ([]entities.CycleCountTask, int64, error)
	GetCountTaskFn      func(ctx context.Context, businessID, id uint) (*entities.CycleCountTask, error)
	StartCountTaskFn    func(ctx context.Context, dto request.StartCountTaskDTO) (*entities.CycleCountTask, error)
	FinishCountTaskFn   func(ctx context.Context, businessID, id uint) (*entities.CycleCountTask, error)

	ListCountLinesFn  func(ctx context.Context, params dtos.ListCycleCountLinesParams) ([]entities.CycleCountLine, int64, error)
	SubmitCountLineFn func(ctx context.Context, dto request.SubmitCountLineDTO) (*response.SubmitCountLineResult, error)

	ListDiscrepanciesFn  func(ctx context.Context, params dtos.ListDiscrepanciesParams) ([]entities.InventoryDiscrepancy, int64, error)
	GetDiscrepancyFn     func(ctx context.Context, businessID, id uint) (*entities.InventoryDiscrepancy, error)
	ApproveDiscrepancyFn func(ctx context.Context, dto request.ApproveDiscrepancyDTO) (*entities.InventoryDiscrepancy, error)
	RejectDiscrepancyFn  func(ctx context.Context, dto request.RejectDiscrepancyDTO) (*entities.InventoryDiscrepancy, error)

	ExportKardexFn func(ctx context.Context, dto request.KardexExportDTO) (*response.KardexExportResult, error)

	CreateLPNFn   func(ctx context.Context, dto request.CreateLPNDTO) (*entities.LicensePlate, error)
	GetLPNFn      func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error)
	ListLPNsFn    func(ctx context.Context, params dtos.ListLPNParams) ([]entities.LicensePlate, int64, error)
	UpdateLPNFn   func(ctx context.Context, dto request.UpdateLPNDTO) (*entities.LicensePlate, error)
	DeleteLPNFn   func(ctx context.Context, businessID, id uint) error
	AddToLPNFn    func(ctx context.Context, dto request.AddToLPNDTO) (*entities.LicensePlateLine, error)
	MoveLPNFn     func(ctx context.Context, dto request.MoveLPNDTO) (*entities.LicensePlate, error)
	DissolveLPNFn func(ctx context.Context, dto request.DissolveLPNDTO) error
	MergeLPNFn    func(ctx context.Context, dto request.MergeLPNDTO) (*entities.LicensePlate, error)

	ScanFn func(ctx context.Context, dto request.ScanDTO) (*response.ScanResponse, error)

	InboundSyncFn  func(ctx context.Context, dto request.InboundSyncDTO) (*response.InboundSyncResult, error)
	ListSyncLogsFn func(ctx context.Context, params dtos.ListSyncLogsParams) ([]entities.InventorySyncLog, int64, error)
}

func (m *UseCaseMock) ValidateCubing(ctx context.Context, dto request.ValidateCubingDTO) (*response.CubingCheckResult, error) {
	if m.ValidateCubingFn != nil {
		return m.ValidateCubingFn(ctx, dto)
	}
	return &response.CubingCheckResult{Fits: true}, nil
}

func (m *UseCaseMock) GetProductInventory(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
	if m.GetProductInventoryFn != nil {
		return m.GetProductInventoryFn(ctx, params)
	}
	return nil, nil
}

func (m *UseCaseMock) ListWarehouseInventory(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error) {
	if m.ListWarehouseInventoryFn != nil {
		return m.ListWarehouseInventoryFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *UseCaseMock) AdjustStock(ctx context.Context, dto request.AdjustStockDTO) (*entities.StockMovement, error) {
	if m.AdjustStockFn != nil {
		return m.AdjustStockFn(ctx, dto)
	}
	return &entities.StockMovement{ID: 1}, nil
}

func (m *UseCaseMock) TransferStock(ctx context.Context, dto request.TransferStockDTO) error {
	if m.TransferStockFn != nil {
		return m.TransferStockFn(ctx, dto)
	}
	return nil
}

func (m *UseCaseMock) ListMovements(ctx context.Context, params dtos.ListMovementsParams) ([]entities.StockMovement, int64, error) {
	if m.ListMovementsFn != nil {
		return m.ListMovementsFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *UseCaseMock) ListMovementTypes(ctx context.Context, params dtos.ListStockMovementTypesParams) ([]entities.StockMovementType, int64, error) {
	if m.ListMovementTypesFn != nil {
		return m.ListMovementTypesFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *UseCaseMock) CreateMovementType(ctx context.Context, dto request.CreateStockMovementTypeDTO) (*entities.StockMovementType, error) {
	if m.CreateMovementTypeFn != nil {
		return m.CreateMovementTypeFn(ctx, dto)
	}
	return &entities.StockMovementType{ID: 1, Code: dto.Code, Name: dto.Name}, nil
}

func (m *UseCaseMock) UpdateMovementType(ctx context.Context, dto request.UpdateStockMovementTypeDTO) (*entities.StockMovementType, error) {
	if m.UpdateMovementTypeFn != nil {
		return m.UpdateMovementTypeFn(ctx, dto)
	}
	return &entities.StockMovementType{ID: dto.ID, Name: dto.Name}, nil
}

func (m *UseCaseMock) DeleteMovementType(ctx context.Context, id uint) error {
	if m.DeleteMovementTypeFn != nil {
		return m.DeleteMovementTypeFn(ctx, id)
	}
	return nil
}

func (m *UseCaseMock) BulkLoadInventory(ctx context.Context, dto request.BulkLoadDTO) (*response.BulkLoadResult, error) {
	if m.BulkLoadInventoryFn != nil {
		return m.BulkLoadInventoryFn(ctx, dto)
	}
	return &response.BulkLoadResult{}, nil
}

func (m *UseCaseMock) ReserveStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error) {
	if m.ReserveStockForOrderFn != nil {
		return m.ReserveStockForOrderFn(ctx, orderID, businessID, warehouseID, items)
	}
	return &response.OrderStockResult{OrderID: orderID, Success: true}, nil
}

func (m *UseCaseMock) ConfirmSaleForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error) {
	if m.ConfirmSaleForOrderFn != nil {
		return m.ConfirmSaleForOrderFn(ctx, orderID, businessID, warehouseID, items)
	}
	return &response.OrderStockResult{OrderID: orderID, Success: true}, nil
}

func (m *UseCaseMock) ReleaseStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error) {
	if m.ReleaseStockForOrderFn != nil {
		return m.ReleaseStockForOrderFn(ctx, orderID, businessID, warehouseID, items)
	}
	return &response.OrderStockResult{OrderID: orderID, Success: true}, nil
}

func (m *UseCaseMock) ReturnStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error) {
	if m.ReturnStockForOrderFn != nil {
		return m.ReturnStockForOrderFn(ctx, orderID, businessID, warehouseID, items)
	}
	return &response.OrderStockResult{OrderID: orderID, Success: true}, nil
}

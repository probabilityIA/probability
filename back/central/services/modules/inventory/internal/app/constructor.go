package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IUseCase interface {
	GetProductInventory(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error)
	ListWarehouseInventory(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error)
	AdjustStock(ctx context.Context, dto request.AdjustStockDTO) (*entities.StockMovement, error)
	TransferStock(ctx context.Context, dto request.TransferStockDTO) error
	ListMovements(ctx context.Context, params dtos.ListMovementsParams) ([]entities.StockMovement, int64, error)

	ListMovementTypes(ctx context.Context, params dtos.ListStockMovementTypesParams) ([]entities.StockMovementType, int64, error)
	CreateMovementType(ctx context.Context, dto request.CreateStockMovementTypeDTO) (*entities.StockMovementType, error)
	UpdateMovementType(ctx context.Context, dto request.UpdateStockMovementTypeDTO) (*entities.StockMovementType, error)
	DeleteMovementType(ctx context.Context, id uint) error

	BulkLoadInventory(ctx context.Context, dto request.BulkLoadDTO) (*response.BulkLoadResult, error)

	ReserveStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error)
	ConfirmSaleForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error)
	ReleaseStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error)
	ReturnStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error)

	ValidateCubing(ctx context.Context, dto request.ValidateCubingDTO) (*response.CubingCheckResult, error)

	CreateLot(ctx context.Context, dto request.CreateLotDTO) (*entities.InventoryLot, error)
	GetLot(ctx context.Context, businessID, lotID uint) (*entities.InventoryLot, error)
	ListLots(ctx context.Context, params dtos.ListLotsParams) ([]entities.InventoryLot, int64, error)
	UpdateLot(ctx context.Context, dto request.UpdateLotDTO) (*entities.InventoryLot, error)
	DeleteLot(ctx context.Context, businessID, lotID uint) error

	CreateSerial(ctx context.Context, dto request.CreateSerialDTO) (*entities.InventorySerial, error)
	GetSerial(ctx context.Context, businessID, serialID uint) (*entities.InventorySerial, error)
	ListSerials(ctx context.Context, params dtos.ListSerialsParams) ([]entities.InventorySerial, int64, error)
	UpdateSerial(ctx context.Context, dto request.UpdateSerialDTO) (*entities.InventorySerial, error)

	ListInventoryStates(ctx context.Context) ([]entities.InventoryState, error)
	ChangeInventoryState(ctx context.Context, dto request.ChangeInventoryStateDTO) (*entities.StockMovement, error)

	ListUoMs(ctx context.Context) ([]entities.UnitOfMeasure, error)
	ListProductUoMs(ctx context.Context, businessID uint, productID string) ([]entities.ProductUoM, error)
	CreateProductUoM(ctx context.Context, dto request.CreateProductUoMDTO) (*entities.ProductUoM, error)
	DeleteProductUoM(ctx context.Context, businessID, id uint) error
	ConvertUoM(ctx context.Context, dto request.ConvertUoMDTO) (*response.ConvertUoMResult, error)

	CreatePutawayRule(ctx context.Context, dto request.CreatePutawayRuleDTO) (*entities.PutawayRule, error)
	ListPutawayRules(ctx context.Context, params dtos.ListPutawayRulesParams) ([]entities.PutawayRule, int64, error)
	UpdatePutawayRule(ctx context.Context, dto request.UpdatePutawayRuleDTO) (*entities.PutawayRule, error)
	DeletePutawayRule(ctx context.Context, businessID, ruleID uint) error
	SuggestPutaway(ctx context.Context, dto request.PutawaySuggestDTO) (*response.PutawaySuggestResult, error)
	ConfirmPutaway(ctx context.Context, dto request.ConfirmPutawayDTO) (*entities.PutawaySuggestion, error)
	ListPutawaySuggestions(ctx context.Context, params dtos.ListPutawaySuggestionsParams) ([]entities.PutawaySuggestion, int64, error)

	CreateReplenishmentTask(ctx context.Context, dto request.CreateReplenishmentTaskDTO) (*entities.ReplenishmentTask, error)
	ListReplenishmentTasks(ctx context.Context, params dtos.ListReplenishmentTasksParams) ([]entities.ReplenishmentTask, int64, error)
	AssignReplenishment(ctx context.Context, dto request.AssignReplenishmentDTO) (*entities.ReplenishmentTask, error)
	CompleteReplenishment(ctx context.Context, dto request.CompleteReplenishmentDTO) (*entities.ReplenishmentTask, error)
	CancelReplenishment(ctx context.Context, businessID, taskID uint, reason string) (*entities.ReplenishmentTask, error)
	DetectReplenishmentNeeds(ctx context.Context, businessID uint) (*response.ReplenishmentDetectResult, error)

	CreateCrossDockLink(ctx context.Context, dto request.CreateCrossDockLinkDTO) (*entities.CrossDockLink, error)
	ListCrossDockLinks(ctx context.Context, params dtos.ListCrossDockLinksParams) ([]entities.CrossDockLink, int64, error)
	ExecuteCrossDock(ctx context.Context, dto request.ExecuteCrossDockDTO) (*entities.CrossDockLink, error)

	RunSlotting(ctx context.Context, dto request.RunSlottingDTO) (*response.SlottingRunResult, error)
	ListVelocities(ctx context.Context, params dtos.ListVelocityParams) ([]entities.ProductVelocity, error)

	CreateCountPlan(ctx context.Context, dto request.CreateCountPlanDTO) (*entities.CycleCountPlan, error)
	GetCountPlan(ctx context.Context, businessID, id uint) (*entities.CycleCountPlan, error)
	ListCountPlans(ctx context.Context, params dtos.ListCycleCountPlansParams) ([]entities.CycleCountPlan, int64, error)
	UpdateCountPlan(ctx context.Context, dto request.UpdateCountPlanDTO) (*entities.CycleCountPlan, error)
	DeleteCountPlan(ctx context.Context, businessID, id uint) error

	GenerateCountTask(ctx context.Context, dto request.GenerateCountTaskDTO) (*response.GenerateCountTaskResult, error)
	ListCountTasks(ctx context.Context, params dtos.ListCycleCountTasksParams) ([]entities.CycleCountTask, int64, error)
	GetCountTask(ctx context.Context, businessID, id uint) (*entities.CycleCountTask, error)
	StartCountTask(ctx context.Context, dto request.StartCountTaskDTO) (*entities.CycleCountTask, error)
	FinishCountTask(ctx context.Context, businessID, id uint) (*entities.CycleCountTask, error)

	ListCountLines(ctx context.Context, params dtos.ListCycleCountLinesParams) ([]entities.CycleCountLine, int64, error)
	SubmitCountLine(ctx context.Context, dto request.SubmitCountLineDTO) (*response.SubmitCountLineResult, error)

	ListDiscrepancies(ctx context.Context, params dtos.ListDiscrepanciesParams) ([]entities.InventoryDiscrepancy, int64, error)
	GetDiscrepancy(ctx context.Context, businessID, id uint) (*entities.InventoryDiscrepancy, error)
	ApproveDiscrepancy(ctx context.Context, dto request.ApproveDiscrepancyDTO) (*entities.InventoryDiscrepancy, error)
	RejectDiscrepancy(ctx context.Context, dto request.RejectDiscrepancyDTO) (*entities.InventoryDiscrepancy, error)

	ExportKardex(ctx context.Context, dto request.KardexExportDTO) (*response.KardexExportResult, error)

	CreateLPN(ctx context.Context, dto request.CreateLPNDTO) (*entities.LicensePlate, error)
	GetLPN(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error)
	ListLPNs(ctx context.Context, params dtos.ListLPNParams) ([]entities.LicensePlate, int64, error)
	UpdateLPN(ctx context.Context, dto request.UpdateLPNDTO) (*entities.LicensePlate, error)
	DeleteLPN(ctx context.Context, businessID, id uint) error
	AddToLPN(ctx context.Context, dto request.AddToLPNDTO) (*entities.LicensePlateLine, error)
	MoveLPN(ctx context.Context, dto request.MoveLPNDTO) (*entities.LicensePlate, error)
	DissolveLPN(ctx context.Context, dto request.DissolveLPNDTO) error
	MergeLPN(ctx context.Context, dto request.MergeLPNDTO) (*entities.LicensePlate, error)

	Scan(ctx context.Context, dto request.ScanDTO) (*response.ScanResponse, error)

	InboundSync(ctx context.Context, dto request.InboundSyncDTO) (*response.InboundSyncResult, error)
	ListSyncLogs(ctx context.Context, params dtos.ListSyncLogsParams) ([]entities.InventorySyncLog, int64, error)
}

type useCase struct {
	repo           ports.IRepository
	publisher      ports.ISyncPublisher
	eventPublisher ports.IInventoryEventPublisher
	log            log.ILogger
}

func New(repo ports.IRepository, publisher ports.ISyncPublisher, eventPublisher ports.IInventoryEventPublisher, logger log.ILogger) IUseCase {
	return &useCase{
		repo:           repo,
		publisher:      publisher,
		eventPublisher: eventPublisher,
		log:            logger,
	}
}

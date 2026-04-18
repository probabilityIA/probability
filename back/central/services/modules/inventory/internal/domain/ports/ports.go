package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

// IRepository define los métodos del repositorio del módulo inventory
type IRepository interface {
	// Inventory Levels
	GetProductInventory(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error)
	ListWarehouseInventory(ctx context.Context, params dtos.ListWarehouseInventoryParams) ([]entities.InventoryLevel, int64, error)
	GetOrCreateLevel(ctx context.Context, productID string, warehouseID uint, locationID *uint, businessID uint) (*entities.InventoryLevel, error)
	UpdateLevel(ctx context.Context, level *entities.InventoryLevel) error

	// Stock Movements
	CreateMovement(ctx context.Context, movement *entities.StockMovement) (*entities.StockMovement, error)
	ListMovements(ctx context.Context, params dtos.ListMovementsParams) ([]entities.StockMovement, int64, error)

	// Stock Movement Types
	ListMovementTypes(ctx context.Context, params dtos.ListStockMovementTypesParams) ([]entities.StockMovementType, int64, error)
	GetMovementTypeByID(ctx context.Context, id uint) (*entities.StockMovementType, error)
	GetMovementTypeIDByCode(ctx context.Context, code string) (uint, error)
	CreateMovementType(ctx context.Context, movType *entities.StockMovementType) (*entities.StockMovementType, error)
	UpdateMovementType(ctx context.Context, movType *entities.StockMovementType) error
	DeleteMovementType(ctx context.Context, id uint) error

	// Product queries (replicadas localmente)
	GetProductByID(ctx context.Context, productID string, businessID uint) (productName string, productSKU string, trackInventory bool, err error)
	GetProductBySKU(ctx context.Context, sku string, businessID uint) (productID string, name string, trackInventory bool, err error)
	EnableProductTrackInventory(ctx context.Context, productID string) error
	UpdateProductStockQuantity(ctx context.Context, productID string, totalQuantity int) error

	// Warehouse validation (replicada localmente)
	WarehouseExists(ctx context.Context, warehouseID uint, businessID uint) (bool, error)

	// Product-Integration queries (replicada para sync)
	GetProductIntegrations(ctx context.Context, productID string, businessID uint) ([]ProductIntegrationInfo, error)

	// Operaciones transaccionales (SELECT FOR UPDATE + commit atómico)
	AdjustStockTx(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error)
	TransferStockTx(ctx context.Context, params dtos.TransferStockTxParams) (*dtos.TransferStockTxResult, error)

	// Operaciones transaccionales de órdenes
	ReserveStockTx(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error)
	ConfirmSaleTx(ctx context.Context, params dtos.ConfirmSaleTxParams) error
	ReleaseStockTx(ctx context.Context, params dtos.ReleaseTxParams) error
	ReturnStockTx(ctx context.Context, params dtos.ReturnStockTxParams) error

	// Warehouse queries
	GetDefaultWarehouseID(ctx context.Context, businessID uint) (uint, error)

	GetLocationCapacity(ctx context.Context, locationID uint) (*LocationCapacityInfo, error)
	GetProductDimensions(ctx context.Context, productID string, businessID uint) (*ProductDimensions, error)
	GetLocationOccupiedQty(ctx context.Context, locationID uint) (int, error)

	CreateLot(ctx context.Context, lot *entities.InventoryLot) (*entities.InventoryLot, error)
	GetLotByID(ctx context.Context, businessID, lotID uint) (*entities.InventoryLot, error)
	ListLots(ctx context.Context, params dtos.ListLotsParams) ([]entities.InventoryLot, int64, error)
	UpdateLot(ctx context.Context, lot *entities.InventoryLot) (*entities.InventoryLot, error)
	DeleteLot(ctx context.Context, businessID, lotID uint) error
	LotExistsByCode(ctx context.Context, businessID uint, productID, code string, excludeID *uint) (bool, error)

	CreateSerial(ctx context.Context, serial *entities.InventorySerial) (*entities.InventorySerial, error)
	GetSerialByID(ctx context.Context, businessID, serialID uint) (*entities.InventorySerial, error)
	ListSerials(ctx context.Context, params dtos.ListSerialsParams) ([]entities.InventorySerial, int64, error)
	UpdateSerial(ctx context.Context, serial *entities.InventorySerial) (*entities.InventorySerial, error)
	SerialExists(ctx context.Context, businessID uint, productID, serial string, excludeID *uint) (bool, error)

	ListInventoryStates(ctx context.Context) ([]entities.InventoryState, error)
	GetInventoryStateByCode(ctx context.Context, code string) (*entities.InventoryState, error)

	ListUoMs(ctx context.Context) ([]entities.UnitOfMeasure, error)
	GetUoMByCode(ctx context.Context, code string) (*entities.UnitOfMeasure, error)
	GetUoMByID(ctx context.Context, uomID uint) (*entities.UnitOfMeasure, error)

	CreateProductUoM(ctx context.Context, pu *entities.ProductUoM) (*entities.ProductUoM, error)
	ListProductUoMs(ctx context.Context, params dtos.ListProductUoMParams) ([]entities.ProductUoM, error)
	DeleteProductUoM(ctx context.Context, businessID, id uint) error
	GetBaseProductUoM(ctx context.Context, businessID uint, productID string) (*entities.ProductUoM, error)

	ChangeStateTx(ctx context.Context, params dtos.ChangeInventoryStateTxParams) (*entities.StockMovement, error)
	ListLotsForReserve(ctx context.Context, productID string, warehouseID, businessID uint, strategy string) ([]entities.InventoryLot, error)

	CreatePutawayRule(ctx context.Context, rule *entities.PutawayRule) (*entities.PutawayRule, error)
	ListPutawayRules(ctx context.Context, params dtos.ListPutawayRulesParams) ([]entities.PutawayRule, int64, error)
	GetPutawayRuleByID(ctx context.Context, businessID, ruleID uint) (*entities.PutawayRule, error)
	UpdatePutawayRule(ctx context.Context, rule *entities.PutawayRule) (*entities.PutawayRule, error)
	DeletePutawayRule(ctx context.Context, businessID, ruleID uint) error
	FindApplicableRule(ctx context.Context, businessID uint, productID string) (*entities.PutawayRule, error)
	PickLocationInZone(ctx context.Context, zoneID uint) (uint, error)

	CreatePutawaySuggestion(ctx context.Context, s *entities.PutawaySuggestion) (*entities.PutawaySuggestion, error)
	GetPutawaySuggestionByID(ctx context.Context, businessID, id uint) (*entities.PutawaySuggestion, error)
	ListPutawaySuggestions(ctx context.Context, params dtos.ListPutawaySuggestionsParams) ([]entities.PutawaySuggestion, int64, error)
	UpdatePutawaySuggestion(ctx context.Context, s *entities.PutawaySuggestion) (*entities.PutawaySuggestion, error)

	CreateReplenishmentTask(ctx context.Context, t *entities.ReplenishmentTask) (*entities.ReplenishmentTask, error)
	GetReplenishmentTaskByID(ctx context.Context, businessID, id uint) (*entities.ReplenishmentTask, error)
	ListReplenishmentTasks(ctx context.Context, params dtos.ListReplenishmentTasksParams) ([]entities.ReplenishmentTask, int64, error)
	UpdateReplenishmentTask(ctx context.Context, t *entities.ReplenishmentTask) (*entities.ReplenishmentTask, error)
	DetectReplenishmentCandidates(ctx context.Context, businessID uint) ([]entities.ReplenishmentTask, error)

	CreateCrossDockLink(ctx context.Context, l *entities.CrossDockLink) (*entities.CrossDockLink, error)
	GetCrossDockLinkByID(ctx context.Context, businessID, id uint) (*entities.CrossDockLink, error)
	ListCrossDockLinks(ctx context.Context, params dtos.ListCrossDockLinksParams) ([]entities.CrossDockLink, int64, error)
	UpdateCrossDockLink(ctx context.Context, l *entities.CrossDockLink) (*entities.CrossDockLink, error)

	ComputeVelocities(ctx context.Context, businessID, warehouseID uint, period string) error
	ListVelocities(ctx context.Context, params dtos.ListVelocityParams) ([]entities.ProductVelocity, error)

	CreateCountPlan(ctx context.Context, p *entities.CycleCountPlan) (*entities.CycleCountPlan, error)
	GetCountPlanByID(ctx context.Context, businessID, id uint) (*entities.CycleCountPlan, error)
	ListCountPlans(ctx context.Context, params dtos.ListCycleCountPlansParams) ([]entities.CycleCountPlan, int64, error)
	UpdateCountPlan(ctx context.Context, p *entities.CycleCountPlan) (*entities.CycleCountPlan, error)
	DeleteCountPlan(ctx context.Context, businessID, id uint) error

	CreateCountTask(ctx context.Context, t *entities.CycleCountTask) (*entities.CycleCountTask, error)
	GetCountTaskByID(ctx context.Context, businessID, id uint) (*entities.CycleCountTask, error)
	ListCountTasks(ctx context.Context, params dtos.ListCycleCountTasksParams) ([]entities.CycleCountTask, int64, error)
	UpdateCountTask(ctx context.Context, t *entities.CycleCountTask) (*entities.CycleCountTask, error)
	GenerateCountLinesForTask(ctx context.Context, task *entities.CycleCountTask, strategy string) ([]entities.CycleCountLine, error)

	CreateCountLine(ctx context.Context, line *entities.CycleCountLine) (*entities.CycleCountLine, error)
	GetCountLineByID(ctx context.Context, businessID, id uint) (*entities.CycleCountLine, error)
	ListCountLines(ctx context.Context, params dtos.ListCycleCountLinesParams) ([]entities.CycleCountLine, int64, error)
	UpdateCountLine(ctx context.Context, line *entities.CycleCountLine) (*entities.CycleCountLine, error)

	CreateDiscrepancy(ctx context.Context, d *entities.InventoryDiscrepancy) (*entities.InventoryDiscrepancy, error)
	GetDiscrepancyByID(ctx context.Context, businessID, id uint) (*entities.InventoryDiscrepancy, error)
	ListDiscrepancies(ctx context.Context, params dtos.ListDiscrepanciesParams) ([]entities.InventoryDiscrepancy, int64, error)
	UpdateDiscrepancy(ctx context.Context, d *entities.InventoryDiscrepancy) (*entities.InventoryDiscrepancy, error)
	ApproveDiscrepancyTx(ctx context.Context, params dtos.ApproveDiscrepancyTxParams) (*entities.InventoryDiscrepancy, error)

	GetKardex(ctx context.Context, params dtos.KardexQueryParams) ([]entities.KardexEntry, error)

	CreateLPN(ctx context.Context, lpn *entities.LicensePlate) (*entities.LicensePlate, error)
	GetLPNByID(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error)
	GetLPNByCode(ctx context.Context, businessID uint, code string) (*entities.LicensePlate, error)
	ListLPNs(ctx context.Context, params dtos.ListLPNParams) ([]entities.LicensePlate, int64, error)
	UpdateLPN(ctx context.Context, lpn *entities.LicensePlate) (*entities.LicensePlate, error)
	DeleteLPN(ctx context.Context, businessID, id uint) error
	LPNExistsByCode(ctx context.Context, businessID uint, code string, excludeID *uint) (bool, error)

	AddLPNLine(ctx context.Context, line *entities.LicensePlateLine) (*entities.LicensePlateLine, error)
	ListLPNLines(ctx context.Context, lpnID uint) ([]entities.LicensePlateLine, error)
	DissolveLPN(ctx context.Context, businessID, id uint) error

	RecordScanEvent(ctx context.Context, event *entities.ScanEvent) (*entities.ScanEvent, error)
	ResolveScanCode(ctx context.Context, businessID uint, code string) (*entities.ScanResolution, error)

	CreateSyncLog(ctx context.Context, log *entities.InventorySyncLog) (*entities.InventorySyncLog, error)
	GetSyncLogByHash(ctx context.Context, businessID uint, direction, hash string) (*entities.InventorySyncLog, error)
	UpdateSyncLogStatus(ctx context.Context, id uint, status, errorMsg string) error
	ListSyncLogs(ctx context.Context, params dtos.ListSyncLogsParams) ([]entities.InventorySyncLog, int64, error)
}

type LocationCapacityInfo struct {
	ID           uint
	MaxWeightKg  *float64
	MaxVolumeCm3 *float64
}

type ProductDimensions struct {
	ID      string
	Weight  float64
	WeightU string
	Length  float64
	Width   float64
	Height  float64
	DimU    string
}

// ProductIntegrationInfo datos de una integración vinculada a un producto
type ProductIntegrationInfo struct {
	IntegrationID       uint
	ExternalProductID   string
	IntegrationTypeCode string
}

// ISyncPublisher publica mensajes de sync a RabbitMQ
type ISyncPublisher interface {
	PublishInventorySync(ctx context.Context, msg InventorySyncMessage) error
}

// IInventoryEventPublisher publica eventos de inventario a Redis SSE
type IInventoryEventPublisher interface {
	PublishInventoryEvent(ctx context.Context, event InventoryEvent) error
}

// InventoryEvent evento de inventario para SSE
type InventoryEvent struct {
	EventType   string                 `json:"event_type"`
	OrderID     string                 `json:"order_id"`
	BusinessID  uint                   `json:"business_id"`
	WarehouseID uint                   `json:"warehouse_id"`
	Timestamp   string                 `json:"timestamp"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// InventorySyncMessage mensaje que se publica a RabbitMQ para sincronizar inventario
type InventorySyncMessage struct {
	ProductID         string `json:"product_id"`
	ExternalProductID string `json:"external_product_id"`
	IntegrationID     uint   `json:"integration_id"`
	BusinessID        uint   `json:"business_id"`
	NewQuantity       int    `json:"new_quantity"`
	WarehouseID       uint   `json:"warehouse_id"`
	Source            string `json:"source"` // manual_adjustment, transfer, sync
	Timestamp         string `json:"timestamp"`
}

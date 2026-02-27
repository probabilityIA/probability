package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

// InventoryLevelResponse respuesta de nivel de inventario
type InventoryLevelResponse struct {
	ID            uint      `json:"id"`
	ProductID     string    `json:"product_id"`
	WarehouseID   uint      `json:"warehouse_id"`
	LocationID    *uint     `json:"location_id"`
	BusinessID    uint      `json:"business_id"`
	Quantity      int       `json:"quantity"`
	ReservedQty   int       `json:"reserved_qty"`
	AvailableQty  int       `json:"available_qty"`
	MinStock      *int      `json:"min_stock"`
	MaxStock      *int      `json:"max_stock"`
	ReorderPoint  *int      `json:"reorder_point"`
	ProductName   string    `json:"product_name,omitempty"`
	ProductSKU    string    `json:"product_sku,omitempty"`
	WarehouseName string    `json:"warehouse_name,omitempty"`
	WarehouseCode string    `json:"warehouse_code,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// StockMovementResponse respuesta de movimiento de stock
type StockMovementResponse struct {
	ID               uint      `json:"id"`
	ProductID        string    `json:"product_id"`
	WarehouseID      uint      `json:"warehouse_id"`
	LocationID       *uint     `json:"location_id"`
	BusinessID       uint      `json:"business_id"`
	MovementTypeID   uint      `json:"movement_type_id"`
	MovementTypeCode string    `json:"movement_type_code"`
	MovementTypeName string    `json:"movement_type_name"`
	Reason           string    `json:"reason"`
	Quantity         int       `json:"quantity"`
	PreviousQty      int       `json:"previous_qty"`
	NewQty           int       `json:"new_qty"`
	ReferenceType    *string   `json:"reference_type"`
	ReferenceID      *string   `json:"reference_id"`
	IntegrationID    *uint     `json:"integration_id"`
	Notes            string    `json:"notes"`
	CreatedByID      *uint     `json:"created_by_id"`
	ProductName      string    `json:"product_name,omitempty"`
	ProductSKU       string    `json:"product_sku,omitempty"`
	WarehouseName    string    `json:"warehouse_name,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

// InventoryListResponse respuesta paginada de niveles de inventario
type InventoryListResponse struct {
	Data       []InventoryLevelResponse `json:"data"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	PageSize   int                      `json:"page_size"`
	TotalPages int                      `json:"total_pages"`
}

// MovementListResponse respuesta paginada de movimientos
type MovementListResponse struct {
	Data       []StockMovementResponse `json:"data"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	PageSize   int                     `json:"page_size"`
	TotalPages int                     `json:"total_pages"`
}

// InventoryLevelFromEntity convierte entidad a response
func InventoryLevelFromEntity(e *entities.InventoryLevel) InventoryLevelResponse {
	return InventoryLevelResponse{
		ID:            e.ID,
		ProductID:     e.ProductID,
		WarehouseID:   e.WarehouseID,
		LocationID:    e.LocationID,
		BusinessID:    e.BusinessID,
		Quantity:      e.Quantity,
		ReservedQty:   e.ReservedQty,
		AvailableQty:  e.AvailableQty,
		MinStock:      e.MinStock,
		MaxStock:      e.MaxStock,
		ReorderPoint:  e.ReorderPoint,
		ProductName:   e.ProductName,
		ProductSKU:    e.ProductSKU,
		WarehouseName: e.WarehouseName,
		WarehouseCode: e.WarehouseCode,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

// StockMovementFromEntity convierte entidad a response
func StockMovementFromEntity(e *entities.StockMovement) StockMovementResponse {
	return StockMovementResponse{
		ID:               e.ID,
		ProductID:        e.ProductID,
		WarehouseID:      e.WarehouseID,
		LocationID:       e.LocationID,
		BusinessID:       e.BusinessID,
		MovementTypeID:   e.MovementTypeID,
		MovementTypeCode: e.MovementTypeCode,
		MovementTypeName: e.MovementTypeName,
		Reason:           e.Reason,
		Quantity:         e.Quantity,
		PreviousQty:      e.PreviousQty,
		NewQty:           e.NewQty,
		ReferenceType:    e.ReferenceType,
		ReferenceID:      e.ReferenceID,
		IntegrationID:    e.IntegrationID,
		Notes:            e.Notes,
		CreatedByID:      e.CreatedByID,
		ProductName:      e.ProductName,
		ProductSKU:       e.ProductSKU,
		WarehouseName:    e.WarehouseName,
		CreatedAt:        e.CreatedAt,
	}
}

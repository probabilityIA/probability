package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

type LotResponse struct {
	ID              uint       `json:"id"`
	BusinessID      uint       `json:"business_id"`
	ProductID       string     `json:"product_id"`
	LotCode         string     `json:"lot_code"`
	ManufactureDate *time.Time `json:"manufacture_date"`
	ExpirationDate  *time.Time `json:"expiration_date"`
	ReceivedAt      *time.Time `json:"received_at"`
	SupplierID      *uint      `json:"supplier_id"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type SerialResponse struct {
	ID                uint       `json:"id"`
	BusinessID        uint       `json:"business_id"`
	ProductID         string     `json:"product_id"`
	SerialNumber      string     `json:"serial_number"`
	LotID             *uint      `json:"lot_id"`
	CurrentLocationID *uint      `json:"current_location_id"`
	CurrentStateID    *uint      `json:"current_state_id"`
	ReceivedAt        *time.Time `json:"received_at"`
	SoldAt            *time.Time `json:"sold_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type InventoryStateResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsTerminal  bool   `json:"is_terminal"`
}

type UoMResponse struct {
	ID       uint   `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	IsActive bool   `json:"is_active"`
}

type ProductUoMResponse struct {
	ID               uint    `json:"id"`
	ProductID        string  `json:"product_id"`
	UomID            uint    `json:"uom_id"`
	UomCode          string  `json:"uom_code"`
	UomName          string  `json:"uom_name"`
	ConversionFactor float64 `json:"conversion_factor"`
	IsBase           bool    `json:"is_base"`
	Barcode          string  `json:"barcode"`
	IsActive         bool    `json:"is_active"`
}

type LotListResponse struct {
	Data       []LotResponse `json:"data"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

type SerialListResponse struct {
	Data       []SerialResponse `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

func LotFromEntity(e *entities.InventoryLot) LotResponse {
	return LotResponse{
		ID:              e.ID,
		BusinessID:      e.BusinessID,
		ProductID:       e.ProductID,
		LotCode:         e.LotCode,
		ManufactureDate: e.ManufactureDate,
		ExpirationDate:  e.ExpirationDate,
		ReceivedAt:      e.ReceivedAt,
		SupplierID:      e.SupplierID,
		Status:          e.Status,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

func SerialFromEntity(e *entities.InventorySerial) SerialResponse {
	return SerialResponse{
		ID:                e.ID,
		BusinessID:        e.BusinessID,
		ProductID:         e.ProductID,
		SerialNumber:      e.SerialNumber,
		LotID:             e.LotID,
		CurrentLocationID: e.CurrentLocationID,
		CurrentStateID:    e.CurrentStateID,
		ReceivedAt:        e.ReceivedAt,
		SoldAt:            e.SoldAt,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}

func StateFromEntity(e *entities.InventoryState) InventoryStateResponse {
	return InventoryStateResponse{
		ID:          e.ID,
		Code:        e.Code,
		Name:        e.Name,
		Description: e.Description,
		IsTerminal:  e.IsTerminal,
	}
}

func UoMFromEntity(e *entities.UnitOfMeasure) UoMResponse {
	return UoMResponse{
		ID:       e.ID,
		Code:     e.Code,
		Name:     e.Name,
		Type:     e.Type,
		IsActive: e.IsActive,
	}
}

func ProductUoMFromEntity(e *entities.ProductUoM) ProductUoMResponse {
	return ProductUoMResponse{
		ID:               e.ID,
		ProductID:        e.ProductID,
		UomID:            e.UomID,
		UomCode:          e.UomCode,
		UomName:          e.UomName,
		ConversionFactor: e.ConversionFactor,
		IsBase:           e.IsBase,
		Barcode:          e.Barcode,
		IsActive:         e.IsActive,
	}
}

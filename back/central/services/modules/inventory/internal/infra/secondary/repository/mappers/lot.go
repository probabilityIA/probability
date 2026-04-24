package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func LotModelToEntity(m *models.InventoryLot) *entities.InventoryLot {
	return &entities.InventoryLot{
		ID:              m.ID,
		BusinessID:      m.BusinessID,
		ProductID:       m.ProductID,
		LotCode:         m.LotCode,
		ManufactureDate: m.ManufactureDate,
		ExpirationDate:  m.ExpirationDate,
		ReceivedAt:      m.ReceivedAt,
		SupplierID:      m.SupplierID,
		Status:          m.Status,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func LotEntityToModel(e *entities.InventoryLot) *models.InventoryLot {
	return &models.InventoryLot{
		BusinessID:      e.BusinessID,
		ProductID:       e.ProductID,
		LotCode:         e.LotCode,
		ManufactureDate: e.ManufactureDate,
		ExpirationDate:  e.ExpirationDate,
		ReceivedAt:      e.ReceivedAt,
		SupplierID:      e.SupplierID,
		Status:          e.Status,
	}
}

func SerialModelToEntity(m *models.InventorySerial) *entities.InventorySerial {
	return &entities.InventorySerial{
		ID:                m.ID,
		BusinessID:        m.BusinessID,
		ProductID:         m.ProductID,
		SerialNumber:      m.SerialNumber,
		LotID:             m.LotID,
		CurrentLocationID: m.CurrentLocationID,
		CurrentStateID:    m.CurrentStateID,
		ReceivedAt:        m.ReceivedAt,
		SoldAt:            m.SoldAt,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func SerialEntityToModel(e *entities.InventorySerial) *models.InventorySerial {
	return &models.InventorySerial{
		BusinessID:        e.BusinessID,
		ProductID:         e.ProductID,
		SerialNumber:      e.SerialNumber,
		LotID:             e.LotID,
		CurrentLocationID: e.CurrentLocationID,
		CurrentStateID:    e.CurrentStateID,
		ReceivedAt:        e.ReceivedAt,
		SoldAt:            e.SoldAt,
	}
}

func StateModelToEntity(m *models.InventoryState) *entities.InventoryState {
	return &entities.InventoryState{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		IsTerminal:  m.IsTerminal,
		IsActive:    m.IsActive,
	}
}

func UoMModelToEntity(m *models.UnitOfMeasure) *entities.UnitOfMeasure {
	return &entities.UnitOfMeasure{
		ID:       m.ID,
		Code:     m.Code,
		Name:     m.Name,
		Type:     m.Type,
		IsActive: m.IsActive,
	}
}

func ProductUoMModelToEntity(m *models.ProductUoM) *entities.ProductUoM {
	e := &entities.ProductUoM{
		ID:               m.ID,
		ProductID:        m.ProductID,
		UomID:            m.UomID,
		BusinessID:       m.BusinessID,
		ConversionFactor: m.ConversionFactor,
		IsBase:           m.IsBase,
		Barcode:          m.Barcode,
		IsActive:         m.IsActive,
	}
	if m.Uom.ID != 0 {
		e.UomCode = m.Uom.Code
		e.UomName = m.Uom.Name
		e.UomType = m.Uom.Type
	}
	return e
}

func ProductUoMEntityToModel(e *entities.ProductUoM) *models.ProductUoM {
	return &models.ProductUoM{
		ProductID:        e.ProductID,
		UomID:            e.UomID,
		BusinessID:       e.BusinessID,
		ConversionFactor: e.ConversionFactor,
		IsBase:           e.IsBase,
		Barcode:          e.Barcode,
		IsActive:         e.IsActive,
	}
}

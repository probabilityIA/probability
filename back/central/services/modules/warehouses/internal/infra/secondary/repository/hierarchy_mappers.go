package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func zoneModelToEntity(m *models.WarehouseZone) *entities.WarehouseZone {
	return &entities.WarehouseZone{
		ID:          m.ID,
		WarehouseID: m.WarehouseID,
		BusinessID:  m.BusinessID,
		Code:        m.Code,
		Name:        m.Name,
		Purpose:     m.Purpose,
		IsActive:    m.IsActive,
		ColorHex:    m.ColorHex,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func aisleModelToEntity(m *models.WarehouseAisle) *entities.WarehouseAisle {
	return &entities.WarehouseAisle{
		ID:         m.ID,
		ZoneID:     m.ZoneID,
		BusinessID: m.BusinessID,
		Code:       m.Code,
		Name:       m.Name,
		IsActive:   m.IsActive,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func rackModelToEntity(m *models.WarehouseRack) *entities.WarehouseRack {
	return &entities.WarehouseRack{
		ID:          m.ID,
		AisleID:     m.AisleID,
		BusinessID:  m.BusinessID,
		Code:        m.Code,
		Name:        m.Name,
		LevelsCount: m.LevelsCount,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func rackLevelModelToEntity(m *models.WarehouseRackLevel) *entities.WarehouseRackLevel {
	return &entities.WarehouseRackLevel{
		ID:         m.ID,
		RackID:     m.RackID,
		BusinessID: m.BusinessID,
		Code:       m.Code,
		Ordinal:    m.Ordinal,
		IsActive:   m.IsActive,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

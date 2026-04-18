package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func PutawayRuleModelToEntity(m *models.PutawayRule) *entities.PutawayRule {
	return &entities.PutawayRule{
		ID:           m.ID,
		BusinessID:   m.BusinessID,
		ProductID:    m.ProductID,
		CategoryID:   m.CategoryID,
		TargetZoneID: m.TargetZoneID,
		Priority:     m.Priority,
		Strategy:     m.Strategy,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func PutawayRuleEntityToModel(e *entities.PutawayRule) *models.PutawayRule {
	return &models.PutawayRule{
		BusinessID:   e.BusinessID,
		ProductID:    e.ProductID,
		CategoryID:   e.CategoryID,
		TargetZoneID: e.TargetZoneID,
		Priority:     e.Priority,
		Strategy:     e.Strategy,
		IsActive:     e.IsActive,
	}
}

func PutawaySuggestionModelToEntity(m *models.PutawaySuggestion) *entities.PutawaySuggestion {
	return &entities.PutawaySuggestion{
		ID:                    m.ID,
		BusinessID:            m.BusinessID,
		ProductID:             m.ProductID,
		RecommendedLocationID: m.RecommendedLocationID,
		Quantity:              m.Quantity,
		Status:                m.Status,
		RuleID:                m.RuleID,
		Reason:                m.Reason,
		ActualLocationID:      m.ActualLocationID,
		ConfirmedAt:           m.ConfirmedAt,
		ConfirmedByID:         m.ConfirmedByID,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
	}
}

func PutawaySuggestionEntityToModel(e *entities.PutawaySuggestion) *models.PutawaySuggestion {
	return &models.PutawaySuggestion{
		BusinessID:            e.BusinessID,
		ProductID:             e.ProductID,
		RecommendedLocationID: e.RecommendedLocationID,
		Quantity:              e.Quantity,
		Status:                e.Status,
		RuleID:                e.RuleID,
		Reason:                e.Reason,
		ActualLocationID:      e.ActualLocationID,
		ConfirmedAt:           e.ConfirmedAt,
		ConfirmedByID:         e.ConfirmedByID,
	}
}

func ReplenishmentTaskModelToEntity(m *models.ReplenishmentTask) *entities.ReplenishmentTask {
	return &entities.ReplenishmentTask{
		ID:             m.ID,
		BusinessID:     m.BusinessID,
		ProductID:      m.ProductID,
		WarehouseID:    m.WarehouseID,
		FromLocationID: m.FromLocationID,
		ToLocationID:   m.ToLocationID,
		Quantity:       m.Quantity,
		Status:         m.Status,
		TriggeredBy:    m.TriggeredBy,
		AssignedToID:   m.AssignedToID,
		AssignedAt:     m.AssignedAt,
		CompletedAt:    m.CompletedAt,
		Notes:          m.Notes,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func ReplenishmentTaskEntityToModel(e *entities.ReplenishmentTask) *models.ReplenishmentTask {
	return &models.ReplenishmentTask{
		BusinessID:     e.BusinessID,
		ProductID:      e.ProductID,
		WarehouseID:    e.WarehouseID,
		FromLocationID: e.FromLocationID,
		ToLocationID:   e.ToLocationID,
		Quantity:       e.Quantity,
		Status:         e.Status,
		TriggeredBy:    e.TriggeredBy,
		AssignedToID:   e.AssignedToID,
		AssignedAt:     e.AssignedAt,
		CompletedAt:    e.CompletedAt,
		Notes:          e.Notes,
	}
}

func CrossDockLinkModelToEntity(m *models.CrossDockLink) *entities.CrossDockLink {
	return &entities.CrossDockLink{
		ID:                m.ID,
		BusinessID:        m.BusinessID,
		InboundShipmentID: m.InboundShipmentID,
		OutboundOrderID:   m.OutboundOrderID,
		ProductID:         m.ProductID,
		Quantity:          m.Quantity,
		Status:            m.Status,
		ExecutedAt:        m.ExecutedAt,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func CrossDockLinkEntityToModel(e *entities.CrossDockLink) *models.CrossDockLink {
	return &models.CrossDockLink{
		BusinessID:        e.BusinessID,
		InboundShipmentID: e.InboundShipmentID,
		OutboundOrderID:   e.OutboundOrderID,
		ProductID:         e.ProductID,
		Quantity:          e.Quantity,
		Status:            e.Status,
		ExecutedAt:        e.ExecutedAt,
	}
}

func ProductVelocityModelToEntity(m *models.ProductVelocity) *entities.ProductVelocity {
	return &entities.ProductVelocity{
		ID:          m.ID,
		BusinessID:  m.BusinessID,
		ProductID:   m.ProductID,
		WarehouseID: m.WarehouseID,
		Period:      m.Period,
		UnitsMoved:  m.UnitsMoved,
		Rank:        m.Rank,
		ComputedAt:  m.ComputedAt,
	}
}

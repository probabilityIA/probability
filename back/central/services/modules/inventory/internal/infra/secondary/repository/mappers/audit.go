package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func CountPlanModelToEntity(m *models.CycleCountPlan) *entities.CycleCountPlan {
	return &entities.CycleCountPlan{
		ID:            m.ID,
		BusinessID:    m.BusinessID,
		WarehouseID:   m.WarehouseID,
		Name:          m.Name,
		Strategy:      m.Strategy,
		FrequencyDays: m.FrequencyDays,
		NextRunAt:     m.NextRunAt,
		IsActive:      m.IsActive,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func CountPlanEntityToModel(e *entities.CycleCountPlan) *models.CycleCountPlan {
	return &models.CycleCountPlan{
		BusinessID:    e.BusinessID,
		WarehouseID:   e.WarehouseID,
		Name:          e.Name,
		Strategy:      e.Strategy,
		FrequencyDays: e.FrequencyDays,
		NextRunAt:     e.NextRunAt,
		IsActive:      e.IsActive,
	}
}

func CountTaskModelToEntity(m *models.CycleCountTask) *entities.CycleCountTask {
	return &entities.CycleCountTask{
		ID:           m.ID,
		PlanID:       m.PlanID,
		BusinessID:   m.BusinessID,
		WarehouseID:  m.WarehouseID,
		ScopeType:    m.ScopeType,
		ScopeID:      m.ScopeID,
		Status:       m.Status,
		AssignedToID: m.AssignedToID,
		StartedAt:    m.StartedAt,
		FinishedAt:   m.FinishedAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func CountTaskEntityToModel(e *entities.CycleCountTask) *models.CycleCountTask {
	return &models.CycleCountTask{
		PlanID:       e.PlanID,
		BusinessID:   e.BusinessID,
		WarehouseID:  e.WarehouseID,
		ScopeType:    e.ScopeType,
		ScopeID:      e.ScopeID,
		Status:       e.Status,
		AssignedToID: e.AssignedToID,
		StartedAt:    e.StartedAt,
		FinishedAt:   e.FinishedAt,
	}
}

func CountLineModelToEntity(m *models.CycleCountLine) *entities.CycleCountLine {
	return &entities.CycleCountLine{
		ID:          m.ID,
		TaskID:      m.TaskID,
		BusinessID:  m.BusinessID,
		ProductID:   m.ProductID,
		LocationID:  m.LocationID,
		LotID:       m.LotID,
		ExpectedQty: m.ExpectedQty,
		CountedQty:  m.CountedQty,
		Variance:    m.Variance,
		Status:      m.Status,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func CountLineEntityToModel(e *entities.CycleCountLine) *models.CycleCountLine {
	return &models.CycleCountLine{
		TaskID:      e.TaskID,
		BusinessID:  e.BusinessID,
		ProductID:   e.ProductID,
		LocationID:  e.LocationID,
		LotID:       e.LotID,
		ExpectedQty: e.ExpectedQty,
		CountedQty:  e.CountedQty,
		Variance:    e.Variance,
		Status:      e.Status,
	}
}

func DiscrepancyModelToEntity(m *models.InventoryDiscrepancy) *entities.InventoryDiscrepancy {
	return &entities.InventoryDiscrepancy{
		ID:                   m.ID,
		TaskID:               m.TaskID,
		LineID:               m.LineID,
		BusinessID:           m.BusinessID,
		Status:               m.Status,
		ResolutionMovementID: m.ResolutionMovementID,
		ReviewedByID:         m.ReviewedByID,
		ReviewedAt:           m.ReviewedAt,
		Notes:                m.Notes,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

func DiscrepancyEntityToModel(e *entities.InventoryDiscrepancy) *models.InventoryDiscrepancy {
	return &models.InventoryDiscrepancy{
		TaskID:               e.TaskID,
		LineID:               e.LineID,
		BusinessID:           e.BusinessID,
		Status:               e.Status,
		ResolutionMovementID: e.ResolutionMovementID,
		ReviewedByID:         e.ReviewedByID,
		ReviewedAt:           e.ReviewedAt,
		Notes:                e.Notes,
	}
}

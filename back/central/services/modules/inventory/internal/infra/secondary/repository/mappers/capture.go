package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func LPNModelToEntity(m *models.LicensePlate) *entities.LicensePlate {
	lpn := &entities.LicensePlate{
		ID:                m.ID,
		BusinessID:        m.BusinessID,
		Code:              m.Code,
		LpnType:           m.LpnType,
		CurrentLocationID: m.CurrentLocationID,
		Status:            m.Status,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
	if len(m.Lines) > 0 {
		lpn.Lines = make([]entities.LicensePlateLine, len(m.Lines))
		for i := range m.Lines {
			lpn.Lines[i] = *LPNLineModelToEntity(&m.Lines[i])
		}
	}
	return lpn
}

func LPNEntityToModel(e *entities.LicensePlate) *models.LicensePlate {
	return &models.LicensePlate{
		BusinessID:        e.BusinessID,
		Code:              e.Code,
		LpnType:           e.LpnType,
		CurrentLocationID: e.CurrentLocationID,
		Status:            e.Status,
	}
}

func LPNLineModelToEntity(m *models.LicensePlateLine) *entities.LicensePlateLine {
	return &entities.LicensePlateLine{
		ID:         m.ID,
		LpnID:      m.LpnID,
		BusinessID: m.BusinessID,
		ProductID:  m.ProductID,
		LotID:      m.LotID,
		SerialID:   m.SerialID,
		Qty:        m.Qty,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func LPNLineEntityToModel(e *entities.LicensePlateLine) *models.LicensePlateLine {
	return &models.LicensePlateLine{
		LpnID:      e.LpnID,
		BusinessID: e.BusinessID,
		ProductID:  e.ProductID,
		LotID:      e.LotID,
		SerialID:   e.SerialID,
		Qty:        e.Qty,
	}
}

func ScanEventModelToEntity(m *models.ScanEvent) *entities.ScanEvent {
	return &entities.ScanEvent{
		ID:          m.ID,
		BusinessID:  m.BusinessID,
		UserID:      m.UserID,
		DeviceID:    m.DeviceID,
		ScannedCode: m.ScannedCode,
		CodeType:    m.CodeType,
		Action:      m.Action,
		ScannedAt:   m.ScannedAt,
		CreatedAt:   m.CreatedAt,
	}
}

func SyncLogModelToEntity(m *models.InventorySyncLog) *entities.InventorySyncLog {
	return &entities.InventorySyncLog{
		ID:            m.ID,
		BusinessID:    m.BusinessID,
		IntegrationID: m.IntegrationID,
		Direction:     m.Direction,
		PayloadHash:   m.PayloadHash,
		Status:        m.Status,
		Error:         m.Error,
		SyncedAt:      m.SyncedAt,
		CreatedAt:     m.CreatedAt,
	}
}

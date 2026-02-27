package models

import (
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"gorm.io/gorm"
)

// IntegrationChannelStatus es el modelo local para integration_channel_statuses
type IntegrationChannelStatus struct {
	gorm.Model

	IntegrationTypeID uint
	Code              string
	Name              string
	Description       string
	IsActive          bool
	DisplayOrder      int

	IntegrationType IntegrationType `gorm:"foreignKey:IntegrationTypeID"`
}

// TableName especifica el nombre de la tabla
func (IntegrationChannelStatus) TableName() string {
	return "integration_channel_statuses"
}

// ToDomain convierte el modelo a entidad de dominio
func (m *IntegrationChannelStatus) ToDomain() entities.ChannelStatusInfo {
	result := entities.ChannelStatusInfo{
		ID:                m.ID,
		IntegrationTypeID: m.IntegrationTypeID,
		Code:              m.Code,
		Name:              m.Name,
		Description:       m.Description,
		IsActive:          m.IsActive,
		DisplayOrder:      m.DisplayOrder,
	}
	if m.IntegrationType.ID > 0 {
		it := m.IntegrationType.ToDomain()
		result.IntegrationType = &it
	}
	return result
}

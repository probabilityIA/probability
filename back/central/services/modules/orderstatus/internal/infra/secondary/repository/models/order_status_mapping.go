package models

import (
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// OrderStatusMapping es el modelo GORM para order_status_mappings
// Contiene tags GORM para mapeo de BD
type OrderStatusMapping struct {
	gorm.Model

	// Mapeo
	IntegrationTypeID uint   `gorm:"not null;index;uniqueIndex:idx_status_mapping,priority:1"`
	OriginalStatus    string `gorm:"size:128;not null;uniqueIndex:idx_status_mapping,priority:2"`
	OrderStatusID     uint   `gorm:"not null;index"`

	// Configuración
	IsActive    bool           `gorm:"default:true;index"`
	Description string         `gorm:"type:text"`
	Metadata    datatypes.JSON `gorm:"type:jsonb"`

	// Relaciones
	IntegrationType IntegrationType `gorm:"foreignKey:IntegrationTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OrderStatus     OrderStatus     `gorm:"foreignKey:OrderStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName especifica el nombre de la tabla
func (OrderStatusMapping) TableName() string {
	return "order_status_mappings"
}

// ToDomain convierte el modelo de infra a entidad de dominio
func (m *OrderStatusMapping) ToDomain() entities.OrderStatusMapping {
	result := entities.OrderStatusMapping{
		ID:                m.ID,
		IntegrationTypeID: m.IntegrationTypeID,
		OriginalStatus:    m.OriginalStatus,
		OrderStatusID:     m.OrderStatusID,
		IsActive:          m.IsActive,
		Description:       m.Description,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}

	// Incluir información del IntegrationType si está cargado
	if m.IntegrationType.ID != 0 {
		result.IntegrationType = &entities.IntegrationTypeInfo{
			ID:       m.IntegrationType.ID,
			Code:     m.IntegrationType.Code,
			Name:     m.IntegrationType.Name,
			ImageURL: m.IntegrationType.ImageURL,
		}
	}

	// Incluir información del OrderStatus si está cargado
	if m.OrderStatus.ID != 0 {
		result.OrderStatus = &entities.OrderStatusInfo{
			ID:          m.OrderStatus.ID,
			Code:        m.OrderStatus.Code,
			Name:        m.OrderStatus.Name,
			Description: m.OrderStatus.Description,
			Category:    m.OrderStatus.Category,
			Color:       m.OrderStatus.Color,
			Priority:    m.OrderStatus.Priority,
		}
	}

	return result
}

// FromDomain convierte entidad de dominio a modelo de infra
func FromDomain(e entities.OrderStatusMapping) *OrderStatusMapping {
	return &OrderStatusMapping{
		Model: gorm.Model{
			ID:        e.ID,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
		IntegrationTypeID: e.IntegrationTypeID,
		OriginalStatus:    e.OriginalStatus,
		OrderStatusID:     e.OrderStatusID,
		IsActive:          e.IsActive,
		Description:       e.Description,
	}
}

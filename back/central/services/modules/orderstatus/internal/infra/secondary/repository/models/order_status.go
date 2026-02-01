package models

import (
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// OrderStatus es el modelo GORM para order_statuses
type OrderStatus struct {
	gorm.Model

	Code        string `gorm:"size:64;unique;not null;index"`
	Name        string `gorm:"size:128;not null"`
	Description string `gorm:"type:text"`
	Category    string `gorm:"size:64;index"`
	IsActive    bool   `gorm:"default:true;index"`

	// UI/UX
	Icon     string         `gorm:"size:255"`
	Color    string         `gorm:"size:32"`
	Metadata datatypes.JSON `gorm:"type:jsonb"`
}

// TableName especifica el nombre de la tabla
func (OrderStatus) TableName() string {
	return "order_statuses"
}

// ToDomain convierte el modelo de infra a entidad de dominio
func (m *OrderStatus) ToDomain() entities.OrderStatusInfo {
	return entities.OrderStatusInfo{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		Category:    m.Category,
		Color:       m.Color,
	}
}

package models

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/domain/entities"
	"gorm.io/gorm"
)

// PaymentStatus modelo GORM para la tabla payment_statuses
type PaymentStatus struct {
	ID          uint           `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Campos del dominio
	Code        string `gorm:"size:64;unique;not null;index"`
	Name        string `gorm:"size:128;not null"`
	Description string `gorm:"type:text"`
	Category    string `gorm:"size:64;index"`
	IsActive    bool   `gorm:"default:true;index"`
	Icon        string `gorm:"size:255"`
	Color       string `gorm:"size:32"`
}

// TableName define el nombre de la tabla en la BD
func (PaymentStatus) TableName() string {
	return "payment_statuses"
}

// ToDomain convierte el modelo de infra a entidad de dominio
func (m *PaymentStatus) ToDomain() entities.PaymentStatus {
	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	return entities.PaymentStatus{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		Category:    m.Category,
		IsActive:    m.IsActive,
		Icon:        m.Icon,
		Color:       m.Color,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}

// FromDomain convierte entidad de dominio a modelo de infra
func FromDomain(e entities.PaymentStatus) *PaymentStatus {
	model := &PaymentStatus{
		ID:          e.ID,
		Code:        e.Code,
		Name:        e.Name,
		Description: e.Description,
		Category:    e.Category,
		IsActive:    e.IsActive,
		Icon:        e.Icon,
		Color:       e.Color,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}

	if e.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{
			Time:  *e.DeletedAt,
			Valid: true,
		}
	}

	return model
}

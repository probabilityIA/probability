package models

import (
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"gorm.io/gorm"
)

// IntegrationType es el modelo GORM para integration_types
type IntegrationType struct {
	gorm.Model

	Code     string `gorm:"size:50;unique;not null;index"`
	Name     string `gorm:"size:100;not null"`
	ImageURL string `gorm:"size:500"`
}

// TableName especifica el nombre de la tabla
func (IntegrationType) TableName() string {
	return "integration_types"
}

// ToDomain convierte el modelo a entidad de dominio
func (m *IntegrationType) ToDomain() entities.IntegrationTypeInfo {
	return entities.IntegrationTypeInfo{
		ID:       m.ID,
		Code:     m.Code,
		Name:     m.Name,
		ImageURL: m.ImageURL,
	}
}

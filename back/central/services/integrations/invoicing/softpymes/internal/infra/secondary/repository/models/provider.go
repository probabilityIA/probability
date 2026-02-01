package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Provider representa una instancia configurada de un proveedor de facturación Softpymes
// para un negocio específico. Contiene las credenciales y configuración necesarias.
// MODELO GORM con tags de infraestructura
type Provider struct {
	gorm.Model

	// Relaciones
	BusinessID     uint `gorm:"not null;index;uniqueIndex:idx_business_provider,priority:1"`
	ProviderTypeID uint `gorm:"not null;index;uniqueIndex:idx_business_provider,priority:2"`

	// Identificación
	Name        string `gorm:"size:100;not null"` // "Softpymes - Tienda Principal"
	Description string `gorm:"size:500"`          // Descripción opcional

	// Estado
	IsActive  bool `gorm:"default:true;index"`  // Si está activa
	IsDefault bool `gorm:"default:false;index"` // Si es el proveedor por defecto para este negocio

	// Configuración (JSON - no contiene información sensible)
	// Ejemplo Softpymes: {"referer": "900123456", "branch_code": "001"}
	Config datatypes.JSON `gorm:"type:jsonb"`

	// Credenciales (JSON - almacenadas sin encriptar por ahora, integrationCore las manejará)
	// Ejemplo: {"api_key": "value", "api_secret": "value"}
	Credentials datatypes.JSON `gorm:"type:jsonb"`

	// Metadata
	CreatedByID uint  `gorm:"index"`
	UpdatedByID *uint `gorm:"index"`
}

// TableName especifica el nombre de la tabla para Provider
func (Provider) TableName() string {
	return "invoicing_providers"
}

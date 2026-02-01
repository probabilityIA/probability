package models

import (
	"gorm.io/gorm"
)

// ProviderType representa un tipo de proveedor de facturación electrónica (Softpymes, Siigo, etc.)
// MODELO GORM con tags de infraestructura
type ProviderType struct {
	gorm.Model

	// Identificación
	Name        string `gorm:"size:100;not null;unique"`       // "Softpymes"
	Code        string `gorm:"size:50;not null;unique;index"` // "softpymes"
	Description string `gorm:"size:500"`                       // Descripción del proveedor
	Icon        string `gorm:"size:100"`                       // Icono para UI
	ImageURL    string `gorm:"size:500"`                       // URL del logo

	// Información del proveedor
	ApiBaseURL         string `gorm:"size:255"`              // URL base de la API
	DocumentationURL   string `gorm:"size:500"`              // URL de documentación
	IsActive           bool   `gorm:"default:true;index"`    // Si está activo y disponible
	SupportedCountries string `gorm:"size:255;default:'CO'"` // Países soportados (CSV)
}

// TableName especifica el nombre de la tabla para ProviderType
func (ProviderType) TableName() string {
	return "invoicing_provider_types"
}

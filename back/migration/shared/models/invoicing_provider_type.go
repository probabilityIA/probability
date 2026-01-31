package models

import "gorm.io/gorm"

// ───────────────────────────────────────────
//
//	INVOICING PROVIDER TYPES - Tipos de proveedores de facturación
//
// ───────────────────────────────────────────

// InvoicingProviderType representa un tipo de proveedor de facturación electrónica
// Ejemplo: Softpymes, Siigo, Facturama, Alegra, etc.
type InvoicingProviderType struct {
	gorm.Model

	// Identificación
	Name        string `gorm:"size:100;not null;unique"` // "Softpymes", "Siigo", "Facturama"
	Code        string `gorm:"size:50;not null;unique"`  // "softpymes", "siigo", "facturama"
	Description string `gorm:"size:500"`                 // Descripción del proveedor
	Icon        string `gorm:"size:100"`                 // Icono para UI
	ImageURL    string `gorm:"size:500"`                 // URL del logo del proveedor

	// Información del proveedor
	ApiBaseURL      string `gorm:"size:255"`        // URL base de la API
	DocumentationURL string `gorm:"size:255"`       // URL de la documentación
	IsActive        bool   `gorm:"default:true"`    // Si el tipo está activo y disponible
	SupportedCountries string `gorm:"size:500"`    // Países soportados (separados por coma: "CO,MX,PE")

	// Relaciones
	InvoicingProviders []InvoicingProvider `gorm:"foreignKey:ProviderTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName especifica el nombre de la tabla para InvoicingProviderType
func (InvoicingProviderType) TableName() string {
	return "invoicing_provider_types"
}

package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ───────────────────────────────────────────
//
//	INVOICING PROVIDERS - Instancias de proveedores de facturación por negocio
//
// ───────────────────────────────────────────

// InvoicingProvider representa una instancia configurada de un proveedor de facturación
// para un negocio específico. Contiene las credenciales y configuración necesarias.
type InvoicingProvider struct {
	gorm.Model

	// Relación con Business
	BusinessID uint     `gorm:"not null;index;uniqueIndex:idx_business_provider,priority:1"`
	Business   Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Relación con InvoicingProviderType
	ProviderTypeID uint                   `gorm:"not null;index;uniqueIndex:idx_business_provider,priority:2"`
	ProviderType   InvoicingProviderType `gorm:"foreignKey:ProviderTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Identificación
	Name        string `gorm:"size:100;not null"` // "Softpymes - Tienda Principal"
	Description string `gorm:"size:500"`          // Descripción opcional

	// Estado
	IsActive  bool `gorm:"default:true;index"`  // Si está activa
	IsDefault bool `gorm:"default:false;index"` // Si es el proveedor por defecto para este negocio

	// Configuración (JSON - no contiene información sensible)
	// Ejemplo Softpymes: {"referer": "900123456", "branch_code": "001"}
	// Ejemplo Siigo: {"document_type": "FV", "numeration_id": "123"}
	Config datatypes.JSON `gorm:"type:jsonb"`

	// Credenciales encriptadas (JSON)
	// Contiene API keys, secrets, tokens encriptados con AES-256
	// Ejemplo: {"api_key": "encrypted_value", "api_secret": "encrypted_value"}
	Credentials datatypes.JSON `gorm:"type:jsonb"`

	// Metadata
	CreatedByID uint  `gorm:"index"`
	UpdatedByID *uint `gorm:"index"`

	// Relaciones
	CreatedBy       User                 `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedBy       *User                `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	InvoicingConfigs []InvoicingConfig   `gorm:"foreignKey:InvoicingProviderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Invoices        []Invoice            `gorm:"foreignKey:InvoicingProviderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName especifica el nombre de la tabla para InvoicingProvider
func (InvoicingProvider) TableName() string {
	return "invoicing_providers"
}

package entities

import "time"

// InvoicingProvider representa una instancia configurada de un proveedor de facturación
// para un negocio específico. Contiene las credenciales y configuración necesarias.
// Entidad PURA de dominio - SIN TAGS de infraestructura
type InvoicingProvider struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Relaciones (solo IDs)
	BusinessID     uint
	ProviderTypeID uint

	// Identificación
	Name        string
	Description string

	// Estado
	IsActive  bool
	IsDefault bool

	// Configuración (será JSON en infraestructura)
	Config      map[string]interface{}
	Credentials map[string]interface{} // Encriptado en infraestructura

	// Metadata
	CreatedByID uint
	UpdatedByID *uint
}

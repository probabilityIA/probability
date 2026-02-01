package entities

import "time"

// ProviderType representa un tipo de proveedor de facturación electrónica (Softpymes, etc.)
// Entidad PURA de dominio - SIN TAGS de infraestructura
type ProviderType struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Identificación
	Name        string
	Code        string
	Description string
	Icon        string
	ImageURL    string

	// Información del proveedor
	ApiBaseURL         string
	DocumentationURL   string
	IsActive           bool
	SupportedCountries string
}

package domain

import "gorm.io/datatypes"

// CreateIntegrationDTO representa los datos para crear una integración
type CreateIntegrationDTO struct {
	Name        string
	Code        string
	Type        string
	Category    string
	BusinessID  *uint
	IsActive    bool
	IsDefault   bool
	Config      datatypes.JSON
	Credentials map[string]interface{} // Se encriptará antes de guardar
	Description string
	CreatedByID uint
}

// UpdateIntegrationDTO representa los datos para actualizar una integración
type UpdateIntegrationDTO struct {
	Name        *string
	Code        *string
	IsActive    *bool
	IsDefault   *bool
	Config      *datatypes.JSON
	Credentials *map[string]interface{} // Se encriptará antes de guardar
	Description *string
	UpdatedByID uint
}

// IntegrationFilters representa los filtros para listar integraciones
type IntegrationFilters struct {
	Page       int
	PageSize   int
	Type       *string
	Category   *string
	BusinessID *uint
	IsActive   *bool
	Search     *string // Búsqueda por nombre o código
}

// IntegrationWithCredentials representa una integración con credenciales desencriptadas
// Solo se usa internamente, nunca se expone en respuestas HTTP
type IntegrationWithCredentials struct {
	Integration
	DecryptedCredentials DecryptedCredentials
}

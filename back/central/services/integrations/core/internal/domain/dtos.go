package domain

import (
	"mime/multipart"

	"gorm.io/datatypes"
)

// CreateIntegrationTypeDTO representa los datos para crear un tipo de integración
type CreateIntegrationTypeDTO struct {
	Name              string
	Code              string
	Description       string
	Icon              string
	Category          string
	IsActive          bool
	ConfigSchema      datatypes.JSON
	CredentialsSchema datatypes.JSON
	ImageFile         *multipart.FileHeader // Archivo de imagen para subir a S3
}

// UpdateIntegrationTypeDTO representa los datos para actualizar un tipo de integración
type UpdateIntegrationTypeDTO struct {
	Name              *string
	Code              *string
	Description       *string
	Icon              *string
	Category          *string
	IsActive          *bool
	ConfigSchema      *datatypes.JSON
	CredentialsSchema *datatypes.JSON
	ImageFile         *multipart.FileHeader // Archivo de imagen para subir a S3
	RemoveImage       bool                  // Flag para eliminar la imagen existente
}

// CreateIntegrationDTO representa los datos para crear una integración
type CreateIntegrationDTO struct {
	Name              string
	Code              string
	IntegrationTypeID uint // ID del tipo de integración (obligatorio)
	Category          string
	BusinessID        *uint
	StoreID           string // Identificador externo (ej: shop domain para Shopify)
	IsActive          bool
	IsDefault         bool
	Config            datatypes.JSON
	Credentials       map[string]interface{} // Se encriptará antes de guardar
	Description       string
	CreatedByID       uint
}

// UpdateIntegrationDTO representa los datos para actualizar una integración
type UpdateIntegrationDTO struct {
	Name              *string
	Code              *string
	IntegrationTypeID *uint   // ID del tipo de integración (opcional en update)
	StoreID           *string // Identificador externo (ej: shop domain para Shopify)
	IsActive          *bool
	IsDefault         *bool
	Config            *datatypes.JSON
	Credentials       *map[string]interface{} // Se encriptará antes de guardar
	Description       *string
	UpdatedByID       uint
}

// IntegrationFilters representa los filtros para listar integraciones
type IntegrationFilters struct {
	Page                int
	PageSize            int
	IntegrationTypeID   *uint   // Filtrar por ID del tipo de integración
	IntegrationTypeCode *string // Filtrar por código del tipo de integración (alternativa)
	Category            *string
	BusinessID          *uint
	IsActive            *bool
	Search              *string // Búsqueda por nombre o código
	StoreID             *string // Identificador externo (p.e. shop domain)
}

// IntegrationWithCredentials representa una integración con credenciales desencriptadas
// Solo se usa internamente, nunca se expone en respuestas HTTP
type IntegrationWithCredentials struct {
	Integration
	DecryptedCredentials DecryptedCredentials
}

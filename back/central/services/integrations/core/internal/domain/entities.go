package domain

import (
	"time"

	"gorm.io/datatypes"
)

type IntegrationCategory struct {
	ID               uint
	Code             string
	Name             string
	Description      string
	Icon             string
	Color            string
	DisplayOrder     int
	ParentCategoryID *uint
	IsActive         bool
	IsVisible        bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type IntegrationType struct {
	ID                uint
	Name              string
	Code              string
	Description       string
	Icon              string
	ImageURL          string
	CategoryID        uint
	Category          *IntegrationCategory
	IsActive          bool
	InDevelopment     bool
	ConfigSchema      datatypes.JSON
	CredentialsSchema datatypes.JSON
	SetupInstructions string
	BaseURL           string
	BaseURLTest       string
	// Credenciales de plataforma encriptadas — opacas al dominio, procesadas por el use case
	PlatformCredentialsEncrypted []byte
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
}

type Integration struct {
	ID                uint
	Name              string
	Code              string
	Category          string // Código de categoría (derivado de IntegrationType.Category.Code)
	IntegrationTypeID uint
	IntegrationType   *IntegrationType
	BusinessID        *uint
	BusinessName      *string // Nombre del business asociado (cargado via Preload)
	StoreID           string
	IsActive          bool
	IsDefault         bool
	IsTesting         bool // Si está en modo de pruebas (usa base_url_test)
	Config            datatypes.JSON
	Credentials       datatypes.JSON
	Description       string
	CreatedByID       uint
	UpdatedByID       *uint
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type DecryptedCredentials map[string]interface{}

type IntegrationConfig map[string]interface{}

// PublicIntegration representa una integración en formato público (sin credenciales, con config deserializada)
type PublicIntegration struct {
	ID              uint
	BusinessID      *uint
	Name            string
	StoreID         string
	IntegrationType int
	Config          map[string]interface{}
	// URL resolution fields (desde integration_types)
	IsTesting   bool   // Si la integración está en modo pruebas
	BaseURL     string // URL de producción (integration_types.base_url)
	BaseURLTest string // URL de pruebas (integration_types.base_url_test)
}

package domain

import (
	"time"

	"gorm.io/datatypes"

	"github.com/secamc93/probability/back/central/shared/productmatch"
)

type CredentialRevealAudit struct {
	UserID            uint
	BusinessID        uint
	IntegrationTypeID uint
	IntegrationCode   string
	IPAddress         string
	UserAgent         string
}

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
	ProductMatchOptions          datatypes.JSON
	DefaultProductMatchRules     datatypes.JSON
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
	ProductMatchRules datatypes.JSON
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
	IsActive        bool
	Config          map[string]interface{}
	// URL resolution fields (desde integration_types)
	IsTesting   bool   // Si la integración está en modo pruebas
	BaseURL     string // URL de producción (integration_types.base_url)
	BaseURLTest string // URL de pruebas (integration_types.base_url_test)
	// Reglas de match de producto ya resueltas (integracion > tipo > default global),
	// en orden de prioridad. Nunca vacio.
	ProductMatchRules []productmatch.Rule
}

type StoreIDOwner struct {
	IntegrationID uint
	BusinessID    *uint
}

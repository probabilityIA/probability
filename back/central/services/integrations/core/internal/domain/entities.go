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
	ConfigSchema      datatypes.JSON
	CredentialsSchema datatypes.JSON
	SetupInstructions string
	CreatedAt         time.Time
	UpdatedAt         time.Time
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

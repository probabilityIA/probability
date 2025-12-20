package domain

import (
	"time"

	"gorm.io/datatypes"
)

type IntegrationType struct {
	ID                uint
	Name              string
	Code              string
	Description       string
	Icon              string
	ImageURL          string
	Category          string
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
	IntegrationTypeID uint
	IntegrationType   *IntegrationType
	Category          string
	BusinessID        *uint
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

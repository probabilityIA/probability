package domain

import (
	"time"

	"gorm.io/datatypes"
)

// IntegrationType representa un tipo de integración disponible
type IntegrationType struct {
	ID                uint           `json:"id"`
	Name              string         `json:"name"`
	Code              string         `json:"code"`
	Description       string         `json:"description"`
	Icon              string         `json:"icon"`
	Category          string         `json:"category"` // "internal" | "external"
	IsActive          bool           `json:"is_active"`
	ConfigSchema      datatypes.JSON `json:"config_schema"`      // JSON schema para campos de configuración requeridos
	CredentialsSchema datatypes.JSON `json:"credentials_schema"` // JSON schema para credenciales requeridas
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

// Integration representa una integración del sistema
type Integration struct {
	ID                uint             `json:"id"`
	Name              string           `json:"name"`
	Code              string           `json:"code"`
	IntegrationTypeID uint             `json:"integration_type_id"`        // Relación con IntegrationType
	IntegrationType   *IntegrationType `json:"integration_type,omitempty"` // Relación cargada
	Category          string           `json:"category"`                   // "internal" | "external" (redundante pero útil para queries)
	BusinessID        *uint            `json:"business_id"`                // NULL = global (como WhatsApp)
	IsActive          bool             `json:"is_active"`
	IsDefault         bool             `json:"is_default"`
	Config            datatypes.JSON   `json:"config"` // Configuración en JSON (no sensible)
	Credentials       datatypes.JSON   `json:"-"`      // Credenciales encriptadas (no se expone)
	Description       string           `json:"description"`
	CreatedByID       uint             `json:"created_by_id"`
	UpdatedByID       *uint            `json:"updated_by_id"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

// DecryptedCredentials representa las credenciales desencriptadas (solo en memoria)
type DecryptedCredentials map[string]interface{}

// IntegrationConfig representa la configuración de una integración (estructura flexible)
type IntegrationConfig map[string]interface{}

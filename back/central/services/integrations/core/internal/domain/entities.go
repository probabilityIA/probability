package domain

import (
	"time"

	"gorm.io/datatypes"
)

// Integration representa una integración del sistema
type Integration struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Code        string         `json:"code"`
	Type        string         `json:"type"`        // "whatsapp", "shopify", "mercado_libre"
	Category    string         `json:"category"`    // "internal" | "external"
	BusinessID  *uint          `json:"business_id"` // NULL = global (como WhatsApp)
	IsActive    bool           `json:"is_active"`
	IsDefault   bool           `json:"is_default"`
	Config      datatypes.JSON `json:"config"` // Configuración en JSON (no sensible)
	Credentials datatypes.JSON `json:"-"`      // Credenciales encriptadas (no se expone)
	Description string         `json:"description"`
	CreatedByID uint           `json:"created_by_id"`
	UpdatedByID *uint          `json:"updated_by_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// DecryptedCredentials representa las credenciales desencriptadas (solo en memoria)
type DecryptedCredentials map[string]interface{}

// IntegrationConfig representa la configuración de una integración (estructura flexible)
type IntegrationConfig map[string]interface{}

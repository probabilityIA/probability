package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ───────────────────────────────────────────
//
//	PAYMENT METHODS - Métodos de pago del sistema
//
// ───────────────────────────────────────────

// PaymentMethod representa un método de pago en Probability
type PaymentMethod struct {
	gorm.Model

	// Identificación
	Code        string `gorm:"size:64;unique;not null;index"` // "credit_card", "paypal", "cash"
	Name        string `gorm:"size:128;not null"`             // "Tarjeta de Crédito"
	Description string `gorm:"type:text"`                     // Descripción detallada

	// Categorización
	Category string `gorm:"size:64;index"` // "card", "digital_wallet", "bank_transfer", "cash"
	Provider string `gorm:"size:64"`       // "stripe", "paypal", "mercadopago"

	// Configuración
	IsActive    bool `gorm:"default:true;index"` // Si está activo
	RequiresKYC bool `gorm:"default:false"`      // Si requiere verificación de identidad

	// UI/UX
	Icon     string         `gorm:"size:255"`   // URL del ícono
	Color    string         `gorm:"size:32"`    // Color hex para UI
	Metadata datatypes.JSON `gorm:"type:jsonb"` // Metadata adicional
}

// TableName especifica el nombre de la tabla
func (PaymentMethod) TableName() string {
	return "payment_methods"
}

// ───────────────────────────────────────────
//
//	PAYMENT METHOD MAPPINGS - Mapeo por integración
//
// ───────────────────────────────────────────

// PaymentMethodMapping mapea métodos de pago de integraciones externas
// a los métodos de pago unificados de Probability
type PaymentMethodMapping struct {
	gorm.Model

	// Mapeo
	IntegrationType string `gorm:"size:50;not null;index;uniqueIndex:idx_payment_mapping,priority:1"` // "shopify", "whatsapp"
	OriginalMethod  string `gorm:"size:128;not null;uniqueIndex:idx_payment_mapping,priority:2"`      // "shopify_payments"
	PaymentMethodID uint   `gorm:"not null;index"`                                                    // FK a payment_methods

	// Configuración
	IsActive bool           `gorm:"default:true;index"` // Si el mapeo está activo
	Priority int            `gorm:"default:0"`          // Prioridad en caso de múltiples mapeos
	Metadata datatypes.JSON `gorm:"type:jsonb"`         // Metadata adicional del mapeo

	// Relación
	PaymentMethod PaymentMethod `gorm:"foreignKey:PaymentMethodID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName especifica el nombre de la tabla
func (PaymentMethodMapping) TableName() string {
	return "payment_method_mappings"
}

// ───────────────────────────────────────────
//
//	ORDER STATUS MAPPINGS - Mapeo de estados
//
// ───────────────────────────────────────────

// OrderStatusMapping mapea estados de órdenes de integraciones externas
// a los estados unificados de Probability
type OrderStatusMapping struct {
	gorm.Model

	// Mapeo
	IntegrationType string `gorm:"size:50;not null;index;uniqueIndex:idx_status_mapping,priority:1"` // "shopify", "whatsapp"
	OriginalStatus  string `gorm:"size:128;not null;uniqueIndex:idx_status_mapping,priority:2"`      // "paid", "fulfilled"
	MappedStatus    string `gorm:"size:64;not null;index"`                                           // "processing", "shipped"

	// Configuración
	IsActive    bool           `gorm:"default:true;index"` // Si el mapeo está activo
	Priority    int            `gorm:"default:0"`          // Prioridad en caso de múltiples mapeos
	Description string         `gorm:"type:text"`          // Descripción del mapeo
	Metadata    datatypes.JSON `gorm:"type:jsonb"`         // Metadata adicional
}

// TableName especifica el nombre de la tabla
func (OrderStatusMapping) TableName() string {
	return "order_status_mappings"
}

// ───────────────────────────────────────────
//
//	HELPER METHODS
//
// ───────────────────────────────────────────

// IsValidMappedStatus verifica si el estado mapeado es válido
// Debe coincidir con las constantes en orders/domain/status.go
func (m *OrderStatusMapping) IsValidMappedStatus() bool {
	validStatuses := []string{
		"pending", "processing", "shipped", "delivered",
		"completed", "cancelled", "refunded", "failed", "on_hold",
	}

	for _, status := range validStatuses {
		if m.MappedStatus == status {
			return true
		}
	}
	return false
}

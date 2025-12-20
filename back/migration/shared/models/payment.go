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
//	ORDER STATUSES - Estados de órdenes del sistema
//
// ───────────────────────────────────────────

// OrderStatus representa un estado de orden en Probability
type OrderStatus struct {
	gorm.Model

	// Identificación
	Code        string `gorm:"size:64;unique;not null;index"` // "pending", "processing", "completed"
	Name        string `gorm:"size:128;not null"`             // "Pendiente", "En Procesamiento", "Completada"
	Description string `gorm:"type:text"`                     // Descripción detallada del estado

	// Categorización
	Category string `gorm:"size:64;index"` // "active", "completed", "cancelled", "refunded"

	// Configuración
	IsActive bool `gorm:"default:true;index"` // Si está activo

	// UI/UX
	Icon     string         `gorm:"size:255"`   // URL del ícono
	Color    string         `gorm:"size:32"`    // Color hex para UI
	Metadata datatypes.JSON `gorm:"type:jsonb"` // Metadata adicional
}

// TableName especifica el nombre de la tabla
func (OrderStatus) TableName() string {
	return "order_statuses"
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
	IntegrationTypeID uint   `gorm:"not null;index;uniqueIndex:idx_status_mapping,priority:1"`    // FK a integration_types
	OriginalStatus    string `gorm:"size:128;not null;uniqueIndex:idx_status_mapping,priority:2"` // "paid", "fulfilled"
	OrderStatusID     uint   `gorm:"not null;index"`                                              // FK a order_statuses

	// Configuración
	IsActive    bool           `gorm:"default:true;index"` // Si el mapeo está activo
	Priority    int            `gorm:"default:0"`          // Prioridad en caso de múltiples mapeos
	Description string         `gorm:"type:text"`          // Descripción del mapeo
	Metadata    datatypes.JSON `gorm:"type:jsonb"`         // Metadata adicional

	// Relaciones
	IntegrationType IntegrationType `gorm:"foreignKey:IntegrationTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OrderStatus     OrderStatus     `gorm:"foreignKey:OrderStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName especifica el nombre de la tabla
func (OrderStatusMapping) TableName() string {
	return "order_status_mappings"
}

// ───────────────────────────────────────────
//
//	PAYMENT STATUSES - Estados de pago del sistema
//
// ───────────────────────────────────────────

// PaymentStatus representa un estado de pago en Probability
type PaymentStatus struct {
	gorm.Model

	// Identificación
	Code        string `gorm:"size:64;unique;not null;index"` // "pending", "authorized", "paid", "refunded"
	Name        string `gorm:"size:128;not null"`             // "Pendiente", "Autorizado", "Pagado", "Reembolsado"
	Description string `gorm:"type:text"`                     // Descripción detallada del estado

	// Categorización
	Category string `gorm:"size:64;index"` // "pending", "completed", "refunded", "failed"

	// Configuración
	IsActive bool `gorm:"default:true;index"` // Si está activo

	// UI/UX
	Icon     string         `gorm:"size:255"`   // URL del ícono
	Color    string         `gorm:"size:32"`    // Color hex para UI
	Metadata datatypes.JSON `gorm:"type:jsonb"` // Metadata adicional
}

// TableName especifica el nombre de la tabla
func (PaymentStatus) TableName() string {
	return "payment_statuses"
}

// ───────────────────────────────────────────
//
//	FULFILLMENT STATUSES - Estados de fulfillment del sistema
//
// ───────────────────────────────────────────

// FulfillmentStatus representa un estado de fulfillment en Probability
type FulfillmentStatus struct {
	gorm.Model

	// Identificación
	Code        string `gorm:"size:64;unique;not null;index"` // "unfulfilled", "partial", "fulfilled", "shipped"
	Name        string `gorm:"size:128;not null"`             // "No Cumplida", "Parcial", "Cumplida", "Enviada"
	Description string `gorm:"type:text"`                     // Descripción detallada del estado

	// Categorización
	Category string `gorm:"size:64;index"` // "pending", "in_progress", "completed", "cancelled"

	// Configuración
	IsActive bool `gorm:"default:true;index"` // Si está activo

	// UI/UX
	Icon     string         `gorm:"size:255"`   // URL del ícono
	Color    string         `gorm:"size:32"`    // Color hex para UI
	Metadata datatypes.JSON `gorm:"type:jsonb"` // Metadata adicional
}

// TableName especifica el nombre de la tabla
func (FulfillmentStatus) TableName() string {
	return "fulfillment_statuses"
}

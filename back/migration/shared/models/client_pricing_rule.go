package models

import "gorm.io/gorm"

// ClientPricingRule representa un ajuste de precio personalizado por cliente.
// Si ProductID es NULL, la regla aplica globalmente a TODOS los productos del negocio para ese cliente.
// Si ProductID tiene valor, la regla aplica solo a ese producto específico.
// La regla producto-específica tiene prioridad sobre la global (no se apilan).
type ClientPricingRule struct {
	gorm.Model
	BusinessID      uint    `gorm:"not null;index;uniqueIndex:idx_pricing_rule_biz_client_product,priority:1"`
	ClientID        uint    `gorm:"not null;index;uniqueIndex:idx_pricing_rule_biz_client_product,priority:2"`
	ProductID       *string `gorm:"type:varchar(64);index;uniqueIndex:idx_pricing_rule_biz_client_product,priority:3"` // NULL = regla global
	AdjustmentType  string  `gorm:"size:20;not null;default:'percentage'"`                                             // "percentage" o "fixed"
	AdjustmentValue float64 `gorm:"type:decimal(15,2);not null;default:0"`                                             // positivo=incremento, negativo=descuento
	IsActive        bool    `gorm:"default:true;index"`
	Priority        int     `gorm:"default:0"`  // mayor = se evalúa primero en caso de conflicto
	Description     string  `gorm:"size:255"`

	// Relations
	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Client   Client   `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Product  *Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

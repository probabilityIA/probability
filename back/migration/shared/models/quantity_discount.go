package models

import "gorm.io/gorm"

// QuantityDiscount representa un descuento por volumen de compra.
// Si ProductID es NULL, el descuento aplica a todos los productos del negocio.
// Si ProductID tiene valor, aplica solo a ese producto.
// Los descuentos por cantidad son iguales para todos los clientes (no dependen de ClientID).
// Se apilan sobre el precio ajustado del cliente (ClientPricingRule).
// Múltiples tiers: MinQuantity=10 -> 5%, MinQuantity=50 -> 10%, etc.
// Se aplica el tier con mayor MinQuantity que el pedido cumpla.
type QuantityDiscount struct {
	gorm.Model
	BusinessID      uint    `gorm:"not null;index;uniqueIndex:idx_qty_discount_biz_product_min,priority:1"`
	ProductID       *string `gorm:"type:varchar(64);index;uniqueIndex:idx_qty_discount_biz_product_min,priority:2"` // NULL = aplica a todos los productos
	MinQuantity     int     `gorm:"not null;uniqueIndex:idx_qty_discount_biz_product_min,priority:3"`
	DiscountPercent float64 `gorm:"type:decimal(5,2);not null"` // ej: 5.00 = 5% de descuento
	IsActive        bool    `gorm:"default:true;index"`
	Description     string  `gorm:"size:255"`

	// Relations
	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Product  *Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

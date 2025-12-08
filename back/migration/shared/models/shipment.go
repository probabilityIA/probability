package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ───────────────────────────────────────────
//
//	SHIPMENTS - Envíos de la orden
//
// ───────────────────────────────────────────

// Shipment representa un envío asociado a una orden
type Shipment struct {
	gorm.Model

	// Relación con la orden
	OrderID string `gorm:"type:varchar(36);not null;index"` // UUID de la orden

	// Información de tracking
	TrackingNumber *string `gorm:"size:128;index"` // Número de rastreo
	TrackingURL    *string `gorm:"size:512"`       // URL de rastreo
	Carrier        *string `gorm:"size:128"`       // Transportista (ej: "FedEx", "DHL")
	CarrierCode    *string `gorm:"size:50"`        // Código del transportista

	// Información de guía
	GuideID  *string `gorm:"size:128;index"` // ID de guía de envío
	GuideURL *string `gorm:"size:512"`       // URL de la guía

	// Estado del envío
	Status      string     `gorm:"size:64;not null;index;default:'pending'"` // "pending", "in_transit", "delivered", "failed"
	ShippedAt   *time.Time `gorm:"index"`                                    // Cuándo se envió
	DeliveredAt *time.Time // Cuándo se entregó

	// Información de dirección
	ShippingAddressID *uint    // FK a addresses (opcional, puede usar la de la orden)
	ShippingAddress   *Address `gorm:"foreignKey:ShippingAddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	// Costos
	ShippingCost  *float64 `gorm:"type:decimal(12,2)"` // Costo de envío
	InsuranceCost *float64 `gorm:"type:decimal(12,2)"` // Costo de seguro
	TotalCost     *float64 `gorm:"type:decimal(12,2)"` // Costo total

	// Dimensiones y peso
	Weight *float64 `gorm:"type:decimal(10,2)"` // Peso en kg
	Height *float64 `gorm:"type:decimal(10,2)"` // Alto en cm
	Width  *float64 `gorm:"type:decimal(10,2)"` // Ancho en cm
	Length *float64 `gorm:"type:decimal(10,2)"` // Largo en cm

	// Información de fulfillment
	WarehouseID   *uint  `gorm:"index"`         // ID del almacén
	WarehouseName string `gorm:"size:128"`      // Nombre del almacén
	DriverID      *uint  `gorm:"index"`         // ID del conductor
	DriverName    string `gorm:"size:255"`      // Nombre del conductor
	IsLastMile    bool   `gorm:"default:false"` // Si es última milla

	// Información adicional
	EstimatedDelivery *time.Time     `gorm:"index"`      // Entrega estimada
	DeliveryNotes     *string        `gorm:"type:text"`  // Notas de entrega
	Metadata          datatypes.JSON `gorm:"type:jsonb"` // Metadata adicional del canal

	// Relación
	Order Order `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName especifica el nombre de la tabla
func (Shipment) TableName() string {
	return "shipments"
}

package models

import (
	"crypto/rand"
	"fmt"
	mathrand "math/rand"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ───────────────────────────────────────────
//
//	ORDERS - Órdenes unificadas del sistema
//	Modelo desnormalizado para fácil migración a DynamoDB
//
// ───────────────────────────────────────────

// Order representa una orden unificada en el sistema Probability
// Este modelo es auto-contenido y no depende de relaciones externas
// para facilitar la migración nocturna a DynamoDB
type Order struct {
	// ID único de la orden (UUID)
	ID        string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// ============================================
	// IDENTIFICADORES DE INTEGRACIÓN
	// ============================================
	BusinessID      *uint  `gorm:"index"`                  // ID del negocio (null = global)
	IntegrationID   uint   `gorm:"not null;index"`         // ID de la integración
	IntegrationType string `gorm:"size:50;not null;index"` // "shopify", "whatsapp", etc.

	// ============================================
	// IDENTIFICADORES DE LA ORDEN
	// ============================================
	Platform       string `gorm:"size:50;not null;index"`                                                     // Plataforma origen
	ExternalID     string `gorm:"size:255;not null;index;uniqueIndex:idx_integration_external_id,priority:2"` // ID en plataforma externa
	OrderNumber    string `gorm:"size:128;index"`                                                             // Número visible de la orden
	InternalNumber string `gorm:"size:128;unique;index"`                                                      // Número interno Probability

	// ============================================
	// INFORMACIÓN FINANCIERA
	// ============================================
	Subtotal     float64  `gorm:"type:decimal(12,2);not null;default:0"` // Subtotal antes de impuestos
	Tax          float64  `gorm:"type:decimal(12,2);not null;default:0"` // Impuestos
	Discount     float64  `gorm:"type:decimal(12,2);not null;default:0"` // Descuentos
	ShippingCost float64  `gorm:"type:decimal(12,2);not null;default:0"` // Costo de envío
	TotalAmount  float64  `gorm:"type:decimal(12,2);not null"`           // Total final
	Currency     string   `gorm:"size:10;default:'USD'"`                 // Moneda
	CodTotal     *float64 `gorm:"type:decimal(12,2)"`                    // Total para pago contra entrega

	// ============================================
	// INFORMACIÓN DEL CLIENTE (Desnormalizado)
	// ============================================
	CustomerID    *uint  `gorm:"index"`          // ID del cliente (referencia opcional)
	CustomerName  string `gorm:"size:255"`       // Nombre completo
	CustomerEmail string `gorm:"size:255;index"` // Email
	CustomerPhone string `gorm:"size:32"`        // Teléfono
	CustomerDNI   string `gorm:"size:64"`        // DNI/Identificación

	// ============================================
	// DIRECCIÓN DE ENVÍO (Desnormalizado)
	// ============================================
	ShippingStreet     string   `gorm:"size:255"` // Calle y número
	ShippingCity       string   `gorm:"size:128"` // Ciudad
	ShippingState      string   `gorm:"size:128"` // Estado/Provincia
	ShippingCountry    string   `gorm:"size:128"` // País
	ShippingPostalCode string   `gorm:"size:32"`  // Código postal
	ShippingLat        *float64 // Latitud
	ShippingLng        *float64 // Longitud

	// ============================================
	// INFORMACIÓN DE PAGO
	// ============================================
	PaymentMethodID uint       `gorm:"not null;index"`      // FK a payment_methods (REQUERIDO)
	IsPaid          bool       `gorm:"default:false;index"` // Si está pagada
	PaidAt          *time.Time // Cuándo se pagó

	// ============================================
	// INFORMACIÓN DE ENVÍO/LOGÍSTICA
	// ============================================
	TrackingNumber *string    `gorm:"size:128;index"` // Número de rastreo
	TrackingLink   *string    `gorm:"size:512"`       // Link de rastreo
	GuideID        *string    `gorm:"size:128"`       // ID de guía de envío
	GuideLink      *string    `gorm:"size:512"`       // Link de guía
	DeliveryDate   *time.Time `gorm:"index"`          // Fecha de entrega programada
	DeliveredAt    *time.Time // Fecha de entrega real

	// ============================================
	// INFORMACIÓN DE FULFILLMENT
	// ============================================
	WarehouseID   *uint  `gorm:"index"`         // ID del almacén
	WarehouseName string `gorm:"size:128"`      // Nombre del almacén (desnormalizado)
	DriverID      *uint  `gorm:"index"`         // ID del conductor
	DriverName    string `gorm:"size:255"`      // Nombre del conductor (desnormalizado)
	IsLastMile    bool   `gorm:"default:false"` // Si es última milla

	// ============================================
	// DIMENSIONES Y PESO
	// ============================================
	Weight *float64 `gorm:"type:decimal(10,2)"` // Peso en kg
	Height *float64 `gorm:"type:decimal(10,2)"` // Alto en cm
	Width  *float64 `gorm:"type:decimal(10,2)"` // Ancho en cm
	Length *float64 `gorm:"type:decimal(10,2)"` // Largo en cm
	Boxes  *string  `gorm:"type:text"`          // Información de cajas/paquetes

	// ============================================
	// TIPO Y ESTADO DE LA ORDEN
	// ============================================
	OrderTypeID    *uint  `gorm:"index"`                                    // ID del tipo (delivery, pickup, etc.)
	OrderTypeName  string `gorm:"size:64"`                                  // Nombre del tipo (desnormalizado)
	Status         string `gorm:"size:64;not null;index;default:'pending'"` // Estado interno
	OriginalStatus string `gorm:"size:64"`                                  // Estado original de la plataforma

	// ============================================
	// INFORMACIÓN ADICIONAL
	// ============================================
	Notes    *string `gorm:"type:text"` // Notas de la orden
	Coupon   *string `gorm:"size:128"`  // Cupón aplicado
	Approved *bool   // Si está aprobada
	UserID   *uint   `gorm:"index"`    // Usuario que procesó
	UserName string  `gorm:"size:255"` // Nombre del usuario (desnormalizado)

	// ============================================
	// FACTURACIÓN
	// ============================================
	Invoiceable     bool    `gorm:"default:false"`  // Si es facturable
	InvoiceURL      *string `gorm:"size:512"`       // URL de la factura
	InvoiceID       *string `gorm:"size:128;index"` // ID de factura externa
	InvoiceProvider *string `gorm:"size:64"`        // Proveedor de facturación (ej: "siigo")

	// ============================================
	// DATOS ESTRUCTURADOS (JSONB)
	// ============================================
	Items    datatypes.JSON `gorm:"type:jsonb"` // Items de la orden
	Metadata datatypes.JSON `gorm:"type:jsonb"` // Metadata adicional de la plataforma

	// Detalles financieros adicionales (descuentos por item, promociones, etc.)
	FinancialDetails datatypes.JSON `gorm:"type:jsonb"`

	// Detalles de envío adicionales (proveedor, zona, etc.)
	ShippingDetails datatypes.JSON `gorm:"type:jsonb"`

	// Detalles de pago adicionales (transacción, referencia, etc.)
	PaymentDetails datatypes.JSON `gorm:"type:jsonb"`

	// Detalles de fulfillment adicionales (picking, packing, etc.)
	FulfillmentDetails datatypes.JSON `gorm:"type:jsonb"`

	// ============================================
	// TIMESTAMPS
	// ============================================
	OccurredAt time.Time `gorm:"index"` // Cuándo ocurrió en la plataforma
	ImportedAt time.Time `gorm:"index"` // Cuándo se importó a Probability

	// ============================================
	// RELACIONES (Solo para integridad referencial)
	// ============================================
	Business      *Business     `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Integration   Integration   `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PaymentMethod PaymentMethod `gorm:"foreignKey:PaymentMethodID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Relaciones con tablas relacionadas (Modelo Canónico)
	OrderItems      []OrderItem            `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Addresses       []Address              `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Payments        []Payment              `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Shipments       []Shipment             `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ChannelMetadata []OrderChannelMetadata `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName especifica el nombre de la tabla para Order
func (Order) TableName() string {
	return "orders"
}

// BeforeCreate genera el ID UUID y número interno antes de crear
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	// Generar UUID para el ID si no existe
	if o.ID == "" {
		o.ID = generateUUID()
	}

	// Generar número interno si no existe
	if o.InternalNumber == "" {
		o.InternalNumber = fmt.Sprintf("ORD-%d-%s",
			time.Now().Unix(),
			generateRandomString(6))
	}

	return nil
}

// ───────────────────────────────────────────
//
//	ORDER HISTORY - Historial de cambios
//
// ───────────────────────────────────────────

// OrderHistory registra cada cambio de estado de una orden
type OrderHistory struct {
	gorm.Model
	OrderID        string         `gorm:"type:varchar(36);not null;index"` // UUID de la orden
	PreviousStatus string         `gorm:"size:64"`
	NewStatus      string         `gorm:"size:64;not null"`
	ChangedBy      *uint          `gorm:"index"`      // ID del usuario que hizo el cambio
	ChangedByName  string         `gorm:"size:255"`   // Nombre del usuario (desnormalizado)
	Reason         *string        `gorm:"type:text"`  // Razón del cambio
	Metadata       datatypes.JSON `gorm:"type:jsonb"` // Metadata adicional del cambio

	// Relación
	Order Order `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName especifica el nombre de la tabla para OrderHistory
func (OrderHistory) TableName() string {
	return "order_history"
}

// ───────────────────────────────────────────
//
//	HELPER FUNCTIONS
//
// ───────────────────────────────────────────

// generateUUID genera un UUID v4 único
func generateUUID() string {
	// Generar 16 bytes aleatorios
	b := make([]byte, 16)
	rand.Read(b)

	// Configurar versión (4) y variante (RFC4122)
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant RFC4122

	// Formatear como UUID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// generateRandomString genera una cadena aleatoria de longitud n
func generateRandomString(n int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[mathrand.Intn(len(letters))]
	}
	return string(b)
}

// ───────────────────────────────────────────
//
//	ORDER ITEMS - Items/Productos de la orden
//
// ───────────────────────────────────────────

// OrderItem representa la relación entre una orden y un producto del catálogo
// Esta tabla solo guarda información específica de la orden (cantidad, precios de la venta, descuentos)
// Toda la información del producto (nombre, SKU, descripción, etc.) se obtiene de la tabla products
type OrderItem struct {
	gorm.Model

	// Relaciones
	OrderID   string  `gorm:"type:varchar(36);not null;index"` // FK a orders.id
	ProductID *string `gorm:"type:varchar(64);index"`          // FK a products.id (puede ser NULL si el producto se elimina)
	Order     Order   `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Product   Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	// ============================================
	// INFORMACIÓN ESPECÍFICA DE LA ORDEN
	// (No está en products, es específica de esta venta)
	// ============================================

	// Cantidad y precios de la venta (pueden diferir del precio actual del producto)
	Quantity   int     `gorm:"not null;default:1"`          // Cantidad comprada
	UnitPrice  float64 `gorm:"type:decimal(12,2);not null"` // Precio unitario al momento de la venta
	TotalPrice float64 `gorm:"type:decimal(12,2);not null"` // Precio total (quantity * unit_price)
	Currency   string  `gorm:"size:10;default:'USD'"`       // Moneda de la venta

	// Descuentos y ajustes aplicados en esta orden
	Discount float64  `gorm:"type:decimal(12,2);default:0"` // Descuento aplicado en esta orden
	Tax      float64  `gorm:"type:decimal(12,2);default:0"` // Impuesto de esta orden
	TaxRate  *float64 `gorm:"type:decimal(5,4)"`            // Tasa de impuesto aplicada (ej: 0.19 para 19%)

	// ============================================
	// INFORMACIÓN ESPECÍFICA DEL CONTEXTO
	// (Específica de esta orden, no del producto)
	// ============================================

	VariantID         *string        `gorm:"size:255"`   // ID de la variante en el sistema externo (si aplica)
	FulfillmentStatus *string        `gorm:"size:64"`    // Estado de fulfillment de este item en esta orden
	Metadata          datatypes.JSON `gorm:"type:jsonb"` // Metadata adicional del canal para este item específico
}

// TableName especifica el nombre de la tabla
func (OrderItem) TableName() string {
	return "order_items"
}

// ───────────────────────────────────────────
//
//	ADDRESSES - Direcciones de envío y facturación
//
// ───────────────────────────────────────────

// Address representa una dirección (puede ser de envío o facturación)
type Address struct {
	gorm.Model

	// Tipo de dirección
	Type string `gorm:"size:20;not null;index"` // "shipping" o "billing"

	// Relación con la orden
	OrderID string `gorm:"type:varchar(36);not null;index"` // UUID de la orden

	// Información del destinatario
	FirstName string `gorm:"size:128"` // Nombre
	LastName  string `gorm:"size:128"` // Apellido
	Company   string `gorm:"size:255"` // Empresa (opcional)
	Phone     string `gorm:"size:32"`  // Teléfono

	// Dirección física
	Street     string `gorm:"size:255;not null"` // Calle y número
	Street2    string `gorm:"size:255"`          // Segunda línea (apt, suite, etc.)
	City       string `gorm:"size:128;not null"` // Ciudad
	State      string `gorm:"size:128"`          // Estado/Provincia
	Country    string `gorm:"size:128;not null"` // País (código ISO)
	PostalCode string `gorm:"size:32"`           // Código postal

	// Coordenadas geográficas
	Latitude  *float64 `gorm:"type:decimal(10,8)"` // Latitud
	Longitude *float64 `gorm:"type:decimal(11,8)"` // Longitud

	// Información adicional
	Instructions *string        `gorm:"type:text"`     // Instrucciones de entrega
	IsDefault    bool           `gorm:"default:false"` // Si es dirección por defecto
	Metadata     datatypes.JSON `gorm:"type:jsonb"`    // Metadata adicional del canal

	// Relación
	Order Order `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName especifica el nombre de la tabla
func (Address) TableName() string {
	return "addresses"
}

// ───────────────────────────────────────────
//
//	PAYMENTS - Pagos de la orden
//
// ───────────────────────────────────────────

// Payment representa un pago asociado a una orden
type Payment struct {
	gorm.Model

	// Relación con la orden
	OrderID string `gorm:"type:varchar(36);not null;index"` // UUID de la orden

	// Método de pago
	PaymentMethodID uint `gorm:"not null;index"` // FK a payment_methods

	// Información financiera
	Amount       float64  `gorm:"type:decimal(12,2);not null"` // Monto pagado
	Currency     string   `gorm:"size:10;default:'USD'"`       // Moneda
	ExchangeRate *float64 `gorm:"type:decimal(10,4)"`          // Tasa de cambio (si aplica)

	// Estado del pago
	Status      string     `gorm:"size:64;not null;index"` // "pending", "completed", "failed", "refunded"
	PaidAt      *time.Time `gorm:"index"`                  // Cuándo se pagó
	ProcessedAt *time.Time // Cuándo se procesó

	// Identificadores del canal
	TransactionID    *string `gorm:"size:255;index"` // ID de transacción del canal
	PaymentReference *string `gorm:"size:255"`       // Referencia de pago
	Gateway          *string `gorm:"size:64"`        // Gateway utilizado (ej: "stripe", "paypal")

	// Información adicional
	RefundAmount  *float64       `gorm:"type:decimal(12,2)"` // Monto reembolsado
	RefundedAt    *time.Time     // Cuándo se reembolsó
	FailureReason *string        `gorm:"type:text"`  // Razón de fallo (si aplica)
	Metadata      datatypes.JSON `gorm:"type:jsonb"` // Metadata adicional del canal

	// Relaciones
	Order         Order         `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PaymentMethod PaymentMethod `gorm:"foreignKey:PaymentMethodID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName especifica el nombre de la tabla
func (Payment) TableName() string {
	return "payments"
}

// ───────────────────────────────────────────
//
//	ORDER CHANNEL METADATA - Datos crudos del canal
//
// ───────────────────────────────────────────

// OrderChannelMetadata almacena los datos crudos originales de la orden
// tal como los recibió el canal de venta (Shopify, Mercado Libre, etc.)
// Esta tabla es crucial para trazabilidad y flexibilidad
type OrderChannelMetadata struct {
	gorm.Model

	// Relación con la orden
	OrderID string `gorm:"type:varchar(36);not null;index"` // UUID de la orden

	// Identificación del canal
	ChannelSource string `gorm:"size:50;not null;index"` // "shopify", "mercado_libre", "paris", etc.
	IntegrationID uint   `gorm:"not null;index"`         // ID de la integración

	// Datos crudos
	RawData datatypes.JSON `gorm:"type:jsonb;not null"` // Payload completo de la API del canal

	// Metadata adicional
	Version     string     `gorm:"size:20"` // Versión del API/webhook recibido
	ReceivedAt  time.Time  `gorm:"index"`   // Cuándo se recibió el dato
	ProcessedAt *time.Time // Cuándo se procesó (si aplica)
	IsLatest    bool       `gorm:"default:true;index"` // Si es la versión más reciente

	// Información de sincronización
	LastSyncedAt *time.Time // Última vez que se sincronizó con el canal
	SyncStatus   string     `gorm:"size:64;default:'pending'"` // Estado de sincronización

	// Relaciones
	Order       Order       `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Integration Integration `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName especifica el nombre de la tabla
func (OrderChannelMetadata) TableName() string {
	return "order_channel_metadata"
}

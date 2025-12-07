package domain

import (
	"time"

	"gorm.io/datatypes"
)

// ───────────────────────────────────────────
//
//	DATOS FICTICIOS - Clientes y Productos
//
// ───────────────────────────────────────────

// FakeCustomer representa un cliente ficticio para generar órdenes
type FakeCustomer struct {
	Name  string
	Email string
	Phone string
	DNI   string
}

// FakeProduct representa un producto ficticio para generar órdenes
type FakeProduct struct {
	ID          string
	SKU         string
	Name        string
	Title       string
	Price       float64
	Weight      float64
	ImageURL    string
	ProductURL  string
	Description string
}

// FakeCustomers contiene una lista de clientes ficticios
var FakeCustomers = []FakeCustomer{
	{Name: "Juan Pérez", Email: "juan.perez@example.com", Phone: "+56912345678", DNI: "12345678-9"},
	{Name: "María González", Email: "maria.gonzalez@example.com", Phone: "+56987654321", DNI: "98765432-1"},
	{Name: "Carlos Rodríguez", Email: "carlos.rodriguez@example.com", Phone: "+56911223344", DNI: "11223344-5"},
	{Name: "Ana Martínez", Email: "ana.martinez@example.com", Phone: "+56955667788", DNI: "55667788-9"},
	{Name: "Luis Fernández", Email: "luis.fernandez@example.com", Phone: "+56999887766", DNI: "99887766-5"},
	{Name: "Laura Sánchez", Email: "laura.sanchez@example.com", Phone: "+56944332211", DNI: "44332211-7"},
	{Name: "Pedro López", Email: "pedro.lopez@example.com", Phone: "+56933445566", DNI: "33445566-3"},
	{Name: "Carmen Torres", Email: "carmen.torres@example.com", Phone: "+56977889900", DNI: "77889900-1"},
}

// FakeProducts contiene una lista de productos ficticios
var FakeProducts = []FakeProduct{
	{ID: "prod-001", SKU: "SKU-001", Name: "Camiseta Básica", Title: "Camiseta Básica - Talla M", Price: 19990, Weight: 0.2, ImageURL: "https://example.com/camiseta.jpg", ProductURL: "https://example.com/productos/camiseta", Description: "Camiseta de algodón 100%"},
	{ID: "prod-002", SKU: "SKU-002", Name: "Pantalón Jeans", Title: "Pantalón Jeans - Talla 32", Price: 49990, Weight: 0.5, ImageURL: "https://example.com/jeans.jpg", ProductURL: "https://example.com/productos/jeans", Description: "Pantalón jeans clásico"},
	{ID: "prod-003", SKU: "SKU-003", Name: "Zapatillas Deportivas", Title: "Zapatillas Deportivas - Talla 42", Price: 79990, Weight: 0.8, ImageURL: "https://example.com/zapatillas.jpg", ProductURL: "https://example.com/productos/zapatillas", Description: "Zapatillas para running"},
	{ID: "prod-004", SKU: "SKU-004", Name: "Chaqueta Impermeable", Title: "Chaqueta Impermeable - Talla L", Price: 89990, Weight: 0.6, ImageURL: "https://example.com/chaqueta.jpg", ProductURL: "https://example.com/productos/chaqueta", Description: "Chaqueta resistente al agua"},
	{ID: "prod-005", SKU: "SKU-005", Name: "Gorra Deportiva", Title: "Gorra Deportiva - One Size", Price: 14990, Weight: 0.1, ImageURL: "https://example.com/gorra.jpg", ProductURL: "https://example.com/productos/gorra", Description: "Gorra ajustable"},
	{ID: "prod-006", SKU: "SKU-006", Name: "Mochila Urbana", Title: "Mochila Urbana - 30L", Price: 59990, Weight: 0.7, ImageURL: "https://example.com/mochila.jpg", ProductURL: "https://example.com/productos/mochila", Description: "Mochila resistente y espaciosa"},
	{ID: "prod-007", SKU: "SKU-007", Name: "Reloj Inteligente", Title: "Reloj Inteligente - Smart Watch", Price: 149990, Weight: 0.05, ImageURL: "https://example.com/reloj.jpg", ProductURL: "https://example.com/productos/reloj", Description: "Reloj inteligente con GPS"},
	{ID: "prod-008", SKU: "SKU-008", Name: "Auriculares Inalámbricos", Title: "Auriculares Inalámbricos - Bluetooth", Price: 69990, Weight: 0.15, ImageURL: "https://example.com/auriculares.jpg", ProductURL: "https://example.com/productos/auriculares", Description: "Auriculares con cancelación de ruido"},
	{ID: "prod-009", SKU: "SKU-009", Name: "Cargador Portátil", Title: "Cargador Portátil - 20000mAh", Price: 29990, Weight: 0.4, ImageURL: "https://example.com/cargador.jpg", ProductURL: "https://example.com/productos/cargador", Description: "Batería externa de alta capacidad"},
	{ID: "prod-010", SKU: "SKU-010", Name: "Fundas para Celular", Title: "Fundas para Celular - Pack x3", Price: 9990, Weight: 0.05, ImageURL: "https://example.com/fundas.jpg", ProductURL: "https://example.com/productos/fundas", Description: "Pack de 3 fundas protectoras"},
}

// FakeAddresses contiene direcciones ficticias
var FakeAddresses = []struct {
	Street     string
	Street2    string
	City       string
	State      string
	Country    string
	PostalCode string
	Latitude   float64
	Longitude  float64
}{
	{Street: "Av. Providencia 1234", Street2: "Depto 5A", City: "Santiago", State: "Región Metropolitana", Country: "CL", PostalCode: "7500000", Latitude: -33.4489, Longitude: -70.6693},
	{Street: "Calle Los Rosales 567", Street2: "", City: "Valparaíso", State: "Valparaíso", Country: "CL", PostalCode: "2340000", Latitude: -33.0472, Longitude: -71.6127},
	{Street: "Av. Libertador 890", Street2: "Casa 12", City: "Concepción", State: "Biobío", Country: "CL", PostalCode: "4030000", Latitude: -36.8201, Longitude: -73.0444},
	{Street: "Pasaje Las Flores 234", Street2: "", City: "La Serena", State: "Coquimbo", Country: "CL", PostalCode: "1700000", Latitude: -29.9027, Longitude: -71.2519},
	{Street: "Av. Costanera 456", Street2: "Torre B, Piso 8", City: "Viña del Mar", State: "Valparaíso", Country: "CL", PostalCode: "2520000", Latitude: -33.0246, Longitude: -71.5518},
}

// FakePaymentMethods contiene métodos de pago ficticios
var FakePaymentMethods = []struct {
	ID   uint
	Name string
}{
	{ID: 1, Name: "Tarjeta de Crédito"},
	{ID: 2, Name: "Tarjeta de Débito"},
	{ID: 3, Name: "Transferencia Bancaria"},
	{ID: 4, Name: "Efectivo"},
	{ID: 5, Name: "Mercado Pago"},
}

// FakeCarriers contiene transportistas ficticios
var FakeCarriers = []struct {
	Name string
	Code string
}{
	{Name: "Chilexpress", Code: "CHX"},
	{Name: "Correos de Chile", Code: "COR"},
	{Name: "Starken", Code: "STK"},
	{Name: "Blue Express", Code: "BLX"},
}

// GenerateOrderRequest representa la solicitud para generar órdenes
type GenerateOrderRequest struct {
	Count           int    `json:"count" binding:"required,min=1,max=100"` // Cantidad de órdenes a generar
	IntegrationID   uint   `json:"integration_id" binding:"required"`      // ID de la integración
	BusinessID      *uint  `json:"business_id"`                            // ID del negocio (opcional)
	Platform        string `json:"platform" binding:"max=50"`              // Plataforma (default: "test")
	Status          string `json:"status" binding:"max=64"`                // Estado inicial (default: "pending")
	IncludePayment  bool   `json:"include_payment"`                        // Si incluir información de pago
	IncludeShipment bool   `json:"include_shipment"`                       // Si incluir información de envío
}

// GenerateOrderResponse representa la respuesta de generación
type GenerateOrderResponse struct {
	Generated int      `json:"generated"` // Cantidad de órdenes generadas
	Published int      `json:"published"` // Cantidad de órdenes publicadas
	Failed    int      `json:"failed"`    // Cantidad de órdenes que fallaron
	OrderIDs  []string `json:"order_ids"` // IDs externos de las órdenes generadas
}

// ───────────────────────────────────────────
//
//	CANONICAL ORDER DTO - Duplicado para uso en test
//	(No podemos importar internal desde otro módulo)
//
// ───────────────────────────────────────────

// CanonicalOrderDTO representa la estructura canónica que todas las integraciones
// deben enviar después de mapear sus datos específicos
type CanonicalOrderDTO struct {
	// Identificadores de integración
	BusinessID      *uint  `json:"business_id"`
	IntegrationID   uint   `json:"integration_id" binding:"required"`
	IntegrationType string `json:"integration_type" binding:"required,max=50"`

	// Identificadores de la orden
	Platform       string `json:"platform" binding:"required,max=50"`
	ExternalID     string `json:"external_id" binding:"required,max=255"`
	OrderNumber    string `json:"order_number" binding:"max=128"`
	InternalNumber string `json:"internal_number" binding:"max=128"`

	// Información financiera
	Subtotal     float64  `json:"subtotal" binding:"required,min=0"`
	Tax          float64  `json:"tax" binding:"min=0"`
	Discount     float64  `json:"discount" binding:"min=0"`
	ShippingCost float64  `json:"shipping_cost" binding:"min=0"`
	TotalAmount  float64  `json:"total_amount" binding:"required,min=0"`
	Currency     string   `json:"currency" binding:"max=10"`
	CodTotal     *float64 `json:"cod_total"`

	// Información del cliente
	CustomerID    *uint  `json:"customer_id"`
	CustomerName  string `json:"customer_name" binding:"max=255"`
	CustomerEmail string `json:"customer_email" binding:"max=255"`
	CustomerPhone string `json:"customer_phone" binding:"max=32"`
	CustomerDNI   string `json:"customer_dni" binding:"max=64"`

	// Tipo y estado
	OrderTypeID    *uint  `json:"order_type_id"`
	OrderTypeName  string `json:"order_type_name" binding:"max=64"`
	Status         string `json:"status" binding:"max=64"`
	OriginalStatus string `json:"original_status" binding:"max=64"`

	// Información adicional
	Notes    *string `json:"notes"`
	Coupon   *string `json:"coupon"`
	Approved *bool   `json:"approved"`
	UserID   *uint   `json:"user_id"`
	UserName string  `json:"user_name" binding:"max=255"`

	// Facturación
	Invoiceable     bool    `json:"invoiceable"`
	InvoiceURL      *string `json:"invoice_url"`
	InvoiceID       *string `json:"invoice_id"`
	InvoiceProvider *string `json:"invoice_provider"`

	// Timestamps
	OccurredAt time.Time `json:"occurred_at"`
	ImportedAt time.Time `json:"imported_at"`

	// Datos estructurados (JSONB) - Para compatibilidad
	Items              datatypes.JSON `json:"items,omitempty"`
	Metadata           datatypes.JSON `json:"metadata,omitempty"`
	FinancialDetails   datatypes.JSON `json:"financial_details,omitempty"`
	ShippingDetails    datatypes.JSON `json:"shipping_details,omitempty"`
	PaymentDetails     datatypes.JSON `json:"payment_details,omitempty"`
	FulfillmentDetails datatypes.JSON `json:"fulfillment_details,omitempty"`

	// ============================================
	// TABLAS RELACIONADAS
	// ============================================

	// Items de la orden
	OrderItems []CanonicalOrderItemDTO `json:"order_items" binding:"dive"`

	// Direcciones
	Addresses []CanonicalAddressDTO `json:"addresses" binding:"dive"`

	// Pagos
	Payments []CanonicalPaymentDTO `json:"payments" binding:"dive"`

	// Envíos
	Shipments []CanonicalShipmentDTO `json:"shipments" binding:"dive"`

	// Metadata del canal (datos crudos)
	ChannelMetadata *CanonicalChannelMetadataDTO `json:"channel_metadata"`
}

// CanonicalOrderItemDTO representa un item/producto de la orden
type CanonicalOrderItemDTO struct {
	ProductID    *string        `json:"product_id"`
	ProductSKU   string         `json:"product_sku" binding:"required,max=128"`
	ProductName  string         `json:"product_name" binding:"required,max=255"`
	ProductTitle string         `json:"product_title" binding:"max=255"`
	VariantID    *string        `json:"variant_id"`
	Quantity     int            `json:"quantity" binding:"required,min=1"`
	UnitPrice    float64        `json:"unit_price" binding:"required,min=0"`
	TotalPrice   float64        `json:"total_price" binding:"required,min=0"`
	Currency     string         `json:"currency" binding:"max=10"`
	Discount     float64        `json:"discount" binding:"min=0"`
	Tax          float64        `json:"tax" binding:"min=0"`
	TaxRate      *float64       `json:"tax_rate"`
	ImageURL     *string        `json:"image_url"`
	ProductURL   *string        `json:"product_url"`
	Weight       *float64       `json:"weight"`
	Metadata     datatypes.JSON `json:"metadata,omitempty"`
}

// CanonicalAddressDTO representa una dirección (envío o facturación)
type CanonicalAddressDTO struct {
	Type         string         `json:"type" binding:"required,oneof=shipping billing"` // "shipping" o "billing"
	FirstName    string         `json:"first_name" binding:"max=128"`
	LastName     string         `json:"last_name" binding:"max=128"`
	Company      string         `json:"company" binding:"max=255"`
	Phone        string         `json:"phone" binding:"max=32"`
	Street       string         `json:"street" binding:"required,max=255"`
	Street2      string         `json:"street2" binding:"max=255"`
	City         string         `json:"city" binding:"required,max=128"`
	State        string         `json:"state" binding:"max=128"`
	Country      string         `json:"country" binding:"required,max=128"`
	PostalCode   string         `json:"postal_code" binding:"max=32"`
	Latitude     *float64       `json:"latitude"`
	Longitude    *float64       `json:"longitude"`
	Instructions *string        `json:"instructions"`
	Metadata     datatypes.JSON `json:"metadata,omitempty"`
}

// CanonicalPaymentDTO representa un pago de la orden
type CanonicalPaymentDTO struct {
	PaymentMethodID  uint           `json:"payment_method_id" binding:"required"`
	Amount           float64        `json:"amount" binding:"required,min=0"`
	Currency         string         `json:"currency" binding:"max=10"`
	ExchangeRate     *float64       `json:"exchange_rate"`
	Status           string         `json:"status" binding:"required,oneof=pending completed failed refunded"`
	PaidAt           *time.Time     `json:"paid_at"`
	ProcessedAt      *time.Time     `json:"processed_at"`
	TransactionID    *string        `json:"transaction_id"`
	PaymentReference *string        `json:"payment_reference"`
	Gateway          *string        `json:"gateway"`
	RefundAmount     *float64       `json:"refund_amount"`
	RefundedAt       *time.Time     `json:"refunded_at"`
	FailureReason    *string        `json:"failure_reason"`
	Metadata         datatypes.JSON `json:"metadata,omitempty"`
}

// CanonicalShipmentDTO representa un envío de la orden
type CanonicalShipmentDTO struct {
	TrackingNumber    *string        `json:"tracking_number"`
	TrackingURL       *string        `json:"tracking_url"`
	Carrier           *string        `json:"carrier"`
	CarrierCode       *string        `json:"carrier_code"`
	GuideID           *string        `json:"guide_id"`
	GuideURL          *string        `json:"guide_url"`
	Status            string         `json:"status" binding:"oneof=pending in_transit delivered failed"`
	ShippedAt         *time.Time     `json:"shipped_at"`
	DeliveredAt       *time.Time     `json:"delivered_at"`
	ShippingAddressID *uint          `json:"shipping_address_id"`
	ShippingCost      *float64       `json:"shipping_cost"`
	InsuranceCost     *float64       `json:"insurance_cost"`
	TotalCost         *float64       `json:"total_cost"`
	Weight            *float64       `json:"weight"`
	Height            *float64       `json:"height"`
	Width             *float64       `json:"width"`
	Length            *float64       `json:"length"`
	WarehouseID       *uint          `json:"warehouse_id"`
	WarehouseName     string         `json:"warehouse_name" binding:"max=128"`
	DriverID          *uint          `json:"driver_id"`
	DriverName        string         `json:"driver_name" binding:"max=255"`
	IsLastMile        bool           `json:"is_last_mile"`
	EstimatedDelivery *time.Time     `json:"estimated_delivery"`
	DeliveryNotes     *string        `json:"delivery_notes"`
	Metadata          datatypes.JSON `json:"metadata,omitempty"`
}

// CanonicalChannelMetadataDTO representa los datos crudos del canal
type CanonicalChannelMetadataDTO struct {
	ChannelSource string         `json:"channel_source" binding:"required,max=50"`
	RawData       datatypes.JSON `json:"raw_data" binding:"required"`
	Version       string         `json:"version" binding:"max=20"`
	ReceivedAt    time.Time      `json:"received_at"`
	ProcessedAt   *time.Time     `json:"processed_at"`
	IsLatest      bool           `json:"is_latest"`
	LastSyncedAt  *time.Time     `json:"last_synced_at"`
	SyncStatus    string         `json:"sync_status" binding:"max=64"`
}

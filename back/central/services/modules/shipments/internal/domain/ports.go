package domain

import (
	"context"
	"time"
)

//
//	REPOSITORY INTERFACE
//

// IRepository define todos los métodos de repositorio del módulo shipments
type IRepository interface {
	// CRUD Operations
	CreateShipment(ctx context.Context, shipment *Shipment) error
	GetShipmentByID(ctx context.Context, id uint) (*Shipment, error)
	GetShipmentByTrackingNumber(ctx context.Context, trackingNumber string) (*Shipment, error)
	GetShipmentsByOrderID(ctx context.Context, orderID string) ([]Shipment, error)
	ListShipments(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]Shipment, int64, error)
	UpdateShipment(ctx context.Context, shipment *Shipment) error
	DeleteShipment(ctx context.Context, id uint) error

	// Validation
	ShipmentExists(ctx context.Context, orderID string, trackingNumber string) (bool, error)

	// Carrier Resolution (replicated query — module isolation)
	GetActiveShippingCarrier(ctx context.Context, businessID uint) (*CarrierInfo, error)

	// Business name lookup (for descriptive error messages)
	GetBusinessName(ctx context.Context, businessID uint) (string, error)

	// Business ID Resolution (replicated queries — module isolation)
	GetOrderBusinessID(ctx context.Context, orderUUID string) (uint, error)
	GetShipmentBusinessIDByTracking(ctx context.Context, trackingNumber string) (uint, error)
	GetShipmentBusinessIDByID(ctx context.Context, shipmentID uint) (uint, error)

	// Updates guide_link, tracking_number, and carrier on the orders table after guide generation.
	UpdateOrderGuideLink(ctx context.Context, orderID string, guideLink string, trackingNumber string, carrier string) error

	UpdateOrderStatusByOrderID(ctx context.Context, orderID string, status string) error
	ClearOrderGuideData(ctx context.Context, orderID string) error

	// EnsureAllBusinessesActive sets all existing businesses to 'paid' status as a one-time migration
	EnsureAllBusinessesActive(ctx context.Context) error

	// GetOrderIntegrationID retrieves the integration_id for an order (replicated query — module isolation)
	GetOrderIntegrationID(ctx context.Context, orderUUID string) (uint, error)

	ListShipmentsForSync(ctx context.Context, filter SyncShipmentsFilter) ([]SyncShipmentRow, error)
	GetBusinessActiveIntegration(ctx context.Context, businessID uint, providerCode string) (uint, string, error)

	DebitWalletForGuide(ctx context.Context, businessID uint, amount float64, trackingNumber string) error

	// Origin Addresses
	CreateOriginAddress(ctx context.Context, address *OriginAddress) error
	GetOriginAddressByID(ctx context.Context, id uint) (*OriginAddress, error)
	ListOriginAddressesByBusiness(ctx context.Context, businessID uint) ([]OriginAddress, error)
	GetDefaultOriginAddress(ctx context.Context, businessID uint) (*OriginAddress, error)
	UpdateOriginAddress(ctx context.Context, address *OriginAddress) error
	DeleteOriginAddress(ctx context.Context, id uint) error
	SetDefaultOriginAddress(ctx context.Context, businessID, addressID uint) error

	ListCODShipments(ctx context.Context, filter CODFilter) ([]Shipment, int64, error)
	GetOrderCODInfo(ctx context.Context, orderID string) (*OrderCODInfo, error)
	MarkOrderPaidCOD(ctx context.Context, orderID string, amount float64, paymentMethodID uint, notes string) error

	GetOrderPublicTrackingByNumber(ctx context.Context, orderNumber string) (*OrderPublicTracking, error)
}

type OrderPublicTracking struct {
	ID                 string
	OrderNumber        string
	BusinessID         uint
	BusinessName       string
	Status             string
	IsPaid             bool
	TotalAmount        float64
	CodTotal           *float64
	Currency           string
	CustomerName       string
	CustomerPhone      string
	ShippingStreet     string
	ShippingCity       string
	ShippingState      string
	ShippingPostalCode string
	CreatedAt          time.Time
	OccurredAt         *time.Time
}

type OrderCODInfo struct {
	OrderID            string
	BusinessID         uint
	CodTotal           *float64
	TotalAmount        float64
	Currency           string
	IsPaid             bool
	PaidAt             *time.Time
	PaymentMethodID    uint
	PaymentMethodCode  string
}

//
//	CARRIER RESOLUTION
//

// CarrierInfo holds the resolved carrier for a business
type CarrierInfo struct {
	IntegrationID     uint   // integrations.id
	IntegrationTypeID uint   // integration_types.id (e.g. 12 = EnvioClick)
	ProviderCode      string // integration_types.code (e.g. "envioclick")
	BaseURL           string // integration_types.base_url (production URL)
	IsTesting         bool   // integrations.is_testing (use sandbox URL)
	BaseURLTest       string // integrations.config->>'base_url_test' (sandbox URL)
}

// ICarrierResolver resolves which shipping carrier a business has active
type ICarrierResolver interface {
	GetActiveShippingCarrier(ctx context.Context, businessID uint) (*CarrierInfo, error)
}

//
//	TRANSPORT REQUEST PUBLISHER
//

// TransportRequestMessage is the message published to the transport queue
type TransportRequestMessage struct {
	ShipmentID        *uint                  `json:"shipment_id,omitempty"`
	Provider          string                 `json:"provider"`
	IntegrationTypeID uint                   `json:"integration_type_id"`
	Operation         string                 `json:"operation"`
	CorrelationID     string                 `json:"correlation_id"`
	BusinessID        uint                   `json:"business_id"`
	IntegrationID     uint                   `json:"integration_id"`
	BaseURL           string                 `json:"base_url,omitempty"`
	IsTest            bool                   `json:"is_test,omitempty"`
	Timestamp         time.Time              `json:"timestamp"`
	Payload           map[string]interface{} `json:"payload"`
}

// ITransportRequestPublisher defines the contract for publishing transport requests
type ITransportRequestPublisher interface {
	PublishTransportRequest(ctx context.Context, request *TransportRequestMessage) error
}

//
//	SSE PUBLISHER
//

// GuideNotificationData holds enriched data for guide_generated events (used for WhatsApp notifications)
type GuideNotificationData struct {
	CustomerName  string
	CustomerPhone string
	OrderNumber   string
	BusinessName  string
	IntegrationID uint
	CodTotal      *float64
	TrackingURL   string
}

// IShipmentSSEPublisher defines the contract for publishing shipment SSE events via Redis
type IShipmentSSEPublisher interface {
	PublishQuoteReceived(ctx context.Context, businessID uint, correlationID string, data map[string]interface{})
	PublishQuoteFailed(ctx context.Context, businessID uint, correlationID string, errorMsg string)
	PublishGuideGenerated(ctx context.Context, businessID uint, shipmentID uint, correlationID string, trackingNumber string, labelURL string, carrier string, notification *GuideNotificationData)
	PublishGuideFailed(ctx context.Context, businessID uint, shipmentID uint, correlationID string, errorMsg string)
	PublishTrackingUpdated(ctx context.Context, businessID uint, correlationID string, data map[string]interface{})
	PublishTrackingFailed(ctx context.Context, businessID uint, correlationID string, errorMsg string)
	PublishShipmentCancelled(ctx context.Context, businessID uint, shipmentID uint)
	PublishCancelFailed(ctx context.Context, businessID uint, shipmentID uint, correlationID string, errorMsg string)
}

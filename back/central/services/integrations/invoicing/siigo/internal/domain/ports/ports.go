package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

type ISiigoClient interface {
	TestAuthentication(ctx context.Context, username, accessKey, accountID, partnerID, baseURL string) error

	CreateInvoice(ctx context.Context, req *dtos.CreateInvoiceRequest) (*dtos.CreateInvoiceResult, error)

	GetCustomerByIdentification(ctx context.Context, credentials dtos.Credentials, identification string) (*dtos.CustomerResult, error)

	CreateCustomer(ctx context.Context, credentials dtos.Credentials, req *dtos.CreateCustomerRequest) (*dtos.CustomerResult, error)

	ListInvoices(ctx context.Context, credentials dtos.Credentials, params dtos.ListInvoicesParams) (*dtos.ListInvoicesResult, error)

	GetInvoiceByID(ctx context.Context, credentials dtos.Credentials, invoiceID string) (*dtos.InvoiceDetail, error)

	GetStampErrors(ctx context.Context, credentials dtos.Credentials, invoiceID string) ([]dtos.StampError, error)

	AnnulInvoice(ctx context.Context, credentials dtos.Credentials, invoiceID string) (*dtos.AnnulInvoiceResult, error)

	ListProducts(ctx context.Context, credentials dtos.Credentials, page, pageSize int) ([]dtos.ProductItem, error)

	ListPaymentTypes(ctx context.Context, credentials dtos.Credentials, documentType string) ([]dtos.PaymentTypeItem, error)

	CreateCashReceipt(ctx context.Context, req *dtos.CreateCashReceiptRequest) (*dtos.CreateCashReceiptResult, error)

	CreateJournal(ctx context.Context, req *dtos.CreateJournalRequest) (*dtos.CreateJournalResult, error)
}

type IInvoiceUseCase interface {
	ProcessOrderForInvoicing(ctx context.Context, event *OrderEventMessage) error

	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
}

type OrderEventMessage struct {
	EventID       string                 `json:"event_id"`
	EventType     string                 `json:"event_type"`
	OrderID       string                 `json:"order_id"`
	BusinessID    *uint                  `json:"business_id"`
	IntegrationID *uint                  `json:"integration_id"`
	Timestamp     time.Time              `json:"timestamp"`
	Order         *OrderSnapshot         `json:"order"`
	Changes       map[string]interface{} `json:"changes,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

type OrderSnapshot struct {
	ID              string              `json:"id"`
	OrderNumber     string              `json:"order_number"`
	InternalNumber  string              `json:"internal_number"`
	ExternalID      string              `json:"external_id"`
	TotalAmount     float64             `json:"total_amount"`
	Currency        string              `json:"currency"`
	PaymentMethodID uint                `json:"payment_method_id"`
	PaymentStatusID *uint               `json:"payment_status_id,omitempty"`
	Subtotal        float64             `json:"subtotal"`
	Tax             float64             `json:"tax"`
	Discount        float64             `json:"discount"`
	ShippingCost    float64             `json:"shipping_cost"`
	CustomerName    string              `json:"customer_name"`
	CustomerEmail   string              `json:"customer_email,omitempty"`
	CustomerPhone   string              `json:"customer_phone,omitempty"`
	CustomerDNI     string              `json:"customer_dni,omitempty"`
	Platform        string              `json:"platform"`
	IntegrationID   uint                `json:"integration_id"`
	Items           []OrderItemSnapshot `json:"items,omitempty"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

type OrderItemSnapshot struct {
	ProductID  *string  `json:"product_id,omitempty"`
	SKU        string   `json:"sku"`
	Name       string   `json:"name"`
	Quantity   int      `json:"quantity"`
	UnitPrice  float64  `json:"unit_price"`
	TotalPrice float64  `json:"total_price"`
	Tax        float64  `json:"tax"`
	TaxRate    *float64 `json:"tax_rate,omitempty"`
	Discount   float64  `json:"discount"`
}

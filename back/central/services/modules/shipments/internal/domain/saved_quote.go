package domain

import "time"

const (
	QuoteSourceShopify     = "shopify"
	QuoteSourcePanel       = "panel"
	QuoteSourceWooCommerce = "woocommerce"

	QuoteStatusCreated        = "created"
	QuoteStatusAssociated     = "associated"
	QuoteStatusGuideGenerated = "guide_generated"
	QuoteStatusExpired        = "expired"
	QuoteStatusFailed         = "failed"
)

type SavedQuote struct {
	ID                  uint
	BusinessID          uint
	IntegrationID       uint
	Source              string
	CorrelationID       string
	OrderUUID           *string
	ExternalOrderRef    string
	RequestPayload      map[string]interface{}
	Rates               []map[string]interface{}
	SelectedCarrier     string
	SelectedServiceCode string
	SelectedIDRate      *int64
	Status              string
	ExpiresAt           *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type SaveQuoteInput struct {
	BusinessID       uint
	IntegrationID    uint
	Source           string
	CorrelationID    string
	OrderUUID        *string
	ExternalOrderRef string
	RequestPayload   map[string]interface{}
	Rates            []map[string]interface{}
	TTL              time.Duration
}

type SavedQuoteFilter struct {
	BusinessID uint
	Source     string
	Status     string
	OrderUUID  string
	Page       int
	PageSize   int
}

type SavedQuoteResponse struct {
	ID                  uint                     `json:"id"`
	BusinessID          uint                     `json:"business_id"`
	IntegrationID       uint                     `json:"integration_id"`
	Source              string                   `json:"source"`
	CorrelationID       string                   `json:"correlation_id,omitempty"`
	OrderUUID           *string                  `json:"order_uuid,omitempty"`
	ExternalOrderRef    string                   `json:"external_order_ref,omitempty"`
	Rates               []map[string]interface{} `json:"rates"`
	SelectedCarrier     string                   `json:"selected_carrier,omitempty"`
	SelectedServiceCode string                   `json:"selected_service_code,omitempty"`
	Status              string                   `json:"status"`
	ExpiresAt           *time.Time               `json:"expires_at,omitempty"`
	CreatedAt           time.Time                `json:"created_at"`
}

type SavedQuotesListResponse struct {
	Data       []SavedQuoteResponse `json:"data"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

type OrderSelectedShipping struct {
	Code   string
	Title  string
	Source string
	Price  string
}

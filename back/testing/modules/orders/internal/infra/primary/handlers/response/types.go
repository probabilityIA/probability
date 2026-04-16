package response

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	SKU      string  `json:"sku"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

type Integration struct {
	ID                  uint   `json:"id"`
	Name                string `json:"name"`
	Code                string `json:"code"`
	Category            string `json:"category"`
	CategoryID          uint   `json:"category_id"`
	IntegrationTypeID   uint   `json:"integration_type_id"`
	IntegrationTypeCode string `json:"integration_type_code"`
}

type PaymentMethod struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type OrderStatus struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type ReferenceData struct {
	Products       []Product           `json:"products"`
	Integrations   []Integration       `json:"integrations"`
	PaymentMethods []PaymentMethod     `json:"payment_methods"`
	OrderStatuses  []OrderStatus       `json:"order_statuses"`
	WebhookTopics  map[string][]string `json:"webhook_topics"`
}

type WebhookPayload struct {
	URL        string            `json:"url"`
	Method     string            `json:"method"`
	Headers    map[string]string `json:"headers"`
	Body       interface{}       `json:"body"`
	RawBody    string            `json:"raw_body,omitempty"`
	HMACSecret string            `json:"hmac_secret,omitempty"`
}

type OrderError struct {
	Index   int    `json:"index"`
	Message string `json:"message"`
}

type GenerateResult struct {
	Total    int              `json:"total"`
	Payloads []WebhookPayload `json:"payloads"`
	Errors   []OrderError     `json:"errors"`
}

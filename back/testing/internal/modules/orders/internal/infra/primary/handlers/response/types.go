package response

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	SKU      string  `json:"sku"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

type Integration struct {
	ID                uint   `json:"id"`
	Name              string `json:"name"`
	Code              string `json:"code"`
	Category          string `json:"category"`
	IntegrationTypeID uint   `json:"integration_type_id"`
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
	Products       []Product       `json:"products"`
	Integrations   []Integration   `json:"integrations"`
	PaymentMethods []PaymentMethod `json:"payment_methods"`
	OrderStatuses  []OrderStatus   `json:"order_statuses"`
}

type CreatedOrder struct {
	ID           string  `json:"id"`
	OrderNumber  string  `json:"order_number"`
	Total        float64 `json:"total"`
	CustomerName string  `json:"customer_name"`
}

type OrderError struct {
	Index   int    `json:"index"`
	Message string `json:"message"`
}

type APIRequest struct {
	Method string                 `json:"method"`
	URL    string                 `json:"url"`
	Body   map[string]interface{} `json:"body"`
}

type APIResponse struct {
	StatusCode int    `json:"status_code"`
	Body       string `json:"body"`
}

type APICallLog struct {
	Index      int         `json:"index"`
	Success    bool        `json:"success"`
	Timestamp  string      `json:"timestamp"`
	DurationMs int64       `json:"duration_ms"`
	Request    APIRequest  `json:"request"`
	Response   APIResponse `json:"response"`
}

type GenerateResult struct {
	Total   int            `json:"total"`
	Created int            `json:"created"`
	Failed  int            `json:"failed"`
	Orders  []CreatedOrder `json:"orders"`
	Errors  []OrderError   `json:"errors"`
	APILogs []APICallLog   `json:"api_logs"`
}

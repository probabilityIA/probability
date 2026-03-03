package entities

type Product struct {
	ID       string
	Name     string
	SKU      string
	Price    float64
	Currency string
}

type Integration struct {
	ID             uint
	Name           string
	Code           string
	Category       string
	IntegrationTypeID uint
}

type PaymentMethod struct {
	ID   uint
	Code string
	Name string
}

type OrderStatus struct {
	ID   uint
	Code string
	Name string
}

type ReferenceData struct {
	Products       []Product
	Integrations   []Integration
	PaymentMethods []PaymentMethod
	OrderStatuses  []OrderStatus
}

type APIRequest struct {
	Method string
	URL    string
	Body   map[string]interface{}
}

type APIResponse struct {
	StatusCode int
	Body       string
}

type APICallLog struct {
	Index      int
	Success    bool
	Timestamp  string
	DurationMs int64
	Request    APIRequest
	Response   APIResponse
}

type GenerateResult struct {
	Total     int
	Created   int
	Failed    int
	Orders    []CreatedOrder
	Errors    []OrderError
	APILogs   []APICallLog
}

type CreatedOrder struct {
	ID          string
	OrderNumber string
	Total       float64
	CustomerName string
}

type OrderError struct {
	Index   int
	Message string
}

package request

// SiigoInvoice representa el body de una solicitud de creación de factura en Siigo
// Endpoint: POST /v1/invoices
type SiigoInvoice struct {
	Document SiigoDocument    `json:"document"`
	Customer SiigoCustomerRef `json:"customer"`
	Date     string           `json:"date"` // "YYYY-MM-DD"
	Currency SiigoCurrency    `json:"currency,omitempty"`
	Items    []SiigoItem      `json:"items"`
	Payments []SiigoPayment   `json:"payments,omitempty"`
	Observations string       `json:"observations,omitempty"`
}

// SiigoDocument referencia al tipo de documento (FV = Factura de Venta)
type SiigoDocument struct {
	ID int `json:"id"`
}

// SiigoCustomerRef referencia al cliente en Siigo
type SiigoCustomerRef struct {
	PersonType     string          `json:"person_type"`              // "Person" o "Company"
	IDType         SiigoIDType     `json:"id_type"`
	Identification string          `json:"identification"`
	Name           []string        `json:"name"`                     // [first_name, last_name] o [company_name]
	Address        *SiigoAddress   `json:"address,omitempty"`
	Phones         []SiigoPhone    `json:"phones,omitempty"`
	Contacts       []SiigoContact  `json:"contacts,omitempty"`
}

// SiigoIDType tipo de documento de identidad
type SiigoIDType struct {
	Code string `json:"code"` // "13"=CC, "31"=NIT, "22"=CE, etc.
}

// SiigoAddress dirección del cliente
type SiigoAddress struct {
	Address string `json:"address"`
}

// SiigoPhone teléfono del cliente
type SiigoPhone struct {
	Indicative string `json:"indicative,omitempty"` // "+57"
	Number     string `json:"number"`
	Extension  string `json:"extension,omitempty"`
}

// SiigoContact contacto del cliente (para email)
type SiigoContact struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email"`
	SendElectronicInvoice bool `json:"send_electronic_invoice"`
}

// SiigoCurrency moneda de la factura
type SiigoCurrency struct {
	Code         string  `json:"code"`          // "COP", "USD", etc.
	ExchangeRate float64 `json:"exchange_rate,omitempty"`
}

// SiigoItem item de la factura
type SiigoItem struct {
	Code     SiigoProductCode `json:"code"`
	Description string        `json:"description"`
	Quantity float64          `json:"quantity"`
	Price    float64          `json:"price"`
	Discount float64          `json:"discount,omitempty"` // Porcentaje 0-100
	Taxes    []SiigoTax       `json:"taxes,omitempty"`
}

// SiigoProductCode código del producto en Siigo
type SiigoProductCode struct {
	Code string `json:"code"`
}

// SiigoTax impuesto del item
type SiigoTax struct {
	ID int `json:"id"` // ID del impuesto en Siigo (ej: IVA 19% = ID específico por config)
}

// SiigoPayment información de pago
type SiigoPayment struct {
	ID       int     `json:"id"`   // ID del método de pago en Siigo
	Value    float64 `json:"value"`
	DueDate  string  `json:"due_date,omitempty"` // "YYYY-MM-DD"
}

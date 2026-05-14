package request

type SiigoInvoice struct {
	Document     SiigoDocument    `json:"document"`
	Date         string           `json:"date"`
	Customer     SiigoCustomerRef `json:"customer"`
	Seller       int              `json:"seller"`
	CostCenter   int              `json:"cost_center,omitempty"`
	Currency     *SiigoCurrency   `json:"currency,omitempty"`
	Items        []SiigoItem      `json:"items"`
	Payments     []SiigoPayment   `json:"payments,omitempty"`
	Stamp        *SiigoStamp      `json:"stamp,omitempty"`
	Mail         *SiigoMail       `json:"mail,omitempty"`
	Observations string           `json:"observations,omitempty"`
}

type SiigoDocument struct {
	ID int `json:"id"`
}

type SiigoCustomerRef struct {
	Identification string `json:"identification"`
	BranchOffice   int    `json:"branch_office"`
}

type SiigoStamp struct {
	Send bool `json:"send"`
}

type SiigoMail struct {
	Send bool `json:"send"`
}

type SiigoIDType struct {
	Code string `json:"code"`
}

type SiigoFiscalResponsibility struct {
	Code string `json:"code"`
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

type SiigoItem struct {
	Code        string     `json:"code"`
	Description string     `json:"description"`
	Quantity    float64    `json:"quantity"`
	Price       float64    `json:"price"`
	Discount    float64    `json:"discount,omitempty"`
	Taxes       []SiigoTax `json:"taxes,omitempty"`
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

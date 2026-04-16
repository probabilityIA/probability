package response

// GetBillCustomer datos del cliente en la factura detallada
type GetBillCustomer struct {
	Identification string  `json:"identification"`
	DV             string  `json:"dv"`
	Names          string  `json:"names"`
	Company        string  `json:"company"`
	TradeName      string  `json:"trade_name"`
	Email          *string `json:"email"`
	Phone          *string `json:"phone"`
	Address        string  `json:"address"`
}

// GetBillItem item de línea en la factura detallada
type GetBillItem struct {
	CodeReference string `json:"code_reference"`
	Name          string `json:"name"`
	Quantity      string `json:"quantity"`
	Price         string `json:"price"`
	DiscountRate  string `json:"discount_rate"`
	Discount      string `json:"discount"`
	TaxRate       string `json:"tax_rate"`
	TaxAmount     string `json:"tax_amount"`
	Total         string `json:"total"`
}

// GetBillDocument tipo de documento en la factura detallada
type GetBillDocument struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// GetBillPaymentForm forma de pago en la factura detallada
type GetBillPaymentForm struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// GetBillNote nota crédito o débito en la factura detallada
type GetBillNote struct {
	ID     int    `json:"id"`
	Number string `json:"number"`
}

// GetBill datos principales de la factura detallada
type GetBill struct {
	ID            int                `json:"id"`
	Number        string             `json:"number"`
	ReferenceCode *string            `json:"reference_code"`
	CUFE          string             `json:"cufe"`
	QR            string             `json:"qr"`
	QRImage       string             `json:"qr_image"`
	Status        int                `json:"status"`
	Total         string             `json:"total"`
	TaxAmount     string             `json:"tax_amount"`
	GrossValue    string             `json:"gross_value"`
	Discount      string             `json:"discount"`
	Validated     string             `json:"validated"`
	CreatedAt     string             `json:"created_at"`
	Document      GetBillDocument    `json:"document"`
	PaymentForm   GetBillPaymentForm `json:"payment_form"`
	CreditNotes   []GetBillNote      `json:"credit_notes"`
	DebitNotes    []GetBillNote      `json:"debit_notes"`
}

// GetBillData data anidada en la respuesta
type GetBillData struct {
	Customer GetBillCustomer `json:"customer"`
	Bill     GetBill         `json:"bill"`
	Items    []GetBillItem   `json:"items"`
}

// GetBillDetail respuesta completa de GET /v1/bills/show/:number
type GetBillDetail struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    GetBillData `json:"data"`
}

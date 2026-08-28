package response

type CreateInvoiceResponse struct {
	ID           string               `json:"id"`
	Document     InvoiceDocumentRef   `json:"document"`
	Date         string               `json:"date"`
	Number       int                  `json:"number"`
	Prefix       string               `json:"prefix"`
	Name         string               `json:"name"`
	Seller       int                  `json:"seller"`
	Customer     InvoiceCustomerInfo  `json:"customer"`
	Items        []InvoiceItem        `json:"items"`
	Payments     []InvoicePayment     `json:"payments"`
	Stamp        InvoiceStamp         `json:"stamp"`
	Mail         InvoiceMail          `json:"mail"`
	Total        float64              `json:"total"`
	TotalTax     float64              `json:"total_tax"`
	Discount     float64              `json:"discount"`
	Balance      float64              `json:"balance"`
	Observations string               `json:"observations"`
	ErrorCode    string               `json:"error_code,omitempty"`
	Errors       []SiigoError         `json:"Errors,omitempty"`
	PublicURL    string               `json:"public_url,omitempty"`
	Metadata     InvoiceMetadata      `json:"metadata,omitempty"`
	AdditionalIn []InvoiceAdditionalT `json:"additional_fields,omitempty"`
}

type InvoiceDocumentRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type InvoiceCustomerInfo struct {
	ID             string `json:"id"`
	Identification string `json:"identification"`
	Name           string `json:"name"`
	BranchOffice   int    `json:"branch_office"`
}

type InvoiceItem struct {
	ID          string          `json:"id"`
	Code        string          `json:"code"`
	Description string          `json:"description"`
	Quantity    float64         `json:"quantity"`
	Price       float64         `json:"price"`
	Discount    float64         `json:"discount"`
	Total       float64         `json:"total"`
	Taxes       []InvoiceItemTx `json:"taxes"`
}

type InvoiceItemTx struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Percentage float64 `json:"percentage"`
	Value      float64 `json:"value"`
}

type InvoicePayment struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	DueAt string  `json:"due_date,omitempty"`
}

type InvoiceStamp struct {
	Status       string `json:"status"`
	CUFE         string `json:"cufe,omitempty"`
	Observations string `json:"observations,omitempty"`
}

type InvoiceMail struct {
	Status       string `json:"status"`
	Observations string `json:"observations,omitempty"`
}

type InvoiceMetadata struct {
	Created string `json:"created,omitempty"`
	CUFE    string `json:"cufe,omitempty"`
	QR      string `json:"qr,omitempty"`
}

type InvoiceAdditionalT struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SiigoError struct {
	Code    string `json:"Code,omitempty"`
	Message string `json:"Message,omitempty"`
	Params  string `json:"Params,omitempty"`
}

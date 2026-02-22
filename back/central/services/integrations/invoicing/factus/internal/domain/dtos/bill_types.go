package dtos

// ═══════════════════════════════════════════════════════════════
// DOMAIN TYPES — Consulta de facturas Factus
// Sin tags (pureza de dominio)
// ═══════════════════════════════════════════════════════════════

// ListBillsParams parámetros para filtrar la consulta de facturas
type ListBillsParams struct {
	Page           int
	PerPage        int     // Registros por página (filter[per_page])
	Number         string  // Número de factura (ej: "SETP990000203")
	Prefix         string  // Prefijo (ej: "SETP")
	Identification string  // NIT / CC del cliente
	Names          string  // Nombre del cliente
	ReferenceCode  string  // Código de referencia interno
	Status         *int    // Estado: 1=validada, 0=no validada
	StartDate      string  // Fecha inicio YYYY-MM-DD
	EndDate        string  // Fecha fin YYYY-MM-DD
}

// BillDocument tipo de documento de una factura
type BillDocument struct {
	Code string
	Name string
}

// BillPaymentForm forma de pago de una factura
type BillPaymentForm struct {
	Code string
	Name string
}

// BillNote nota crédito o débito asociada a una factura
type BillNote struct {
	ID     int
	Number string
}

// Bill representa una factura electrónica retornada por Factus
type Bill struct {
	ID                        int
	Document                  BillDocument
	Number                    string // "SETP990000203"
	APIClientName             string
	ReferenceCode             *string
	Identification            string // NIT / CC del cliente
	GraphicRepresentationName string
	Company                   string
	TradeName                 *string
	Names                     string
	Email                     *string
	Total                     string // "90000.00" — string en la API de Factus
	Status                    int    // 1 = activa
	Errors                    []string
	SendEmail                 bool
	HasClaim                  bool
	IsNegotiableInstrument    bool
	PaymentForm               BillPaymentForm
	CreatedAt                 string // "17-07-2024 03:54:10 PM"
	CreditNotes               []BillNote
	DebitNotes                []BillNote
}

// BillsPagination información de paginación
type BillsPagination struct {
	Total       int
	PerPage     int
	CurrentPage int
	LastPage    int
	From        int
	To          int
}

// ListBillsResult resultado de consultar facturas en Factus
type ListBillsResult struct {
	Bills      []Bill
	Pagination BillsPagination
}

// BillDetailCustomer datos del cliente en una factura detallada
type BillDetailCustomer struct {
	Identification string
	DV             string
	Names          string
	Company        string
	TradeName      string
	Email          string
	Phone          string
	Address        string
}

// BillDetailItem item de una factura detallada
type BillDetailItem struct {
	CodeReference string
	Name          string
	Quantity      float64
	Price         string
	DiscountRate  string
	Discount      string
	TaxRate       string
	TaxAmount     string
	Total         string
}

// BillDetail representa una factura electrónica con todos sus detalles
// Retornada por GET /v1/bills/show/:number
type BillDetail struct {
	ID           int
	Number       string
	ReferenceCode *string
	CUFE         string
	QRCode       string
	QRImage      string // base64 data URI
	Status       int
	Total        string
	TaxAmount    string
	GrossValue   string
	Discount     string
	Validated    string
	CreatedAt    string
	Document     BillDocument
	PaymentForm  BillPaymentForm
	Customer     BillDetailCustomer
	Items        []BillDetailItem
	CreditNotes  []BillNote
	DebitNotes   []BillNote
}

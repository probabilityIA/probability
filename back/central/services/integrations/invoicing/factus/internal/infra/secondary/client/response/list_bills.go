package response

// BillDocument tipo de documento de una factura listada
type BillDocument struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// BillPaymentForm forma de pago de una factura listada
type BillPaymentForm struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// BillNote nota crédito o débito asociada a una factura listada
type BillNote struct {
	ID     int    `json:"id"`
	Number string `json:"number"`
}

// Bill representa una factura en el listado de GET /v1/bills
type Bill struct {
	ID                        int             `json:"id"`
	Document                  BillDocument    `json:"document"`
	Number                    string          `json:"number"`
	APIClientName             string          `json:"api_client_name"`
	ReferenceCode             *string         `json:"reference_code"`
	Identification            string          `json:"identification"`
	GraphicRepresentationName string          `json:"graphic_representation_name"`
	Company                   string          `json:"company"`
	TradeName                 *string         `json:"trade_name"`
	Names                     string          `json:"names"`
	Email                     *string         `json:"email"`
	Total                     string          `json:"total"`
	Status                    int             `json:"status"`
	Errors                    []string        `json:"errors"`
	SendEmail                 int             `json:"send_email"`
	HasClaim                  int             `json:"has_claim"`
	IsNegotiableInstrument    int             `json:"is_negotiable_instrument"`
	PaymentForm               BillPaymentForm `json:"payment_form"`
	CreatedAt                 string          `json:"created_at"`
	CreditNotes               []BillNote      `json:"credit_notes"`
	DebitNotes                []BillNote      `json:"debit_notes"`
}

// BillsPagination información de paginación del listado
type BillsPagination struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	From        int `json:"from"`
	To          int `json:"to"`
}

// BillsData data anidada en la respuesta
type BillsData struct {
	Data       []Bill          `json:"data"`
	Pagination BillsPagination `json:"pagination"`
}

// Bills respuesta completa de GET /v1/bills
type Bills struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Data    BillsData `json:"data"`
}

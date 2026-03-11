package request

// SiigoJournal representa el body de una solicitud de creación de comprobante contable en Siigo
// Endpoint: POST /v1/journals
type SiigoJournal struct {
	Document     SiigoDocument      `json:"document"`
	Date         string             `json:"date"` // "YYYY-MM-DD"
	Currency     SiigoCurrency      `json:"currency,omitempty"`
	Items        []SiigoJournalItem `json:"items"`
	Observations string             `json:"observations,omitempty"`
}

// SiigoJournalItem item del comprobante contable
type SiigoJournalItem struct {
	Account    SiigoJournalAccount   `json:"account"`
	Customer   *SiigoJournalCustomer `json:"customer,omitempty"`
	Product    *SiigoJournalProduct  `json:"product,omitempty"`
	Description string               `json:"description,omitempty"`
	CostCenter int                   `json:"cost_center,omitempty"`
	Movement   string                `json:"movement"` // "Debit" o "Credit"
	Value      float64               `json:"value"`
	Taxes      []SiigoTax            `json:"taxes,omitempty"`
}

// SiigoJournalAccount cuenta contable del item
type SiigoJournalAccount struct {
	Code string `json:"code"` // Ej: "11050501"
}

// SiigoJournalCustomer cliente asociado al item del journal
type SiigoJournalCustomer struct {
	Identification string `json:"identification"`
}

// SiigoJournalProduct producto asociado al item del journal
type SiigoJournalProduct struct {
	Code      string `json:"code"`
	Quantity  int    `json:"quantity"`
	Warehouse int    `json:"warehouse,omitempty"`
}

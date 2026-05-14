package request

type SiigoJournal struct {
	Document     SiigoDocument      `json:"document"`
	Date         string             `json:"date"`
	Number       int                `json:"number,omitempty"`
	Currency     *SiigoCurrency     `json:"currency,omitempty"`
	Items        []SiigoJournalItem `json:"items"`
	Observations string             `json:"observations,omitempty"`
}

type SiigoJournalItem struct {
	Account     SiigoJournalAccount   `json:"account"`
	Customer    *SiigoJournalCustomer `json:"customer,omitempty"`
	Product     *SiigoJournalProduct  `json:"product,omitempty"`
	Description string                `json:"description,omitempty"`
	CostCenter  int                   `json:"cost_center,omitempty"`
	Value       float64               `json:"value"`
	Tax         *SiigoJournalTax      `json:"tax,omitempty"`
	Due         *SiigoJournalDue      `json:"due,omitempty"`
	FixedAssets int                   `json:"fixed_assets,omitempty"`
}

type SiigoJournalAccount struct {
	Code     string `json:"code"`
	Movement string `json:"movement"`
}

type SiigoJournalCustomer struct {
	Identification string `json:"identification"`
	BranchOffice   int    `json:"branch_office,omitempty"`
}

type SiigoJournalProduct struct {
	Code      string  `json:"code"`
	Quantity  float64 `json:"quantity"`
	Warehouse int     `json:"warehouse,omitempty"`
}

type SiigoJournalTax struct {
	ID         int     `json:"id"`
	Name       string  `json:"name,omitempty"`
	Type       string  `json:"type,omitempty"`
	Percentage float64 `json:"percentage,omitempty"`
}

type SiigoJournalDue struct {
	Prefix      string `json:"prefix,omitempty"`
	Consecutive int    `json:"consecutive,omitempty"`
	Quote       int    `json:"quote,omitempty"`
	Date        string `json:"date,omitempty"`
}

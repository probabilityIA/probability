package domain

import "time"

type Customer struct {
	ID             string
	Identification string
	Name           string
	Email          string
	Phone          string
	CreatedAt      time.Time
}

type InvoiceItem struct {
	Code        string
	Description string
	Quantity    float64
	Price       float64
	Total       float64
}

type Invoice struct {
	ID          string
	Prefix      string
	Number      int
	Name        string
	Date        string
	CustomerID  string
	CustomerNIT string
	Items       []InvoiceItem
	Total       float64
	Balance     float64
	StampStatus string
	Status      string
	Annulled    bool
	CUFE        string
	PublicURL   string
	StampErrors []StampError
	CreatedAt   time.Time
}

type StampError struct {
	Code    string
	Message string
}

type JournalEntry struct {
	ID         string
	Date       string
	DocumentID string
	Items      []map[string]interface{}
	Total      float64
	CreatedAt  time.Time
}

type Product struct {
	ID          string
	Code        string
	Name        string
	Description string
	Price       float64
}

type PaymentType struct {
	ID   int
	Name string
	Type string
}

type Voucher struct {
	ID            string
	Name          string
	Number        int
	InvoiceNumber string
	Value         float64
	Date          string
	CreatedAt     time.Time
}

type AuthToken struct {
	Token     string
	ExpiresAt time.Time
	Username  string
}

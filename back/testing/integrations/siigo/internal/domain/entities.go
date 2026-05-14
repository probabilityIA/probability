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
	ID             string
	Prefix         string
	Number         int
	Name           string
	Date           string
	CustomerID     string
	CustomerNIT    string
	Items          []InvoiceItem
	Total          float64
	CUFE           string
	PublicURL      string
	CreatedAt      time.Time
}

type JournalEntry struct {
	ID        string
	Date      string
	DocumentID string
	Items     []map[string]interface{}
	Total     float64
	CreatedAt time.Time
}

type AuthToken struct {
	Token     string
	ExpiresAt time.Time
	Username  string
}

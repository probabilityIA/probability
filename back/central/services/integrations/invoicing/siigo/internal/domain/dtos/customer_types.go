package dtos

type CustomerResult struct {
	ID             string
	Name           string
	Identification string
	Email          string
	Phone          string
	Address        string
}

type CreateCustomerRequest struct {
	PersonType     string
	IDType         string
	Identification string
	Name           string
	Email          string
	Phone          string
	Address        string
	City           string
	CountryCode    string
	StateCode      string
	CityCode       string
	Credentials    Credentials
}

type ListInvoicesParams struct {
	Page        int
	PageSize    int
	DateFrom    string
	DateTo      string
	Credentials Credentials
}

type ListInvoicesResult struct {
	Items      []InvoiceSummary
	Total      int
	Page       int
	PageSize   int
	TotalPages int
}

type InvoiceSummary struct {
	ID           string
	Number       string
	Prefix       string
	Date         string
	CustomerName string
	CustomerID   string
	Total        float64
	Status       string
	StampStatus  string
	Observations string
	Annulled     bool
}

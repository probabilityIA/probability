package dtos

// CustomerResult resultado de consultar/crear un cliente en Siigo
type CustomerResult struct {
	ID             string
	Name           string
	Identification string
	Email          string
	Phone          string
	Address        string
}

// CreateCustomerRequest datos para crear un cliente en Siigo
type CreateCustomerRequest struct {
	PersonType     string // "Person" o "Company"
	IDType         string // Tipo de documento ("13"=CC, "22"=CE, etc.)
	Identification string
	Name           string
	Email          string
	Phone          string
	Address        string
	City           string
	Credentials    Credentials
}

// ListInvoicesParams parámetros para listar facturas en Siigo
type ListInvoicesParams struct {
	Page        int
	PageSize    int
	DateFrom    string
	DateTo      string
	Credentials Credentials
}

// ListInvoicesResult resultado de listar facturas en Siigo
type ListInvoicesResult struct {
	Items      []InvoiceSummary
	Total      int
	Page       int
	PageSize   int
	TotalPages int
}

// InvoiceSummary resumen de una factura de Siigo
type InvoiceSummary struct {
	ID            string
	Number        string
	Date          string
	CustomerName  string
	CustomerID    string
	Total         float64
	Status        string
}

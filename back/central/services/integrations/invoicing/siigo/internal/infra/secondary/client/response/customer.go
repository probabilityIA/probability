package response

// Customer representa un cliente en Siigo
type Customer struct {
	ID             string               `json:"id"`
	Name           []string             `json:"name"`
	PersonType     string               `json:"person_type"`
	IDType         CustomerIDType       `json:"id_type"`
	Identification string               `json:"identification"`
	Emails         []CustomerEmail      `json:"contacts,omitempty"`
	Phones         []CustomerPhone      `json:"phones,omitempty"`
	Address        *CustomerAddress     `json:"address,omitempty"`
}

// CustomerIDType tipo de documento del cliente
type CustomerIDType struct {
	Code string `json:"code"`
}

// CustomerEmail email de contacto del cliente
type CustomerEmail struct {
	Email string `json:"email"`
}

// CustomerPhone teléfono del cliente
type CustomerPhone struct {
	Number string `json:"number"`
}

// CustomerAddress dirección del cliente
type CustomerAddress struct {
	Address string `json:"address"`
}

// ListCustomersResponse respuesta de listar clientes
type ListCustomersResponse struct {
	Pagination PaginationInfo `json:"_pagination"`
	Results    []Customer     `json:"results"`
}

// PaginationInfo información de paginación
type PaginationInfo struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
}

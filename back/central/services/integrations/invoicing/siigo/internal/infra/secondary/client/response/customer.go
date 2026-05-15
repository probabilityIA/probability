package response

type Customer struct {
	ID                     string                   `json:"id"`
	Type                   string                   `json:"type"`
	PersonType             string                   `json:"person_type"`
	IDType                 CustomerIDType           `json:"id_type"`
	Identification         string                   `json:"identification"`
	BranchOffice           int                      `json:"branch_office"`
	Name                   []string                 `json:"name"`
	CommercialName         string                   `json:"commercial_name,omitempty"`
	Active                 bool                     `json:"active"`
	VatResponsible         bool                     `json:"vat_responsible"`
	FiscalResponsibilities []FiscalResponsibility   `json:"fiscal_responsibilities,omitempty"`
	Address                *CustomerAddress         `json:"address,omitempty"`
	Phones                 []CustomerPhone          `json:"phones,omitempty"`
	Contacts               []CustomerContact        `json:"contacts,omitempty"`
	Metadata               map[string]interface{}   `json:"metadata,omitempty"`
}

type CustomerIDType struct {
	Code string `json:"code"`
	Name string `json:"name,omitempty"`
}

type FiscalResponsibility struct {
	Code string `json:"code"`
	Name string `json:"name,omitempty"`
}

type CustomerContact struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
}

type CustomerPhone struct {
	Indicative string `json:"indicative,omitempty"`
	Number     string `json:"number"`
	Extension  string `json:"extension,omitempty"`
}

type CustomerAddress struct {
	Address    string             `json:"address"`
	City       *CustomerCityRef   `json:"city,omitempty"`
	PostalCode string             `json:"postal_code,omitempty"`
}

type CustomerCityRef struct {
	CountryCode string `json:"country_code,omitempty"`
	CountryName string `json:"country_name,omitempty"`
	StateCode   string `json:"state_code,omitempty"`
	StateName   string `json:"state_name,omitempty"`
	CityCode    string `json:"city_code,omitempty"`
	CityName    string `json:"city_name,omitempty"`
}

type ListCustomersResponse struct {
	Pagination PaginationInfo `json:"pagination"`
	Results    []Customer     `json:"results"`
}

type PaginationInfo struct {
	Page         int `json:"page"`
	PageSize     int `json:"page_size"`
	TotalResults int `json:"total_results"`
}

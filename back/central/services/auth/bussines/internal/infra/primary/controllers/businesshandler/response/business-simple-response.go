package response

// BusinessSimpleResponse representa un negocio en formato simplificado para dropdowns/selectores
type BusinessSimpleResponse struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	LogoURL         string `json:"logo_url,omitempty"`
	PrimaryColor    string `json:"primary_color,omitempty"`
	SecondaryColor  string `json:"secondary_color,omitempty"`
	TertiaryColor   string `json:"tertiary_color,omitempty"`
	QuaternaryColor string `json:"quaternary_color,omitempty"`
}

// GetBusinessesSimpleResponse representa la respuesta para obtener negocios en formato simple
type GetBusinessesSimpleResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    []BusinessSimpleResponse `json:"data"`
}

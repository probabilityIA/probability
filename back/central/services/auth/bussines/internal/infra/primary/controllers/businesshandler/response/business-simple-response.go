package response

// BusinessSimpleResponse representa un negocio en formato simplificado para dropdowns/selectores
type BusinessSimpleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// GetBusinessesSimpleResponse representa la respuesta para obtener negocios en formato simple
type GetBusinessesSimpleResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    []BusinessSimpleResponse `json:"data"`
}

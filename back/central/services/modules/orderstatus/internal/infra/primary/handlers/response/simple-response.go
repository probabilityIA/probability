package response

// OrderStatusSimpleResponse representa un estado de orden en formato simplificado para dropdowns/selectores
type OrderStatusSimpleResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	IsActive bool   `json:"is_active"`
}

// OrderStatusesSimpleResponse representa la respuesta para obtener estados de orden en formato simple
type OrderStatusesSimpleResponse struct {
	Success bool                        `json:"success"`
	Message string                      `json:"message"`
	Data    []OrderStatusSimpleResponse `json:"data"`
}

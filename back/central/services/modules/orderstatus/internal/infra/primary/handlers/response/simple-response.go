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

// OrderStatusCatalogResponse representa un estado de orden completo para el catálogo
type OrderStatusCatalogResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Color       string `json:"color"`
	Priority    int    `json:"priority"`
	IsActive    bool   `json:"is_active"`
}

// OrderStatusesCatalogResponse es el envelope de la lista completa
type OrderStatusesCatalogResponse struct {
	Success bool                         `json:"success"`
	Message string                       `json:"message"`
	Data    []OrderStatusCatalogResponse `json:"data"`
}

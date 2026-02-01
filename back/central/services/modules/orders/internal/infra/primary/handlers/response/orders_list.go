package response

// OrdersList representa la respuesta HTTP paginada de órdenes
// ✅ DTO HTTP - CON TAGS (json)
type OrdersList struct {
	Data       []OrderSummary `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

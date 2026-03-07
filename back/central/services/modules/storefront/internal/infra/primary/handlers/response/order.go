package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
)

// OrderResponse respuesta de orden para el storefront
type OrderResponse struct {
	ID          string              `json:"id"`
	OrderNumber string              `json:"order_number"`
	Status      string              `json:"status"`
	TotalAmount float64             `json:"total_amount"`
	Currency    string              `json:"currency"`
	CreatedAt   time.Time           `json:"created_at"`
	Items       []OrderItemResponse `json:"items"`
}

// OrderItemResponse respuesta de item de orden
type OrderItemResponse struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
	ImageURL    *string `json:"image_url"`
}

// OrdersListResponse respuesta paginada de ordenes
type OrdersListResponse struct {
	Data       []OrderResponse `json:"data"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// OrderFromEntity convierte una entidad a response
func OrderFromEntity(o *entities.StorefrontOrder) OrderResponse {
	items := make([]OrderItemResponse, len(o.Items))
	for i, item := range o.Items {
		items[i] = OrderItemResponse{
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
			ImageURL:    item.ImageURL,
		}
	}
	return OrderResponse{
		ID:          o.ID,
		OrderNumber: o.OrderNumber,
		Status:      o.Status,
		TotalAmount: o.TotalAmount,
		Currency:    o.Currency,
		CreatedAt:   o.CreatedAt,
		Items:       items,
	}
}

// OrdersFromEntities convierte un slice de entidades a responses
func OrdersFromEntities(orders []entities.StorefrontOrder) []OrderResponse {
	result := make([]OrderResponse, len(orders))
	for i := range orders {
		result[i] = OrderFromEntity(&orders[i])
	}
	return result
}

package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
)

type CustomerOrderItemResponse struct {
	ID           uint    `json:"id"`
	CustomerID   uint    `json:"customer_id"`
	BusinessID   uint    `json:"business_id"`
	OrderID      string  `json:"order_id"`
	OrderNumber  string  `json:"order_number"`
	ProductID    *string `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductSKU   string  `json:"product_sku"`
	ProductImage *string `json:"product_image"`
	Quantity     int     `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	TotalPrice   float64 `json:"total_price"`
	OrderStatus  string  `json:"order_status"`
	OrderedAt    time.Time `json:"ordered_at"`
}

func OrderItemFromEntity(o *entities.CustomerOrderItem) CustomerOrderItemResponse {
	return CustomerOrderItemResponse{
		ID:           o.ID,
		CustomerID:   o.CustomerID,
		BusinessID:   o.BusinessID,
		OrderID:      o.OrderID,
		OrderNumber:  o.OrderNumber,
		ProductID:    o.ProductID,
		ProductName:  o.ProductName,
		ProductSKU:   o.ProductSKU,
		ProductImage: o.ProductImage,
		Quantity:     o.Quantity,
		UnitPrice:    o.UnitPrice,
		TotalPrice:   o.TotalPrice,
		OrderStatus:  o.OrderStatus,
		OrderedAt:    o.OrderedAt,
	}
}

type CustomerOrderItemListResponse struct {
	Data       []CustomerOrderItemResponse `json:"data"`
	Total      int64                       `json:"total"`
	Page       int                         `json:"page"`
	PageSize   int                         `json:"page_size"`
	TotalPages int                         `json:"total_pages"`
}

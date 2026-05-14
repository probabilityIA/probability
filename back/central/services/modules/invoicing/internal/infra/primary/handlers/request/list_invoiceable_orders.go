package request

type ListInvoiceableOrdersRequest struct {
	Page                int    `form:"page"`
	PageSize            int    `form:"page_size"`
	BusinessID          uint   `form:"business_id"`
	StartDate           string `form:"start_date"`
	EndDate             string `form:"end_date"`
	OrderNumber         string `form:"order_number"`
	CustomerName        string `form:"customer_name"`
	CustomerEmail       string `form:"customer_email"`
	PaymentStatusID     uint   `form:"payment_status_id"`
	FulfillmentStatusID uint   `form:"fulfillment_status_id"`
	SortBy              string `form:"sort_by"`
	SortOrder           string `form:"sort_order"`
}

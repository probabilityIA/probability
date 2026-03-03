package domain

import "time"

// GetOrdersParams define los parámetros para consultar órdenes de WooCommerce.
type GetOrdersParams struct {
	Status  string
	After   *time.Time
	Before  *time.Time
	Page    int
	PerPage int
	OrderBy string
	Order   string
}

// GetOrdersResult contiene el resultado paginado de la consulta de órdenes.
type GetOrdersResult struct {
	Orders     []WooCommerceOrder
	Total      int
	TotalPages int
}

package domain

import (
	"fmt"
	"net/url"
	"time"
)

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

// ToQueryString construye los query params para la API REST de WooCommerce.
func (p *GetOrdersParams) ToQueryString() string {
	params := url.Values{}

	if p.Status != "" {
		params.Set("status", p.Status)
	}
	if p.After != nil {
		params.Set("after", p.After.Format(time.RFC3339))
	}
	if p.Before != nil {
		params.Set("before", p.Before.Format(time.RFC3339))
	}
	if p.Page > 0 {
		params.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.PerPage > 0 {
		params.Set("per_page", fmt.Sprintf("%d", p.PerPage))
	}
	if p.OrderBy != "" {
		params.Set("orderby", p.OrderBy)
	}
	if p.Order != "" {
		params.Set("order", p.Order)
	}

	return params.Encode()
}

// GetOrdersResult contiene el resultado paginado de la consulta de órdenes.
type GetOrdersResult struct {
	Orders     []WooCommerceOrder
	Total      int
	TotalPages int
}

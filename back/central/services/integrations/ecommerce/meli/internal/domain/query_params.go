package domain

import (
	"fmt"
	"net/url"
	"time"
)

// GetOrdersParams define los parámetros para consultar órdenes de MercadoLibre.
// MeLi usa offset/limit (no page/per_page) y un máximo de 50 por request.
type GetOrdersParams struct {
	Status   string     // "paid", "confirmed", "cancelled", etc.
	DateFrom *time.Time // order.date_created.from (RFC3339)
	DateTo   *time.Time // order.date_created.to (RFC3339)
	Offset   int        // pagination offset (0, 50, 100...)
	Limit    int        // pagination limit (max 50 en MeLi)
	Sort     string     // "date_asc" o "date_desc"
}

// ToQueryString construye los query params para la API de MeLi.
// El sellerID es obligatorio para buscar órdenes.
func (p *GetOrdersParams) ToQueryString(sellerID int64) string {
	params := url.Values{}
	params.Set("seller", fmt.Sprintf("%d", sellerID))

	if p.Status != "" {
		params.Set("order.status", p.Status)
	}
	if p.DateFrom != nil {
		params.Set("order.date_created.from", p.DateFrom.Format(time.RFC3339))
	}
	if p.DateTo != nil {
		params.Set("order.date_created.to", p.DateTo.Format(time.RFC3339))
	}
	if p.Offset > 0 {
		params.Set("offset", fmt.Sprintf("%d", p.Offset))
	}
	if p.Limit > 0 {
		limit := p.Limit
		if limit > 50 {
			limit = 50
		}
		params.Set("limit", fmt.Sprintf("%d", limit))
	}
	if p.Sort != "" {
		params.Set("sort", p.Sort)
	}

	return params.Encode()
}

// GetOrdersResult contiene el resultado paginado de la consulta de órdenes.
type GetOrdersResult struct {
	Orders []MeliOrder
	Total  int
	Offset int
	Limit  int
}

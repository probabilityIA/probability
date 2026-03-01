package domain

import (
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

// GetOrdersResult contiene el resultado paginado de la consulta de órdenes.
type GetOrdersResult struct {
	Orders []MeliOrder
	Total  int
	Offset int
	Limit  int
}

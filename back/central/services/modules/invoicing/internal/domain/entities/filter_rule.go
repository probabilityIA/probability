package entities

import "time"

// FilterRule representa una regla de filtrado para facturación
type FilterRule struct {
	Type      FilterType
	Operator  FilterOperator
	Value     interface{}
	FieldName string // Campo a validar
}

type FilterType string

const (
	// Filtros de monto
	FilterTypeMinAmount FilterType = "min_amount"
	FilterTypeMaxAmount FilterType = "max_amount"

	// Filtros de pago
	FilterTypePaymentStatus  FilterType = "payment_status"
	FilterTypePaymentMethods FilterType = "payment_methods"

	// Filtros de orden
	FilterTypeOrderTypes      FilterType = "order_types"
	FilterTypeExcludeStatuses FilterType = "exclude_statuses"

	// Filtros de productos
	FilterTypeExcludeProducts     FilterType = "exclude_products"
	FilterTypeIncludeProductsOnly FilterType = "include_products_only"
	FilterTypeMinItemsCount       FilterType = "min_items_count"
	FilterTypeMaxItemsCount       FilterType = "max_items_count"

	// Filtros de cliente
	FilterTypeCustomerTypes      FilterType = "customer_types"
	FilterTypeExcludeCustomerIDs FilterType = "exclude_customer_ids"

	// Filtros de ubicación
	FilterTypeShippingRegions FilterType = "shipping_regions"

	// Filtros de fecha
	FilterTypeDateRange FilterType = "date_range"
)

type FilterOperator string

const (
	OperatorGreaterThan      FilterOperator = "gt"
	OperatorGreaterThanEqual FilterOperator = "gte"
	OperatorLessThan         FilterOperator = "lt"
	OperatorLessThanEqual    FilterOperator = "lte"
	OperatorEqual            FilterOperator = "eq"
	OperatorNotEqual         FilterOperator = "ne"
	OperatorIn               FilterOperator = "in"
	OperatorNotIn            FilterOperator = "not_in"
	OperatorContains         FilterOperator = "contains"
	OperatorBetween          FilterOperator = "between"
)

// FilterConfig es la configuración completa de filtros
type FilterConfig struct {
	// Monto
	MinAmount *float64 `json:"min_amount,omitempty"`
	MaxAmount *float64 `json:"max_amount,omitempty"`

	// Pago
	PaymentStatus  *string `json:"payment_status,omitempty"`  // "paid", "unpaid", "partial"
	PaymentMethods []uint  `json:"payment_methods,omitempty"` // IDs de métodos permitidos

	// Orden
	OrderTypes      []string `json:"order_types,omitempty"`      // ["delivery", "pickup"]
	ExcludeStatuses []string `json:"exclude_statuses,omitempty"` // ["cancelled", "refunded"]

	// Productos
	ExcludeProducts     []string `json:"exclude_products,omitempty"`      // SKUs a excluir
	IncludeProductsOnly []string `json:"include_products_only,omitempty"` // Solo estos SKUs
	MinItemsCount       *int     `json:"min_items_count,omitempty"`
	MaxItemsCount       *int     `json:"max_items_count,omitempty"`

	// Cliente
	CustomerTypes      []string `json:"customer_types,omitempty"`       // ["natural", "juridica"]
	ExcludeCustomerIDs []string `json:"exclude_customer_ids,omitempty"` // IDs de clientes a excluir

	// Ubicación
	ShippingRegions []string `json:"shipping_regions,omitempty"` // ["Bogotá", "Medellín"]

	// Fecha
	DateRange *DateRangeFilter `json:"date_range,omitempty"`
}

type DateRangeFilter struct {
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

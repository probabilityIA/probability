package domain

// DashboardStats contiene todas las estadísticas del dashboard
type DashboardStats struct {
	TotalOrders             int64                         `json:"total_orders"`
	OrdersByIntegrationType []OrderCountByIntegrationType `json:"orders_by_integration_type"`
	TopCustomers            []TopCustomer                 `json:"top_customers"`
	OrdersByLocation        []OrderCountByLocation        `json:"orders_by_location"`
}

// OrderCountByIntegrationType representa el conteo de órdenes por tipo de integración
type OrderCountByIntegrationType struct {
	IntegrationType string `json:"integration_type"`
	Count           int64  `json:"count"`
}

// TopCustomer representa un cliente con su conteo de órdenes
type TopCustomer struct {
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	OrderCount    int64  `json:"order_count"`
}

// OrderCountByLocation representa el conteo de órdenes por ubicación geográfica
type OrderCountByLocation struct {
	City       string `json:"city"`
	State      string `json:"state"`
	OrderCount int64  `json:"order_count"`
}

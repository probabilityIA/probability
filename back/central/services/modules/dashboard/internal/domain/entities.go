package domain

// DashboardStats contiene todas las estadísticas del dashboard
type DashboardStats struct {
	// Existentes
	TotalOrders             int64                         `json:"total_orders"`
	OrdersByIntegrationType []OrderCountByIntegrationType `json:"orders_by_integration_type"`
	TopCustomers            []TopCustomer                 `json:"top_customers"`
	OrdersByLocation        []OrderCountByLocation        `json:"orders_by_location"`

	// Nuevas: Transportadores
	TopDrivers        []TopDriver        `json:"top_drivers"`
	DriversByLocation []DriverByLocation `json:"drivers_by_location"`

	// Nuevas: Productos
	TopProducts        []TopProduct        `json:"top_products"`
	ProductsByCategory []ProductByCategory `json:"products_by_category"`
	ProductsByBrand    []ProductByBrand    `json:"products_by_brand"`

	// Nuevas: Envíos
	ShipmentsByStatus    []ShipmentsByStatus    `json:"shipments_by_status"`
	ShipmentsByCarrier   []ShipmentsByCarrier   `json:"shipments_by_carrier"`
	ShipmentsByWarehouse []ShipmentsByWarehouse `json:"shipments_by_warehouse"`

	// Nuevas: Businesses (solo si es super admin)
	OrdersByBusiness []OrdersByBusiness `json:"orders_by_business,omitempty"`
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

// TopDriver representa un transportador con su conteo de órdenes
type TopDriver struct {
	DriverName string `json:"driver_name"`
	DriverID   *uint  `json:"driver_id"`
	OrderCount int64  `json:"order_count"`
}

// DriverByLocation representa transportadores agrupados por ubicación
type DriverByLocation struct {
	DriverName string `json:"driver_name"`
	City       string `json:"city"`
	State      string `json:"state"`
	OrderCount int64  `json:"order_count"`
}

// TopProduct representa un producto con su conteo de órdenes
type TopProduct struct {
	ProductName string  `json:"product_name"`
	ProductID   string  `json:"product_id"`
	SKU         string  `json:"sku"`
	OrderCount  int64   `json:"order_count"`
	TotalSold   float64 `json:"total_sold"`
}

// ProductByCategory representa productos agrupados por categoría
type ProductByCategory struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}

// ProductByBrand representa productos agrupados por marca
type ProductByBrand struct {
	Brand string `json:"brand"`
	Count int64  `json:"count"`
}

// ShipmentsByStatus representa envíos agrupados por estado
type ShipmentsByStatus struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// ShipmentsByCarrier representa envíos agrupados por transportista
type ShipmentsByCarrier struct {
	Carrier string `json:"carrier"`
	Count   int64  `json:"count"`
}

// ShipmentsByWarehouse representa envíos agrupados por almacén
type ShipmentsByWarehouse struct {
	WarehouseName string `json:"warehouse_name"`
	WarehouseID   *uint  `json:"warehouse_id"`
	Count         int64  `json:"count"`
}

// OrdersByBusiness representa órdenes agrupadas por business (solo super admin)
type OrdersByBusiness struct {
	BusinessID   uint   `json:"business_id"`
	BusinessName string `json:"business_name"`
	OrderCount   int64  `json:"order_count"`
}

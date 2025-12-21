// Dashboard Stats Types

export interface DashboardStats {
    // Existentes
    total_orders: number;
    orders_by_integration_type: OrderCountByIntegrationType[];
    top_customers: TopCustomer[];
    orders_by_location: OrderCountByLocation[];

    // Nuevas: Transportadores
    top_drivers: TopDriver[];
    drivers_by_location: DriverByLocation[];

    // Nuevas: Productos
    top_products: TopProduct[];
    products_by_category: ProductByCategory[];
    products_by_brand: ProductByBrand[];

    // Nuevas: Envíos
    shipments_by_status: ShipmentsByStatus[];
    shipments_by_carrier: ShipmentsByCarrier[];
    shipments_by_warehouse: ShipmentsByWarehouse[];

    // Nuevas: Businesses (solo si es super admin)
    orders_by_business?: OrdersByBusiness[];
}

export interface OrderCountByIntegrationType {
    integration_type: string;
    count: number;
}

export interface TopCustomer {
    customer_name: string;
    customer_email: string;
    order_count: number;
}

export interface OrderCountByLocation {
    city: string;
    state: string;
    order_count: number;
}

// Transportadores
export interface TopDriver {
    driver_name: string;
    driver_id?: number;
    order_count: number;
}

export interface DriverByLocation {
    driver_name: string;
    city: string;
    state: string;
    order_count: number;
}

// Productos
export interface TopProduct {
    product_name: string;
    product_id: string;
    sku: string;
    order_count: number;
    total_sold: number;
}

export interface ProductByCategory {
    category: string;
    count: number;
}

export interface ProductByBrand {
    brand: string;
    count: number;
}

// Envíos
export interface ShipmentsByStatus {
    status: string;
    count: number;
}

export interface ShipmentsByCarrier {
    carrier: string;
    count: number;
}

export interface ShipmentsByWarehouse {
    warehouse_name: string;
    warehouse_id?: number;
    count: number;
}

// Businesses (solo super admin)
export interface OrdersByBusiness {
    business_id: number;
    business_name: string;
    order_count: number;
}

export interface DashboardStatsResponse {
    success: boolean;
    message: string;
    data: DashboardStats;
}

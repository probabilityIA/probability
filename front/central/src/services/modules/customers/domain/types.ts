export interface CustomerInfo {
    id: number;
    business_id: number;
    name: string;
    email: string | null;
    phone: string;
    dni: string | null;
    created_at: string;
    updated_at: string;
}

export interface CustomerDetail extends CustomerInfo {
    order_count: number;
    total_spent: number;
    last_order_at: string | null;
}

export interface CustomerSummary {
    id: number;
    customer_id: number;
    business_id: number;
    total_orders: number;
    delivered_orders: number;
    cancelled_orders: number;
    in_progress_orders: number;
    total_spent: number;
    avg_ticket: number;
    total_paid_orders: number;
    avg_delivery_score: number;
    first_order_at: string | null;
    last_order_at: string | null;
    preferred_platform: string;
    last_updated_at: string;
}

export interface CustomerAddress {
    id: number;
    customer_id: number;
    business_id: number;
    street: string;
    city: string;
    state: string;
    country: string;
    postal_code: string;
    times_used: number;
    last_used_at: string;
}

export interface CustomerProduct {
    id: number;
    customer_id: number;
    business_id: number;
    product_id: string;
    product_name: string;
    product_sku: string;
    product_image: string | null;
    times_ordered: number;
    total_quantity: number;
    total_spent: number;
    first_ordered_at: string;
    last_ordered_at: string;
}

export interface CustomerOrderItem {
    id: number;
    customer_id: number;
    business_id: number;
    order_id: string;
    order_number: string;
    product_id: string | null;
    product_name: string;
    product_sku: string;
    product_image: string | null;
    quantity: number;
    unit_price: number;
    total_price: number;
    order_status: string;
    ordered_at: string;
}

export interface CreateCustomerDTO {
    name: string;
    email?: string;
    phone?: string;
    dni?: string | null;
}

export interface UpdateCustomerDTO {
    name: string;
    email?: string;
    phone?: string;
    dni?: string | null;
}

export interface GetCustomersParams {
    page?: number;
    page_size?: number;
    search?: string;
    business_id?: number;
}

export interface PaginationParams {
    page?: number;
    page_size?: number;
    business_id?: number;
}

export interface CustomersListResponse {
    data: CustomerInfo[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface CustomerAddressListResponse {
    data: CustomerAddress[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface CustomerProductListResponse {
    data: CustomerProduct[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface CustomerOrderItemListResponse {
    data: CustomerOrderItem[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface DeleteCustomerResponse {
    message?: string;
    error?: string;
}

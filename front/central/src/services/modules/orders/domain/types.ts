export interface OrderStatusInfo {
    id: number;
    code: string;
    name: string;
    description?: string;
    category?: string;
    color?: string;
}

export interface PaymentStatusInfo {
    id: number;
    code: string;
    name: string;
    description?: string;
    category?: string;
    color?: string;
}

export interface FulfillmentStatusInfo {
    id: number;
    code: string;
    name: string;
    description?: string;
    category?: string;
    color?: string;
}

export interface Order {
    id: string;
    created_at: string;
    updated_at: string;
    deleted_at?: string;

    business_id?: number;
    integration_id: number;
    integration_type: string;
    integration_logo_url?: string;
    integration_name?: string;

    platform: string;
    external_id: string;
    order_number: string;
    internal_number: string;

    subtotal: number;
    tax: number;
    discount: number;
    shipping_cost: number;
    shipping_discount?: number;
    shipping_discount_presentment?: number;
    total_amount: number;
    currency: string;
    cod_total?: number;
    is_cod?: boolean;

    subtotal_presentment?: number;
    tax_presentment?: number;
    discount_presentment?: number;
    shipping_cost_presentment?: number;
    total_amount_presentment?: number;
    currency_presentment?: string;

    customer_id?: number;
    customer_name: string;
    customer_first_name?: string;
    customer_last_name?: string;
    customer_email: string;
    customer_phone: string;
    customer_dni: string;

    shipping_street: string;
    shipping_city: string;
    shipping_state: string;
    shipping_country: string;
    shipping_postal_code: string;
    shipping_house?: string;
    shipping_barrio?: string;
    shipping_lat?: number;
    shipping_lng?: number;
    shipping_geo_confidence?: 'high' | 'medium' | 'low';

    payment_method_id: number;
    is_paid: boolean;
    paid_at?: string;

    tracking_number?: string;
    tracking_link?: string;
    guide_id?: string;
    guide_link?: string;
    delivery_date?: string;
    delivered_at?: string;
    delivery_probability?: number;

    shipment?: {
        id: number;
        carrier?: string;
        tracking_number?: string;
        guide_url?: string;
        status: string;
        total_cost?: number;
    };

    warehouse_id?: number;
    warehouse_name: string;
    driver_id?: number;
    driver_name: string;
    is_last_mile: boolean;

    weight?: number;
    height?: number;
    width?: number;
    length?: number;
    boxes?: string;

    order_type_id?: number;
    order_type_name: string;
    status: string;
    original_status: string;
    status_id?: number;
    order_status?: OrderStatusInfo;

    payment_status_id?: number;
    fulfillment_status_id?: number;
    payment_status?: PaymentStatusInfo;
    fulfillment_status?: FulfillmentStatusInfo;

    notes?: string;
    coupon?: string;
    approved?: boolean;
    user_id?: number;
    user_name: string;

    is_confirmed?: boolean | null;
    novelty?: string;

    is_test?: boolean;

    invoiceable: boolean;
    invoice_url?: string;
    invoice_id?: string;
    invoice_provider?: string;
    invoice_status?: string;

    order_status_url?: string;

    items?: any;
    order_items?: any;
    metadata?: any;
    financial_details?: any;
    shipping_details?: any;
    payment_details?: any;
    fulfillment_details?: any;

    invoice?: {
        id: number;
        invoice_number: string;
        status: string;
        issued_at?: string;
        retention_amount: number;
    };

    occurred_at: string;
    imported_at: string;

    negative_factors?: string[];
    score_breakdown?: {
        final_score: number;
        categories: {
            name: string;
            weight: number;
            raw_score: number;
            weighted_score: number;
            factors: string[] | null;
        }[];
        negative_factors: string[];
    };
}

export interface PaginatedResponse<T> {
    success: boolean;
    message: string;
    data: T[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface SingleResponse<T> {
    success: boolean;
    message: string;
    data: T;
}

export interface ActionResponse {
    success: boolean;
    message: string;
    error?: string;
}

export interface OrderHistory {
    id: number;
    created_at: string;
    order_id: string;
    previous_status: string;
    new_status: string;
    changed_by?: number;
    changed_by_name: string;
    reason?: string;
}

export interface GetOrdersParams {
    page?: number;
    page_size?: number;
    business_id?: number;
    integration_id?: number;
    integration_type?: string;
    status?: string;
    customer_name?: string;
    customer_email?: string;
    customer_phone?: string;
    order_number?: string;
    internal_number?: string;
    platform?: string;
    currency_presentment?: string;
    is_paid?: boolean;
    is_cod?: boolean;
    payment_status_id?: number;
    fulfillment_status_id?: number;
    warehouse_id?: number;
    driver_id?: number;
    start_date?: string;
    end_date?: string;
    invoice_status?: string;
    sort_by?: 'created_at' | 'updated_at' | 'total_amount' | 'order_number';
    sort_order?: 'asc' | 'desc';
}

export interface CreateOrderDTO {
    business_id?: number;
    integration_id: number;
    integration_type: string;

    platform: string;
    external_id: string;
    order_number?: string;
    internal_number?: string;

    subtotal: number;
    tax?: number;
    discount?: number;
    shipping_cost?: number;
    total_amount: number;
    currency?: string;
    cod_total?: number;
    is_cod?: boolean;

    customer_id?: number;
    customer_name?: string;
    customer_first_name?: string;
    customer_last_name?: string;
    customer_email?: string;
    customer_phone?: string;
    customer_dni?: string;
    client_group_id?: number;

    shipping_street?: string;
    shipping_city?: string;
    shipping_state?: string;
    shipping_country?: string;
    shipping_postal_code?: string;
    shipping_house?: string;
    shipping_barrio?: string;
    shipping_lat?: number;
    shipping_lng?: number;

    payment_method_id: number;
    is_paid?: boolean;
    paid_at?: string;

    tracking_number?: string;
    tracking_link?: string;
    guide_id?: string;
    guide_link?: string;
    delivery_date?: string;
    delivered_at?: string;

    warehouse_id?: number;
    warehouse_name?: string;
    driver_id?: number;
    driver_name?: string;
    is_last_mile?: boolean;

    weight?: number;
    height?: number;
    width?: number;
    length?: number;
    boxes?: string;

    order_type_id?: number;
    order_type_name?: string;
    status?: string;
    original_status?: string;

    notes?: string;
    coupon?: string;
    approved?: boolean;
    user_id?: number;
    user_name?: string;

    invoiceable?: boolean;
    invoice_url?: string;
    invoice_id?: string;
    invoice_provider?: string;

    items?: any;
    metadata?: any;
    financial_details?: any;
    shipping_details?: any;
    payment_details?: any;
    fulfillment_details?: any;

    occurred_at?: string;
    imported_at?: string;
}

export interface SimulateShopifyResult {
    total: number;
    sent: number;
    failed: number;
    errors: string[];
}

export interface UpdateOrderDTO {
    subtotal?: number;
    tax?: number;
    discount?: number;
    shipping_cost?: number;
    total_amount?: number;
    currency?: string;
    cod_total?: number;
    is_cod?: boolean;

    customer_name?: string;
    customer_first_name?: string;
    customer_last_name?: string;
    customer_email?: string;
    customer_phone?: string;
    customer_dni?: string;

    shipping_street?: string;
    shipping_city?: string;
    shipping_state?: string;
    shipping_country?: string;
    shipping_postal_code?: string;
    shipping_house?: string;
    shipping_barrio?: string;
    shipping_lat?: number;
    shipping_lng?: number;

    payment_method_id?: number;
    is_paid?: boolean;
    paid_at?: string;

    tracking_number?: string;
    tracking_link?: string;
    guide_id?: string;
    guide_link?: string;
    delivery_date?: string;
    delivered_at?: string;

    warehouse_id?: number;
    warehouse_name?: string;
    driver_id?: number;
    driver_name?: string;
    is_last_mile?: boolean;

    weight?: number;
    height?: number;
    width?: number;
    length?: number;
    boxes?: string;

    order_type_id?: number;
    order_type_name?: string;
    status?: string;
    original_status?: string;
    status_id?: number;
    payment_status_id?: number;
    fulfillment_status_id?: number;

    notes?: string;
    coupon?: string;
    approved?: boolean;
    user_id?: number;
    user_name?: string;

    is_confirmed?: boolean | null;
    confirmation_status?: 'yes' | 'no' | 'pending';
    novelty?: string;

    is_test?: boolean;

    invoiceable?: boolean;
    invoice_url?: string;
    invoice_id?: string;
    invoice_provider?: string;

    items?: any;
    metadata?: any;
    financial_details?: any;
    shipping_details?: any;
    payment_details?: any;
    fulfillment_details?: any;
}

export interface ChangeOrderStatusDTO {
    status: string;
    metadata?: Record<string, unknown>;
}

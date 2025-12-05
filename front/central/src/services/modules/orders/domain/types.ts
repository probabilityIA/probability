export interface Order {
    id: string;
    created_at: string;
    updated_at: string;
    deleted_at?: string;

    // Identificadores de integración
    business_id?: number;
    integration_id: number;
    integration_type: string;

    // Identificadores de la orden
    platform: string;
    external_id: string;
    order_number: string;
    internal_number: string;

    // Información financiera
    subtotal: number;
    tax: number;
    discount: number;
    shipping_cost: number;
    total_amount: number;
    currency: string;
    cod_total?: number;

    // Información del cliente
    customer_id?: number;
    customer_name: string;
    customer_email: string;
    customer_phone: string;
    customer_dni: string;

    // Dirección de envío
    shipping_street: string;
    shipping_city: string;
    shipping_state: string;
    shipping_country: string;
    shipping_postal_code: string;
    shipping_lat?: number;
    shipping_lng?: number;

    // Información de pago
    payment_method_id: number;
    is_paid: boolean;
    paid_at?: string;

    // Información de envío/logística
    tracking_number?: string;
    tracking_link?: string;
    guide_id?: string;
    guide_link?: string;
    delivery_date?: string;
    delivered_at?: string;

    // Información de fulfillment
    warehouse_id?: number;
    warehouse_name: string;
    driver_id?: number;
    driver_name: string;
    is_last_mile: boolean;

    // Dimensiones y peso
    weight?: number;
    height?: number;
    width?: number;
    length?: number;
    boxes?: string;

    // Tipo y estado
    order_type_id?: number;
    order_type_name: string;
    status: string;
    original_status: string;

    // Información adicional
    notes?: string;
    coupon?: string;
    approved?: boolean;
    user_id?: number;
    user_name: string;

    // Facturación
    invoiceable: boolean;
    invoice_url?: string;
    invoice_id?: string;
    invoice_provider?: string;

    // Datos estructurados (JSONB)
    items?: any;
    metadata?: any;
    financial_details?: any;
    shipping_details?: any;
    payment_details?: any;
    fulfillment_details?: any;

    // Timestamps
    occurred_at: string;
    imported_at: string;
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

export interface GetOrdersParams {
    page?: number;
    page_size?: number;
    business_id?: number;
    integration_id?: number;
    integration_type?: string;
    status?: string;
    customer_email?: string;
    customer_phone?: string;
    platform?: string;
    is_paid?: boolean;
    warehouse_id?: number;
    driver_id?: number;
    start_date?: string;
    end_date?: string;
    sort_by?: 'created_at' | 'updated_at' | 'total_amount';
    sort_order?: 'asc' | 'desc';
}

export interface CreateOrderDTO {
    // Identificadores de integración
    business_id?: number;
    integration_id: number;
    integration_type: string;

    // Identificadores de la orden
    platform: string;
    external_id: string;
    order_number?: string;
    internal_number?: string;

    // Información financiera
    subtotal: number;
    tax?: number;
    discount?: number;
    shipping_cost?: number;
    total_amount: number;
    currency?: string;
    cod_total?: number;

    // Información del cliente
    customer_id?: number;
    customer_name?: string;
    customer_email?: string;
    customer_phone?: string;
    customer_dni?: string;

    // Dirección de envío
    shipping_street?: string;
    shipping_city?: string;
    shipping_state?: string;
    shipping_country?: string;
    shipping_postal_code?: string;
    shipping_lat?: number;
    shipping_lng?: number;

    // Información de pago
    payment_method_id: number;
    is_paid?: boolean;
    paid_at?: string;

    // Información de envío/logística
    tracking_number?: string;
    tracking_link?: string;
    guide_id?: string;
    guide_link?: string;
    delivery_date?: string;
    delivered_at?: string;

    // Información de fulfillment
    warehouse_id?: number;
    warehouse_name?: string;
    driver_id?: number;
    driver_name?: string;
    is_last_mile?: boolean;

    // Dimensiones y peso
    weight?: number;
    height?: number;
    width?: number;
    length?: number;
    boxes?: string;

    // Tipo y estado
    order_type_id?: number;
    order_type_name?: string;
    status?: string;
    original_status?: string;

    // Información adicional
    notes?: string;
    coupon?: string;
    approved?: boolean;
    user_id?: number;
    user_name?: string;

    // Facturación
    invoiceable?: boolean;
    invoice_url?: string;
    invoice_id?: string;
    invoice_provider?: string;

    // Datos estructurados (JSONB)
    items?: any;
    metadata?: any;
    financial_details?: any;
    shipping_details?: any;
    payment_details?: any;
    fulfillment_details?: any;

    // Timestamps
    occurred_at?: string;
    imported_at?: string;
}

export interface UpdateOrderDTO {
    // Información financiera
    subtotal?: number;
    tax?: number;
    discount?: number;
    shipping_cost?: number;
    total_amount?: number;
    currency?: string;
    cod_total?: number;

    // Información del cliente
    customer_name?: string;
    customer_email?: string;
    customer_phone?: string;
    customer_dni?: string;

    // Dirección de envío
    shipping_street?: string;
    shipping_city?: string;
    shipping_state?: string;
    shipping_country?: string;
    shipping_postal_code?: string;
    shipping_lat?: number;
    shipping_lng?: number;

    // Información de pago
    payment_method_id?: number;
    is_paid?: boolean;
    paid_at?: string;

    // Información de envío/logística
    tracking_number?: string;
    tracking_link?: string;
    guide_id?: string;
    guide_link?: string;
    delivery_date?: string;
    delivered_at?: string;

    // Información de fulfillment
    warehouse_id?: number;
    warehouse_name?: string;
    driver_id?: number;
    driver_name?: string;
    is_last_mile?: boolean;

    // Dimensiones y peso
    weight?: number;
    height?: number;
    width?: number;
    length?: number;
    boxes?: string;

    // Tipo y estado
    order_type_id?: number;
    order_type_name?: string;
    status?: string;
    original_status?: string;

    // Información adicional
    notes?: string;
    coupon?: string;
    approved?: boolean;
    user_id?: number;
    user_name?: string;

    // Facturación
    invoiceable?: boolean;
    invoice_url?: string;
    invoice_id?: string;
    invoice_provider?: string;

    // Datos estructurados (JSONB)
    items?: any;
    metadata?: any;
    financial_details?: any;
    shipping_details?: any;
    payment_details?: any;
    fulfillment_details?: any;
}

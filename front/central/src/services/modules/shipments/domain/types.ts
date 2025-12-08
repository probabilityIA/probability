export interface Shipment {
    id: number;
    created_at: string;
    updated_at: string;
    order_id: string;
    tracking_number?: string;
    tracking_url?: string;
    carrier?: string;
    carrier_code?: string;
    guide_id?: string;
    guide_url?: string;
    status: 'pending' | 'in_transit' | 'delivered' | 'failed';
    shipped_at?: string;
    delivered_at?: string;
    shipping_cost?: number;
    insurance_cost?: number;
    total_cost?: number;
    weight?: number;
    height?: number;
    width?: number;
    length?: number;
    warehouse_name?: string;
    driver_name?: string;
    is_last_mile: boolean;
    estimated_delivery?: string;
    delivery_notes?: string;
}

export interface GetShipmentsParams {
    page?: number;
    page_size?: number;
    order_id?: string;
    tracking_number?: string;
    carrier?: string;
    status?: string;
    start_date?: string;
    end_date?: string;
    shipped_after?: string;
    shipped_before?: string;
    delivered_after?: string;
    delivered_before?: string;
    sort_by?: string;
    sort_order?: 'asc' | 'desc';
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

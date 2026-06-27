export interface CarrierAggregate {
    carrier: string;
    orders_count: number;
    total_collected: number;
    discount_pct: number;
    total_discount: number;
    total_net: number;
}

export interface MonthlyPoint {
    month: string;
    label: string;
    orders: number;
    collected: number;
    discount: number;
    net: number;
}

export interface CodSummary {
    total_collected: number;
    total_pending: number;
    total_discount: number;
    total_net: number;
    orders_collected: number;
    orders_pending: number;
    by_carrier: CarrierAggregate[];
    monthly: MonthlyPoint[];
}

export interface CodOrder {
    order_id: string;
    order_number: string;
    shipment_id: number;
    has_guide: boolean;
    customer_name: string;
    carrier: string;
    cod_total: number;
    cod_carrier_fee: number;
    shipping_cost: number;
    discount_pct: number;
    discount: number;
    net: number;
    currency: string;
    status: string;
    collected: boolean;
    cut_status: string;
    created_at: string;
    delivered_at: string | null;
}

export interface PaymentCut {
    id: number;
    period_start: string;
    period_end: string;
    status: string;
    orders_count: number;
    total_collected: number;
    total_discount: number;
    total_net: number;
    by_carrier: CarrierAggregate[];
    confirmed_by: number;
    confirmed_by_name: string;
    confirmed_at: string | null;
}

export interface CarrierConfig {
    id: number;
    carrier_name: string;
    discount_percentage: number;
    is_active: boolean;
}

export type RangeKey = 'today' | 'week' | 'month' | '3months' | 'custom';

export interface ReportFilters {
    range: RangeKey;
    start_date?: string;
    end_date?: string;
    carrier?: string;
    business_id?: number;
}

export interface CodOrdersParams extends ReportFilters {
    page?: number;
    page_size?: number;
    collected?: boolean;
    has_guide?: boolean;
    search?: string;
}

export interface SaveCarrierConfigInput {
    carrier_name: string;
    discount_percentage: number;
    is_active: boolean;
}

export interface Paginated<T> {
    success: boolean;
    message?: string;
    data: T[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface SingleResult<T> {
    success: boolean;
    message?: string;
    data: T;
}

export interface CutsResult {
    success: boolean;
    message?: string;
    data: PaymentCut[];
    can_confirm: boolean;
}

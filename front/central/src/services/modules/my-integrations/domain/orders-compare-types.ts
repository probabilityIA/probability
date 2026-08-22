export type OrderCompareAction = 'create' | 'in_sync' | 'only_in_probability';

export interface OrderCompareRow {
    external_id: string;
    number: string;
    customer_name: string;
    channel_status: string;
    raw_status: string;
    local_status?: string;
    order_id?: string;
    order_number?: string;
    total: number;
    local_total?: number;
    currency: string;
    items: number;
    created_at: string;
    url?: string;
    action: OrderCompareAction;
    moves_inventory: boolean;
    inventory_note?: string;
    status_mismatch: boolean;
    total_mismatch: boolean;
}

export interface OrderCompareTotals {
    total: number;
    to_create: number;
    in_sync: number;
    only_in_probability: number;
    without_inventory: number;
    with_status_mismatch: number;
}

export interface OrdersCompareChannel {
    integration_id: number;
    integration_name: string;
    integration_type_id: number;
    supported: boolean;
}

export interface OrdersComparePage {
    rows: OrderCompareRow[];
    totals: OrderCompareTotals;
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
    checked_at: string;
    channel: OrdersCompareChannel;
}

export interface OrdersCompareQuery {
    integrationId: number;
    businessId?: number;
    from?: string;
    to?: string;
    page?: number;
    pageSize?: number;
    limit?: number;
    onlyDiff?: boolean;
    search?: string;
}

export interface OrdersApplyResult {
    queued: string[];
    skipped: string[];
    failed?: Record<string, string>;
    without_inventory: string[];
    note?: string;
}

export const ORDERS_COMPARE_TYPE_IDS = [1, 3, 4, 16, 17, 33];

export interface PaginationParams {
    page?: number;
    page_size?: number;
}

export interface PaginatedResponse<T> {
    data: T[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface InventoryLevel {
    id: number;
    product_id: string;
    warehouse_id: number;
    location_id: number | null;
    state_id?: number | null;
    business_id: number;
    quantity: number;
    reserved_qty: number;
    available_qty: number;
    min_stock: number | null;
    max_stock: number | null;
    reorder_point: number | null;
    product_name?: string;
    product_sku?: string;
    warehouse_name?: string;
    warehouse_code?: string;
    state_name?: string;
    location_name?: string;
    location_code?: string;
    created_at: string;
    updated_at: string;
}

export interface StockMovement {
    id: number;
    product_id: string;
    warehouse_id: number;
    location_id: number | null;
    business_id: number;
    movement_type_id: number;
    movement_type_code: string;
    movement_type_name: string;
    reason: string;
    quantity: number;
    previous_qty: number;
    new_qty: number;
    reference_type: string | null;
    reference_id: string | null;
    integration_id: number | null;
    notes: string;
    created_by_id: number | null;
    product_name?: string;
    product_sku?: string;
    warehouse_name?: string;
    created_at: string;
}

export interface MovementType {
    id: number;
    code: string;
    name: string;
    description: string;
    is_active: boolean;
    direction: string;
    created_at: string;
    updated_at: string;
}

export interface GetInventoryParams extends PaginationParams {
    search?: string;
    low_stock?: boolean;
    business_id?: number;
}

export interface GetMovementsParams extends PaginationParams {
    product_id?: string;
    warehouse_id?: number;
    type?: string;
    business_id?: number;
}

export interface AdjustStockDTO {
    product_id: string;
    warehouse_id: number;
    location_id?: number | null;
    lot_id?: number | null;
    state_id?: number | null;
    uom_id?: number | null;
    quantity: number;
    reason: string;
    notes?: string;
}

export interface TransferStockDTO {
    product_id: string;
    from_warehouse_id: number;
    to_warehouse_id: number;
    from_location_id?: number | null;
    to_location_id?: number | null;
    lot_id?: number | null;
    state_id?: number | null;
    uom_id?: number | null;
    quantity: number;
    reason?: string;
    notes?: string;
}

export interface BulkLoadItem {
    sku: string;
    quantity: number;
    min_stock?: number | null;
    max_stock?: number | null;
    reorder_point?: number | null;
}

export interface BulkLoadDTO {
    warehouse_id: number;
    reason?: string;
    items: BulkLoadItem[];
}

export interface BulkLoadItemResult {
    sku: string;
    product_id: string;
    success: boolean;
    previous_qty: number;
    new_qty: number;
    error?: string;
}

export interface BulkLoadResult {
    total_items: number;
    success_count: number;
    failure_count: number;
    items: BulkLoadItemResult[];
}

export interface BulkLoadAccepted {
    message: string;
    total_items: number;
}

export type InventoryListResponse = PaginatedResponse<InventoryLevel>;
export type MovementListResponse = PaginatedResponse<StockMovement>;
export type MovementTypeListResponse = PaginatedResponse<MovementType>;

// ============================================
// ENTIDADES
// ============================================

export interface InventoryLevel {
    id: number;
    product_id: string;
    warehouse_id: number;
    location_id: number | null;
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

// ============================================
// DTOs
// ============================================

export interface GetInventoryParams {
    page?: number;
    page_size?: number;
    search?: string;
    low_stock?: boolean;
    business_id?: number;
}

export interface GetMovementsParams {
    page?: number;
    page_size?: number;
    product_id?: string;
    warehouse_id?: number;
    type?: string;
    business_id?: number;
}

export interface AdjustStockDTO {
    product_id: string;
    warehouse_id: number;
    location_id?: number | null;
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
    quantity: number;
    reason?: string;
    notes?: string;
}

// ============================================
// RESPONSES
// ============================================

export interface InventoryListResponse {
    data: InventoryLevel[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface MovementListResponse {
    data: StockMovement[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface MovementTypeListResponse {
    data: MovementType[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

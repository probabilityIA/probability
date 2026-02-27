// ============================================
// ENTIDADES
// ============================================

export interface Warehouse {
    id: number;
    business_id: number;
    name: string;
    code: string;
    address: string;
    city: string;
    state: string;
    country: string;
    zip_code: string;
    phone: string;
    contact_name: string;
    contact_email: string;
    is_active: boolean;
    is_default: boolean;
    is_fulfillment: boolean;
    created_at: string;
    updated_at: string;
}

export interface WarehouseDetail extends Warehouse {
    locations: WarehouseLocation[];
}

export interface WarehouseLocation {
    id: number;
    warehouse_id: number;
    name: string;
    code: string;
    type: string;
    is_active: boolean;
    is_fulfillment: boolean;
    capacity: number | null;
    created_at: string;
    updated_at: string;
}

// ============================================
// DTOs
// ============================================

export interface CreateWarehouseDTO {
    name: string;
    code: string;
    address?: string;
    city?: string;
    state?: string;
    country?: string;
    zip_code?: string;
    phone?: string;
    contact_name?: string;
    contact_email?: string;
    is_default?: boolean;
    is_fulfillment?: boolean;
}

export interface UpdateWarehouseDTO {
    name: string;
    code: string;
    address?: string;
    city?: string;
    state?: string;
    country?: string;
    zip_code?: string;
    phone?: string;
    contact_name?: string;
    contact_email?: string;
    is_active?: boolean;
    is_default?: boolean;
    is_fulfillment?: boolean;
}

export interface GetWarehousesParams {
    page?: number;
    page_size?: number;
    search?: string;
    is_active?: boolean;
    is_fulfillment?: boolean;
    business_id?: number;
}

export interface CreateLocationDTO {
    name: string;
    code: string;
    type?: string;
    is_fulfillment?: boolean;
    capacity?: number | null;
}

export interface UpdateLocationDTO {
    name: string;
    code: string;
    type?: string;
    is_active?: boolean;
    is_fulfillment?: boolean;
    capacity?: number | null;
}

// ============================================
// RESPONSES
// ============================================

export interface WarehousesListResponse {
    data: Warehouse[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

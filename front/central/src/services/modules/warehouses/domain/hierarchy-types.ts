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

export interface WarehouseLocationFlags {
    is_picking?: boolean;
    is_bulk?: boolean;
    is_quarantine?: boolean;
    is_damaged?: boolean;
    is_returns?: boolean;
    is_cross_dock?: boolean;
    is_hazmat?: boolean;
}

export interface Zone {
    id: number;
    warehouse_id: number;
    business_id: number;
    code: string;
    name: string;
    purpose: string;
    color_hex: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface Aisle {
    id: number;
    zone_id: number;
    business_id: number;
    code: string;
    name: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface Rack {
    id: number;
    aisle_id: number;
    business_id: number;
    code: string;
    name: string;
    levels_count: number;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface RackLevel {
    id: number;
    rack_id: number;
    business_id: number;
    code: string;
    ordinal: number;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface TreePosition {
    id: number;
    code: string;
    name: string;
    type: string;
    is_active: boolean;
    priority: number;
}

export interface TreeLevel extends RackLevel {
    positions: TreePosition[];
}

export interface TreeRack extends Rack {
    levels: TreeLevel[];
}

export interface TreeAisle extends Aisle {
    racks: TreeRack[];
}

export interface TreeZone extends Zone {
    aisles: TreeAisle[];
}

export interface WarehouseTree {
    warehouse_id: number;
    zones: TreeZone[];
}

export interface CreateZoneDTO {
    warehouse_id: number;
    code: string;
    name: string;
    purpose?: string;
    color_hex?: string;
    is_active?: boolean;
}

export interface UpdateZoneDTO {
    code?: string;
    name?: string;
    purpose?: string;
    color_hex?: string;
    is_active?: boolean;
}

export interface CreateAisleDTO {
    zone_id: number;
    code: string;
    name: string;
    is_active?: boolean;
}

export interface UpdateAisleDTO {
    code?: string;
    name?: string;
    is_active?: boolean;
}

export interface CreateRackDTO {
    aisle_id: number;
    code: string;
    name: string;
    levels_count?: number;
    is_active?: boolean;
}

export interface UpdateRackDTO {
    code?: string;
    name?: string;
    levels_count?: number;
    is_active?: boolean;
}

export interface CreateRackLevelDTO {
    rack_id: number;
    code: string;
    ordinal?: number;
    is_active?: boolean;
}

export interface UpdateRackLevelDTO {
    code?: string;
    ordinal?: number;
    is_active?: boolean;
}

export interface CubingCheckResult {
    fits: boolean;
    reason?: string;
    weight_needed_kg: number;
    weight_max_kg: number;
    volume_needed_cm3: number;
    volume_max_cm3: number;
    occupied_qty: number;
}

export interface ValidateCubingInput {
    product_id: string;
    location_id: number;
    quantity: number;
}

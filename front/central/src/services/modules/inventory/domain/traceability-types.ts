import { PaginatedResponse, PaginationParams } from './types';

export interface InventoryLot {
    id: number;
    business_id: number;
    product_id: string;
    lot_code: string;
    manufacture_date: string | null;
    expiration_date: string | null;
    received_at: string | null;
    supplier_id: number | null;
    status: string;
    created_at: string;
    updated_at: string;
}

export interface InventorySerial {
    id: number;
    business_id: number;
    product_id: string;
    serial_number: string;
    lot_id: number | null;
    current_location_id: number | null;
    current_state_id: number | null;
    received_at: string | null;
    sold_at: string | null;
    created_at: string;
    updated_at: string;
}

export interface InventoryState {
    id: number;
    code: string;
    name: string;
    description: string;
    is_terminal: boolean;
}

export interface UnitOfMeasure {
    id: number;
    code: string;
    name: string;
    type: string;
    is_active: boolean;
}

export interface ProductUoM {
    id: number;
    product_id: string;
    uom_id: number;
    uom_code: string;
    uom_name: string;
    conversion_factor: number;
    is_base: boolean;
    barcode: string;
    is_active: boolean;
}

export interface ConvertUoMResult {
    from_uom_code: string;
    to_uom_code: string;
    input_quantity: number;
    converted_quantity: number;
    base_unit_quantity: number;
    base_uom_code: string;
}

export interface GetLotsParams extends PaginationParams {
    product_id?: string;
    status?: string;
    expiring_in_days?: number;
    business_id?: number;
}

export interface GetSerialsParams extends PaginationParams {
    product_id?: string;
    lot_id?: number;
    state_id?: number;
    location_id?: number;
    business_id?: number;
}

export interface CreateLotDTO {
    product_id: string;
    lot_code: string;
    manufacture_date?: string | null;
    expiration_date?: string | null;
    received_at?: string | null;
    supplier_id?: number | null;
    status?: string;
}

export interface UpdateLotDTO {
    lot_code?: string;
    manufacture_date?: string | null;
    expiration_date?: string | null;
    received_at?: string | null;
    supplier_id?: number | null;
    status?: string;
}

export interface CreateSerialDTO {
    product_id: string;
    serial_number: string;
    lot_id?: number | null;
    location_id?: number | null;
    state_code?: string;
}

export interface UpdateSerialDTO {
    lot_id?: number | null;
    location_id?: number | null;
    state_code?: string;
}

export interface ChangeInventoryStateDTO {
    level_id: number;
    from_state_code: string;
    to_state_code: string;
    quantity: number;
    reason?: string;
}

export interface CreateProductUoMDTO {
    uom_code: string;
    conversion_factor: number;
    is_base: boolean;
    barcode?: string;
}

export interface ConvertUoMInput {
    product_id: string;
    from_uom_code: string;
    to_uom_code: string;
    quantity: number;
}

export type LotListResponse = PaginatedResponse<InventoryLot>;
export type SerialListResponse = PaginatedResponse<InventorySerial>;

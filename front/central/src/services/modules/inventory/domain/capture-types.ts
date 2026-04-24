import { PaginatedResponse, PaginationParams } from './types';

export interface LicensePlate {
    id: number;
    business_id: number;
    code: string;
    lpn_type: string;
    current_location_id: number | null;
    status: string;
    created_at: string;
    updated_at: string;
    Lines?: LicensePlateLine[];
}

export interface LicensePlateLine {
    id: number;
    lpn_id: number;
    business_id: number;
    product_id: string;
    lot_id: number | null;
    serial_id: number | null;
    qty: number;
    created_at: string;
    updated_at: string;
}

export interface ScanResolution {
    code: string;
    code_type: string;
    matched_id?: number | null;
    product_id?: string;
    location_id?: number | null;
    lot_id?: number | null;
    serial_id?: number | null;
    lpn_id?: number | null;
    suggested?: string;
    data?: Record<string, any>;
}

export interface ScanEvent {
    id: number;
    business_id: number;
    user_id: number | null;
    device_id: string;
    scanned_code: string;
    code_type: string;
    action: string;
    scanned_at: string;
    created_at: string;
}

export interface ScanResult {
    resolved: boolean;
    resolution?: ScanResolution;
    event?: ScanEvent;
}

export interface InventorySyncLog {
    id: number;
    business_id: number;
    integration_id: number | null;
    direction: string;
    payload_hash: string;
    status: string;
    error: string;
    synced_at: string | null;
    created_at: string;
}

export interface InboundSyncResult {
    log: InventorySyncLog;
    duplicate: boolean;
}

export interface CreateLPNDTO {
    code: string;
    lpn_type?: string;
    location_id?: number | null;
}

export interface UpdateLPNDTO {
    code?: string;
    lpn_type?: string;
    location_id?: number | null;
    status?: string;
}

export interface AddToLPNInput {
    product_id: string;
    lot_id?: number | null;
    serial_id?: number | null;
    qty: number;
}

export interface MoveLPNInput {
    new_location_id: number;
}

export interface MergeLPNInput {
    target_lpn_id: number;
}

export interface ScanInput {
    code: string;
    device_id?: string;
    action?: string;
}

export interface InboundSyncInput {
    payload: Record<string, any>;
}

export interface GetLPNsParams extends PaginationParams {
    lpn_type?: string;
    status?: string;
    location_id?: number;
    business_id?: number;
}

export interface GetSyncLogsParams extends PaginationParams {
    integration_id?: number;
    direction?: string;
    status?: string;
    business_id?: number;
}

export type LicensePlateListResponse = PaginatedResponse<LicensePlate>;
export type InventorySyncLogListResponse = PaginatedResponse<InventorySyncLog>;

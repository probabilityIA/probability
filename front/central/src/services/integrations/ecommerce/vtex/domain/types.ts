export interface VTEXWarehouseMappingConfig {
    internal_warehouse_id: number;
    vtex_warehouse_id: string;
}

export interface VTEXConfig {
    account_name?: string;
    is_seller?: boolean;
    inventory_sync_enabled?: boolean;
    status_sync_enabled?: boolean;
    vtex_warehouse_mappings?: VTEXWarehouseMappingConfig[];
}

export interface VTEXCredentials {
    app_key: string;
    app_token: string;
}

export interface VTEXIntegrationData {
    name: string;
    config: VTEXConfig;
    credentials: VTEXCredentials;
    is_active: boolean;
}

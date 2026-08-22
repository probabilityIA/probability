export interface TiendanubeConfig {
    cod_includes_shipping?: boolean;
    store_id?: string;
    is_testing?: boolean;
    inventory_sync_enabled?: boolean;
    status_sync_enabled?: boolean;
    inventory_single_warehouse_id?: number;
}

export interface TiendanubeCredentials {
    access_token: string;
}

export interface TiendanubeIntegrationData {
    name: string;
    config: TiendanubeConfig;
    credentials: TiendanubeCredentials;
    is_active: boolean;
}

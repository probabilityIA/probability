export interface WooCommerceConfig {
    store_url: string;
    free_shipping_enabled?: boolean;
    free_shipping_min?: number;
    inventory_sync_enabled?: boolean;
    inventory_warehouse_mode?: 'single' | 'sum';
    inventory_single_warehouse_id?: number;
    inventory_warehouse_ids?: number[];
}

export interface WooCommerceCredentials {
    consumer_key: string;
    consumer_secret: string;
}

export interface WooCommerceIntegrationData {
    name: string;
    config: WooCommerceConfig;
    credentials: WooCommerceCredentials;
    is_active: boolean;
}

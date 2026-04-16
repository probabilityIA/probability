// WooCommerce Integration Configuration Types

export interface WooCommerceConfig {
    store_url: string;          // URL de la tienda (ej: https://mitienda.com)
}

export interface WooCommerceCredentials {
    consumer_key: string;       // Consumer Key de la API REST
    consumer_secret: string;    // Consumer Secret de la API REST
}

export interface WooCommerceIntegrationData {
    name: string;
    config: WooCommerceConfig;
    credentials: WooCommerceCredentials;
    is_active: boolean;
}

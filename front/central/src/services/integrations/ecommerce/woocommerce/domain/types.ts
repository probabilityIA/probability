// WooCommerce Integration Configuration Types

export interface WooCommerceConfig {
    store_url: string;          // URL de la tienda (ej: https://mitienda.com)
    free_shipping_enabled?: boolean;  // Habilita envio gratis por monto minimo
    free_shipping_min?: number;       // Monto minimo de la orden para envio gratis
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

// Tiendanube Ecommerce Integration Configuration Types

export interface TiendanubeConfig {
    store_id?: string;        // ID de la tienda en Tiendanube
}

export interface TiendanubeCredentials {
    access_token: string;     // Token de acceso OAuth
}

export interface TiendanubeIntegrationData {
    name: string;
    config: TiendanubeConfig;
    credentials: TiendanubeCredentials;
    is_active: boolean;
}

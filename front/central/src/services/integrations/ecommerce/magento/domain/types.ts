// Magento / Adobe Commerce Integration Configuration Types

export interface MagentoConfig {
    store_url?: string;       // URL de la tienda Magento
}

export interface MagentoCredentials {
    access_token: string;     // Bearer token de la API REST
}

export interface MagentoIntegrationData {
    name: string;
    config: MagentoConfig;
    credentials: MagentoCredentials;
    is_active: boolean;
}

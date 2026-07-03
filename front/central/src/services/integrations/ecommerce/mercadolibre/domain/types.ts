export interface MercadoLibreConfig {
    seller_id?: string;
}

export interface MercadoLibreCredentials {
    client_id: string;
    client_secret: string;
    access_token: string;
    refresh_token: string;
    seller_id: string;
}

export interface MercadoLibreIntegrationData {
    name: string;
    config: MercadoLibreConfig;
    credentials: MercadoLibreCredentials;
    is_active: boolean;
}

export interface AmazonConfig {
    marketplace_id?: string;   // ID del marketplace (ej: A1AM78C64UM0Y8 para Mexico)
    region?: string;           // Region (na, eu, fe)
}

export interface AmazonCredentials {
    seller_id: string;         // ID del vendedor en Amazon
    refresh_token: string;     // Refresh token del SP-API OAuth
}

export interface AmazonIntegrationData {
    name: string;
    config: AmazonConfig;
    credentials: AmazonCredentials;
    is_active: boolean;
}

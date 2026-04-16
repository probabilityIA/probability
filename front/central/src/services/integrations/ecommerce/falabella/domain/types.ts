// Falabella Integration Configuration Types

export interface FalabellaConfig {
    // Config especifica se agrega cuando se implemente backend completo
}

export interface FalabellaCredentials {
    api_key: string;           // API Key del Seller Center
    user_id: string;           // User ID del vendedor
}

export interface FalabellaIntegrationData {
    name: string;
    config: FalabellaConfig;
    credentials: FalabellaCredentials;
    is_active: boolean;
}

export interface ExitoConfig {
    // Config especifica se agrega cuando se implemente backend completo
}

export interface ExitoCredentials {
    api_key: string;           // API Key del marketplace
    seller_id: string;         // ID del vendedor
}

export interface ExitoIntegrationData {
    name: string;
    config: ExitoConfig;
    credentials: ExitoCredentials;
    is_active: boolean;
}

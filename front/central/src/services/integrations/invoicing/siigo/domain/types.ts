// Siigo Integration Configuration Types

export interface SiigoConfig {
    // Config especifica se agrega cuando se implemente backend completo
}

export interface SiigoCredentials {
    username: string;      // Usuario API de Siigo
    access_key: string;    // Clave de acceso API
    account_id: string;    // ID de suscripcion / Account ID
    partner_id: string;    // Partner ID (header)
    base_url?: string;     // URL API (opcional)
}

export interface SiigoIntegrationData {
    name: string;
    config: SiigoConfig;
    credentials: SiigoCredentials;
    is_active: boolean;
}

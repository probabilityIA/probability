// Alegra Integration Configuration Types

export interface AlegraConfig {
    // Config especifica se agrega cuando se implemente backend completo
}

export interface AlegraCredentials {
    email: string;      // Email de la cuenta Alegra
    token: string;      // Token API / API Key
    base_url?: string;  // URL API (opcional)
}

export interface AlegraIntegrationData {
    name: string;
    config: AlegraConfig;
    credentials: AlegraCredentials;
    is_active: boolean;
}

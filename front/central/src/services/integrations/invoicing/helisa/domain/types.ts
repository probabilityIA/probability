// Helisa Integration Configuration Types

export interface HelisaConfig {
    // Config especifica se agrega cuando se implemente backend completo
}

export interface HelisaCredentials {
    username: string;    // Usuario
    password: string;    // Contrasena
    company_id: string;  // ID de empresa
    base_url?: string;   // URL API (opcional)
}

export interface HelisaIntegrationData {
    name: string;
    config: HelisaConfig;
    credentials: HelisaCredentials;
    is_active: boolean;
}

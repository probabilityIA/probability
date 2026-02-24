// World Office Integration Configuration Types

export interface WorldOfficeConfig {
    // Config especifica se agrega cuando se implemente backend completo
}

export interface WorldOfficeCredentials {
    username: string;      // Usuario
    password: string;      // Contrasena
    company_code: string;  // Codigo de empresa
    base_url?: string;     // URL API (opcional)
}

export interface WorldOfficeIntegrationData {
    name: string;
    config: WorldOfficeConfig;
    credentials: WorldOfficeCredentials;
    is_active: boolean;
}

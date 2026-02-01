// Softpymes Integration Configuration Types

export interface SoftpymesConfig {
    company_nit: string;          // NIT de la empresa
    company_name: string;         // Nombre de la empresa
    referer: string;              // Identificación de instancia del cliente (requerido para header Referer)
    api_url: string;              // URL de la API de Softpymes
    test_mode?: boolean;          // Modo de pruebas
}

export interface SoftpymesCredentials {
    api_key: string;              // API Key de Softpymes
    api_secret: string;           // API Secret de Softpymes
}

export interface SoftpymesIntegrationData {
    name: string;
    config: SoftpymesConfig;
    credentials: SoftpymesCredentials;
    is_active: boolean;
}

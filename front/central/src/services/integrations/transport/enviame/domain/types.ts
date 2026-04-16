// Enviame Integration Configuration Types

export interface EnviameConfig {
    base_url?: string;            // URL base de la API de Enviame (opcional, tiene default)
}

export interface EnviameCredentials {
    api_key: string;              // API Key de Enviame
}

export interface EnviameIntegrationData {
    name: string;
    config: EnviameConfig;
    credentials: EnviameCredentials;
    is_active: boolean;
}

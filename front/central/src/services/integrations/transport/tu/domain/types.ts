// TU Integration Configuration Types

export interface TuConfig {
    base_url?: string;            // URL base de la API de TU (opcional, tiene default)
}

export interface TuCredentials {
    api_key: string;              // API Key de TU
}

export interface TuIntegrationData {
    name: string;
    config: TuConfig;
    credentials: TuCredentials;
    is_active: boolean;
}

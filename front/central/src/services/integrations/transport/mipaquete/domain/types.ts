// MiPaquete Integration Configuration Types

export interface MiPaqueteConfig {
    base_url?: string;            // URL base de la API de MiPaquete (opcional, tiene default)
}

export interface MiPaqueteCredentials {
    api_key: string;              // API Key de MiPaquete
}

export interface MiPaqueteIntegrationData {
    name: string;
    config: MiPaqueteConfig;
    credentials: MiPaqueteCredentials;
    is_active: boolean;
}

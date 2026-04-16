// EnvioClick Integration Configuration Types

export interface EnvioClickConfig {
    use_platform_token?: boolean; // true = usar token compartido de la plataforma (no requiere api_key propia)
    base_url_test?: string;       // URL de pruebas para esta integración (override del base_url del tipo)
}

export interface EnvioClickCredentials {
    api_key?: string;             // API Key de EnvioClick (vacío cuando use_platform_token=true)
}

export interface EnvioClickIntegrationData {
    name: string;
    config: EnvioClickConfig;
    credentials: EnvioClickCredentials;
    is_active: boolean;
}

// VTEX Integration Configuration Types

export interface VTEXConfig {
    account_name?: string;    // Nombre de la cuenta VTEX
    environment?: string;     // Ambiente (produccion, sandbox)
}

export interface VTEXCredentials {
    app_key: string;          // X-VTEX-API-AppKey
    app_token: string;        // X-VTEX-API-AppToken
}

export interface VTEXIntegrationData {
    name: string;
    config: VTEXConfig;
    credentials: VTEXCredentials;
    is_active: boolean;
}

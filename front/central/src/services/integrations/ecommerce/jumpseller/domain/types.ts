export interface JumpsellerCredentials {
    api_key: string;
    api_secret: string;
}

export interface JumpsellerConfig {
    inventory_sync_enabled?: boolean;
    status_sync_enabled?: boolean;
}

export interface JumpsellerIntegrationData {
    name: string;
    config: JumpsellerConfig;
    credentials: JumpsellerCredentials;
    is_active: boolean;
}

export interface ShipitConfig {
    base_url?: string;
    base_url_test?: string;
}

export interface ShipitCredentials {
    email: string;
    access_token: string;
}

export interface ShipitIntegrationData {
    name: string;
    config: ShipitConfig;
    credentials: ShipitCredentials;
    is_active: boolean;
}

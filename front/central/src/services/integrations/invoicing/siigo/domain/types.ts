export interface SiigoConfig {}

export interface SiigoCredentials {
    username: string;
    access_key: string;
    account_id?: string;
    partner_id: string;
}

export interface SiigoIntegrationData {
    name: string;
    config: SiigoConfig;
    credentials: SiigoCredentials;
    is_active: boolean;
}

export interface SiigoWebhookInfo {
    id: string;
    application_id: string;
    url: string;
    topic: string;
    company_key: string;
    active: boolean;
    created_at: string;
}

export interface PayUConfig {
    account_id?: string;   // ID de cuenta PayU
    merchant_id?: string;  // ID de comercio PayU
}

export interface PayUCredentials {
    api_key: string;
    api_login: string;
    environment: 'sandbox' | 'production';
}

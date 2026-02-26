export interface NequiConfig {
    phone_code?: string; // Prefijo telefónico por defecto (ej: "+57")
}

export interface NequiCredentials {
    api_key: string;
    environment: 'sandbox' | 'production';
}

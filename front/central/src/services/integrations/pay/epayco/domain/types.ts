export interface EPaycoConfig {
    // Sin configuración adicional necesaria para ePayco
}

export interface EPaycoCredentials {
    customer_id: string;
    key: string;
    environment: 'test' | 'production';
}

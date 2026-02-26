export interface StripeConfig {
    // Sin configuración adicional necesaria para Stripe
}

export interface StripeCredentials {
    secret_key: string;
    environment: 'test' | 'live';
}

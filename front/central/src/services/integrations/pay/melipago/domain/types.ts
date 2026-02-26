export interface MeliPagoConfig {
    // Sin configuración adicional necesaria para MercadoPago
}

export interface MeliPagoCredentials {
    access_token: string;
    environment: 'sandbox' | 'production';
}

// Softpymes Integration Configuration Types

export interface SoftpymesConfig {
    company_nit: string;          // NIT de la empresa
    company_name: string;         // Nombre de la empresa
    referer: string;              // Identificación de instancia del cliente (requerido para header Referer)
    default_customer_nit?: string; // NIT por defecto para clientes sin DNI (consumidor final: 222222222222)
    resolution_id?: number;       // ID de resolución de facturación (obtenido desde Softpymes /resolutions)
    branch_code?: string;         // Código de sucursal (default: "001")
    customer_branch_code?: string; // Código de sucursal del cliente (default: "001", requerido por Softpymes)
    seller_nit?: string;          // NIT del vendedor (opcional, usa referer por defecto)
    // Configuración de recibo de caja (para registrar pagos y mover cuenta contable)
    send_cash_receipt?: boolean;   // true = enviar recibo de caja después de crear la factura
    payment_type?: string;         // Tipo de pago: "EF"=Efectivo, "TR"=Transferencia, "TC"=T.Crédito, "TD"=T.Débito
    payment_account_number?: string; // Número de cuenta bancaria (requerido para tipo TR)
}

export interface SoftpymesCredentials {
    api_key: string;              // API Key de Softpymes
    api_secret: string;           // API Secret de Softpymes
}

export interface SoftpymesIntegrationData {
    name: string;
    config: SoftpymesConfig;
    credentials: SoftpymesCredentials;
    is_active: boolean;
}

// Factus Integration Configuration Types

export interface FactusConfig {
    numbering_range_id: number;           // ID del rango de numeración en Factus (requerido)
    default_tax_rate?: string;            // Tasa de IVA por defecto (default: "19.00")
    payment_form?: string;                // Forma de pago: 1=Contado, 2=Crédito (default: "1")
    payment_method_code?: string;         // Método de pago: 10=Efectivo, 42=Consignación (default: "10")
    legal_organization_id?: string;       // Tipo de organización legal: 1=Jurídica, 2=Natural (default: "2")
    tribute_id?: string;                  // Régimen tributario DIAN (default: "21" = No responsable IVA)
    identification_document_id?: string;  // Tipo de documento: 3=Cédula, 13=NIT (default: "3")
    municipality_id?: string;             // ID municipio del cliente (opcional)
}

export interface FactusCredentials {
    client_id: string;     // Client ID OAuth2
    client_secret: string; // Client Secret OAuth2
    username: string;      // Email de la cuenta Factus
    password: string;      // Contraseña de la cuenta Factus
    api_url?: string;      // URL base de la API (opcional, default: https://api.factus.com.co)
}

export interface FactusIntegrationData {
    name: string;
    config: FactusConfig;
    credentials: FactusCredentials;
    is_active: boolean;
}

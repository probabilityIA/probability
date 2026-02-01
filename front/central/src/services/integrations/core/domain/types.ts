export interface IntegrationConfig {
    [key: string]: any;
}

export interface IntegrationCredentials {
    [key: string]: any;
}

export interface IntegrationTypeInfo {
    id: number;
    name: string;
    code: string;
    image_url?: string; // URL completa de la imagen del logo
}

export interface Integration {
    id: number;
    name: string;
    code: string;
    integration_type_id: number;
    type: 'whatsapp' | 'shopify' | 'mercado_libre' | string;
    category: 'internal' | 'external' | string;
    business_id: number | null;
    store_id?: string; // Identificador externo (ej: shop domain para Shopify)
    is_active: boolean;
    is_default: boolean;
    config: IntegrationConfig;
    credentials?: IntegrationCredentials; // Solo se incluye cuando se solicita para edición
    description?: string;
    created_by_id: number;
    updated_by_id: number | null;
    created_at: string;
    updated_at: string;
    integration_type?: IntegrationTypeInfo; // Información del tipo de integración si está cargado
}

export interface CreateIntegrationDTO {
    name: string;
    code: string;
    integration_type_id: number;
    type?: string;
    category: string;
    business_id: number | null;
    store_id?: string;
    is_active?: boolean;
    is_default?: boolean;
    config?: IntegrationConfig;
    credentials?: IntegrationCredentials;
    description?: string;
}

export interface UpdateIntegrationDTO {
    name?: string;
    code?: string;
    store_id?: string;
    is_active?: boolean;
    is_default?: boolean;
    config?: IntegrationConfig;
    credentials?: IntegrationCredentials;
    description?: string;
}

export interface GetIntegrationsParams {
    page?: number;
    page_size?: number;
    type?: string;
    category?: string;              // Legacy string category
    category_id?: number;           // NEW - Filter by category ID
    business_id?: number;
    is_active?: boolean;
    search?: string;
}

export interface PaginatedResponse<T> {
    success: boolean;
    message: string;
    data: T[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface SingleResponse<T> {
    success: boolean;
    message: string;
    data: T;
}

export interface ActionResponse {
    success: boolean;
    message: string;
    error?: string;
}

export interface IntegrationType {
    id: number;
    name: string;
    code: string;
    description?: string;
    icon?: string;
    image_url?: string; // URL completa de la imagen del logo
    category: 'internal' | 'external' | string; // Legacy field, kept for backward compatibility
    category_id?: number;             // NEW - FK to IntegrationCategory
    integration_category?: IntegrationCategory;  // NEW - Populated category object
    is_active: boolean;
    config_schema?: any;
    credentials_schema?: any;
    setup_instructions?: string;
    created_at: string;
    updated_at: string;
}

export interface CreateIntegrationTypeDTO {
    name: string;
    code?: string;
    description?: string;
    icon?: string;
    category: string;
    is_active?: boolean;
    config_schema?: any;
    credentials_schema?: any;
    setup_instructions?: string;
    image_file?: File; // Archivo de imagen para subir
}

export interface UpdateIntegrationTypeDTO {
    name?: string;
    code?: string;
    description?: string;
    icon?: string;
    category?: string;
    is_active?: boolean;
    config_schema?: any;
    credentials_schema?: any;
    setup_instructions?: string;
    image_file?: File; // Archivo de imagen para subir
    remove_image?: boolean; // Flag para eliminar la imagen existente
}

// Información del webhook para configurar en plataformas externas
export interface WebhookInfo {
    url: string;
    method: string;
    description: string;
    events?: string[];
}

export interface WebhookResponse {
    success: boolean;
    data: WebhookInfo;
}

// Información de un webhook de Shopify (desde la API)
export interface ShopifyWebhookInfo {
    id: string;
    address: string;
    topic: string;
    format: string;
    created_at: string;
    updated_at: string;
}

// Respuesta al listar webhooks
export interface ListWebhooksResponse {
    success: boolean;
    data: ShopifyWebhookInfo[];
}

// Respuesta al eliminar webhook
export interface DeleteWebhookResponse {
    success: boolean;
    message: string;
}

// Respuesta al verificar webhooks existentes
export interface VerifyWebhooksResponse {
    success: boolean;
    data: ShopifyWebhookInfo[];
    message: string;
}

// Datos del resultado de crear webhooks
export interface CreateWebhookResponseData {
    existing_webhooks: ShopifyWebhookInfo[];
    deleted_webhooks: ShopifyWebhookInfo[];
    created_webhooks: string[];
    webhook_url: string;
}

// Respuesta al crear webhooks
export interface CreateWebhookResponse {
    success: boolean;
    data: CreateWebhookResponseData;
    message: string;
}

// Parámetros para sincronización de órdenes
export interface SyncOrdersParams {
    created_at_min?: string;  // Formato: YYYY-MM-DD o RFC3339
    created_at_max?: string;  // Formato: YYYY-MM-DD o RFC3339
    status?: string;          // any, open, closed, cancelled
    financial_status?: string; // any, paid, pending, refunded, etc.
    fulfillment_status?: string; // any, shipped, partial, unshipped, etc.
}

// ============================================
// Simple Types para Dropdowns/Selectores
// ============================================

export interface IntegrationSimple {
    id: number;
    name: string;
    type: string;
    business_id: number | null;
    is_active: boolean;
}

export interface IntegrationsSimpleResponse {
    success: boolean;
    message: string;
    data: IntegrationSimple[];
}

// ============================================
// Integration Categories
// ============================================

export interface IntegrationCategory {
    id: number;
    code: string;                    // 'ecommerce', 'invoicing', 'messaging', 'system'
    name: string;                    // 'E-commerce', 'Facturación', 'Mensajería', 'Sistema'
    description?: string;
    icon?: string;                   // Icon name (heroicons)
    color?: string;                  // Tailwind color class or hex
    display_order: number;           // Display order in UI
    parent_category_id?: number;     // For nested categories (future)
    is_active: boolean;
    is_visible: boolean;
    created_at: string;
    updated_at: string;
}

export interface IntegrationCategoriesResponse {
    success: boolean;
    message: string;
    data: IntegrationCategory[];
}

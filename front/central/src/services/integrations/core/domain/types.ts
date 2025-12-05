export interface IntegrationConfig {
    [key: string]: any;
}

export interface IntegrationCredentials {
    [key: string]: any;
}

export interface Integration {
    id: number;
    name: string;
    code: string;
    integration_type_id: number;
    type: 'whatsapp' | 'shopify' | 'mercado_libre' | string;
    category: 'internal' | 'external' | string;
    business_id: number | null;
    is_active: boolean;
    is_default: boolean;
    config: IntegrationConfig;
    description?: string;
    created_by_id: number;
    updated_by_id: number | null;
    created_at: string;
    updated_at: string;
}

export interface CreateIntegrationDTO {
    name: string;
    code: string;
    integration_type_id: number;
    type?: string;
    category: string;
    business_id: number | null;
    is_active?: boolean;
    is_default?: boolean;
    config?: IntegrationConfig;
    credentials?: IntegrationCredentials;
    description?: string;
}

export interface UpdateIntegrationDTO {
    name?: string;
    code?: string;
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
    category?: string;
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
    category: 'internal' | 'external' | string;
    is_active: boolean;
    config_schema?: any;
    credentials_schema?: any;
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
}

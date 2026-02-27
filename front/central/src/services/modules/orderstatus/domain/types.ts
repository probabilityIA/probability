export interface OrderStatusInfo {
    id: number;
    code: string;
    name: string;
    description?: string;
    category?: string;
    color?: string;
    priority?: number;
    is_active?: boolean;
}

export interface CreateOrderStatusDTO {
    code: string;
    name: string;
    description?: string;
    category?: string;
    color?: string;
    priority?: number;
    is_active?: boolean;
}

export interface UpdateOrderStatusDTO {
    code: string;
    name: string;
    description?: string;
    category?: string;
    color?: string;
    priority?: number;
    is_active?: boolean;
}

export interface IntegrationTypeInfo {
    id: number;
    code: string;
    name: string;
    image_url?: string;
}

export interface OrderStatusMapping {
    id: number;
    integration_type_id: number;
    integration_type?: IntegrationTypeInfo;
    original_status: string;
    order_status_id: number;
    order_status?: OrderStatusInfo;
    is_active: boolean;
    description: string;
    created_at: string;
    updated_at: string;
}

export interface PaginatedResponse<T> {
    success: boolean;
    message?: string;
    data: T[];
    total: number;
    page?: number;
    page_size?: number;
    total_pages?: number;
}

export interface SingleResponse<T> {
    success: boolean;
    message?: string;
    data: T;
}

export interface ActionResponse {
    success: boolean;
    message: string;
    error?: string;
}

export interface GetOrderStatusMappingsParams {
    page?: number;
    page_size?: number;
    integration_type_id?: number;
    is_active?: boolean;
}

export interface CreateOrderStatusMappingDTO {
    integration_type_id: number;
    original_status: string;
    order_status_id: number;
    description?: string;
}

export interface UpdateOrderStatusMappingDTO {
    original_status: string;
    order_status_id: number;
    description?: string;
}

// ============================================
// Simple Types para Dropdowns/Selectores
// ============================================

export interface OrderStatusSimple {
    id: number;
    name: string;
    code: string;
    is_active: boolean;
}

export interface OrderStatusesSimpleResponse {
    success: boolean;
    message: string;
    data: OrderStatusSimple[];
}

// ============================================
// Channel Statuses - Estados nativos por canal
// ============================================

export interface EcommerceIntegrationType {
    id: number;
    code: string;
    name: string;
    image_url?: string;
}

export interface ChannelStatusInfo {
    id: number;
    integration_type_id: number;
    integration_type?: EcommerceIntegrationType;
    code: string;
    name: string;
    description?: string;
    is_active: boolean;
    display_order: number;
}

export interface CreateChannelStatusDTO {
    integration_type_id: number;
    code: string;
    name: string;
    description?: string;
    is_active: boolean;
    display_order: number;
}

export interface UpdateChannelStatusDTO {
    code: string;
    name: string;
    description?: string;
    is_active: boolean;
    display_order: number;
}

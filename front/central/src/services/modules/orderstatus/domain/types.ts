export interface OrderStatusMapping {
    id: number;
    integration_type: string;
    original_status: string;
    mapped_status: string;
    is_active: boolean;
    priority: number;
    description: string;
    created_at: string;
    updated_at: string;
}

export interface PaginatedResponse<T> {
    success: boolean;
    message?: string;
    data: T[];
    total: number;
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
    integration_type?: string;
    is_active?: boolean;
}

export interface CreateOrderStatusMappingDTO {
    integration_type: string;
    original_status: string;
    mapped_status: string;
    priority?: number;
    description?: string;
}

export interface UpdateOrderStatusMappingDTO {
    original_status: string;
    mapped_status: string;
    priority?: number;
    description?: string;
}

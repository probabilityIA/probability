export interface Business {
    id: number;
    name: string;
    code?: string;
    business_type?: BusinessType;
    business_type_id: number;
    timezone?: string;
    address?: string;
    description?: string;
    logo_url?: string;
    primary_color?: string;
    secondary_color?: string;
    tertiary_color?: string;
    quaternary_color?: string;
    navbar_image_url?: string;
    custom_domain?: string;
    is_active: boolean;
    enable_delivery?: boolean;
    enable_pickup?: boolean;
    enable_reservations?: boolean;
    created_at?: string;
    updated_at?: string;
    phone?: string;
    email?: string;
    website?: string;
}

export interface BusinessType {
    id: number;
    name: string;
    code: string;
    description?: string;
    icon?: string;
    is_active: boolean;
    created_at?: string;
    updated_at?: string;
}

export interface ConfiguredResource {
    resource_id: number;
    resource_name: string;
    is_active: boolean;
}

export interface BusinessConfiguredResources {
    business_id: number;
    resources: ConfiguredResource[];
    total: number;
    active: number;
    inactive: number;
}

export interface Pagination {
    current_page: number;
    per_page: number;
    total: number;
    last_page: number;
    has_next: boolean;
    has_prev: boolean;
}

export interface PaginatedResponse<T> {
    success: boolean;
    message: string;
    data: T[];
    pagination: Pagination;
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

export interface GetBusinessesParams {
    page?: number;
    per_page?: number;
    name?: string;
    business_type_id?: number;
    is_active?: boolean;
}

export interface CreateBusinessDTO {
    name: string;
    code: string;
    business_type_id: number;
    timezone?: string;
    address?: string;
    description?: string;
    logo_file?: File;
    primary_color?: string;
    secondary_color?: string;
    tertiary_color?: string;
    quaternary_color?: string;
    navbar_image_file?: File;
    custom_domain?: string;
    is_active?: boolean;
    enable_delivery?: boolean;
    enable_pickup?: boolean;
    enable_reservations?: boolean;
}

export interface UpdateBusinessDTO extends Partial<CreateBusinessDTO> { }

export interface GetConfiguredResourcesParams {
    page?: number;
    per_page?: number;
    business_id?: number;
    business_type_id?: number;
}

export interface CreateBusinessTypeDTO {
    name: string;
    code: string;
    description?: string;
    icon?: string;
    is_active?: boolean;
}

export interface UpdateBusinessTypeDTO extends Partial<CreateBusinessTypeDTO> { }

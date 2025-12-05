export interface Resource {
    id: number;
    name: string;
    description?: string;
    business_type_id: number;
    business_type_name: string;
    created_at: string;
    updated_at: string;
}

export interface PaginatedResponse<T> {
    success: boolean;
    message: string;
    data: {
        resources: T[];
        total: number;
        page: number;
        page_size: number;
        total_pages: number;
    };
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

export interface GetResourcesParams {
    page?: number;
    page_size?: number;
    name?: string;
    description?: string;
    business_type_id?: number;
    sort_by?: 'name' | 'created_at' | 'updated_at';
    sort_order?: 'asc' | 'desc';
}

export interface CreateResourceDTO {
    name: string;
    description?: string;
    business_type_id?: number | null;
}

export interface UpdateResourceDTO {
    name: string;
    description: string;
    business_type_id?: number | null;
}

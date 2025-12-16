export interface Action {
    id: number;
    name: string;
    description?: string;
    created_at: string;
    updated_at: string;
}

export interface PaginatedResponse<T> {
    success: boolean;
    message: string;
    data: {
        actions: T[];
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

export interface GetActionsParams {
    page?: number;
    page_size?: number;
    name?: string;
    sort_by?: 'name' | 'created_at' | 'updated_at';
    sort_order?: 'asc' | 'desc';
}

export interface CreateActionDTO {
    name: string;
    description?: string;
}

export interface UpdateActionDTO {
    name: string;
    description?: string;
}

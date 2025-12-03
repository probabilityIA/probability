export interface Permission {
    id: number;
    name: string;
    code: string;
    description?: string;
    resource: string;
    action: string;
    resource_id: number;
    action_id: number;
    scope_id: number;
    scope_name: string;
    scope_code: string;
    business_type_id: number;
    business_type_name: string;
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
    data: T[];
    total: number; // The API example shows 'total' at the root for list response, not inside a pagination object like Business module. 
    // However, the Business module had a 'pagination' object. 
    // Let's look closely at the API Example for GET /permissions (Response 200 OK):
    // { "success": true, "data": [...], "total": 2 }
    // It seems different from the Business module which had "pagination": { ... }.
    // I will stick to the API Example provided for this module.
}

export interface SingleResponse<T> {
    success: boolean;
    data: T;
    message?: string;
}

export interface ActionResponse {
    success: boolean;
    message: string;
    error?: string;
}

export interface GetPermissionsParams {
    business_type_id?: number;
    name?: string;
    scope_id?: number;
}

export interface CreatePermissionDTO {
    name: string;
    code?: string;
    description?: string;
    resource_id: number;
    action_id: number;
    scope_id: number;
    business_type_id?: number | null;
}

export interface UpdatePermissionDTO {
    name: string;
    code: string;
    description: string;
    resource_id: number;
    action_id: number;
    scope_id: number;
    business_type_id?: number | null;
}

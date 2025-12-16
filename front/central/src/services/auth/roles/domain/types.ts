export interface Role {
    id: number;
    name: string;
    code: string;
    description: string;
    level: number;
    is_system: boolean;
    scope_id: number;
    scope_name: string;
    scope_code: string;
    business_type_id: number;
    business_type_name: string;
    created_at?: string;
    updated_at?: string;
}

export interface Permission {
    id: number;
    resource: string;
    action: string;
    description: string;
    scope_id: number;
    scope_name: string;
    scope_code: string;
}

export interface RolePermissionsResponse {
    success: boolean;
    message: string;
    role_id: number;
    role_name: string;
    permissions: Permission[];
    count: number;
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

export interface ListResponse<T> {
    success: boolean;
    data: T[];
    count: number;
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

export interface GetRolesParams {
    business_type_id?: number;
    scope_id?: number;
    is_system?: boolean;
    name?: string;
    level?: number;
    page?: number;
    page_size?: number;
}

export interface CreateRoleDTO {
    name: string;
    description: string;
    level: number;
    is_system: boolean;
    scope_id: number;
    business_type_id: number;
}

export type UpdateRoleDTO = Partial<CreateRoleDTO>;

export interface AssignPermissionsDTO {
    permission_ids: number[];
}

export interface AssignPermissionsResponse {
    success: boolean;
    message: string;
    role_id: number;
    permission_ids: number[];
}

export interface BusinessRoleAssignment {
    business_id: number;
    business_name: string;
    role_id: number;
    role_name: string;
}

export interface User {
    id: number;
    name: string;
    email: string;
    phone?: string;
    avatar_url?: string;
    is_active: boolean;
    is_super_user: boolean;
    last_login_at?: string;
    business_role_assignments: BusinessRoleAssignment[];
    created_at: string;
    updated_at: string;
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
    pagination: Pagination;
}

export interface SingleResponse<T> {
    success: boolean;
    data: T;
    message?: string;
    // For create user response which has extra fields
    email?: string;
    password?: string;
}

export interface ActionResponse {
    success: boolean;
    message: string;
    error?: string;
}

export interface GetUsersParams {
    page?: number;
    page_size?: number;
    name?: string;
    email?: string;
    phone?: string;
    user_ids?: string;
    is_active?: boolean;
    role_id?: number;
    business_id?: number;
    created_at?: string;
    sort_by?: string;
    sort_order?: 'asc' | 'desc';
}

export interface CreateUserDTO {
    name: string;
    email: string;
    phone?: string;
    is_active?: boolean;
    avatarFile?: File;
    business_ids?: number[];
}

export interface UpdateUserDTO {
    name?: string;
    email?: string;
    phone?: string;
    is_active?: boolean;
    remove_avatar?: boolean;
    avatarFile?: File;
    business_ids?: number[];
}

export interface RoleAssignment {
    business_id: number;
    role_id: number;
}

export interface AssignRolesDTO {
    assignments: RoleAssignment[];
}

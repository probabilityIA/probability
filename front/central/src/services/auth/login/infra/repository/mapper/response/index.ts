// Interfaces
export interface UserInfo {
    id: number;
    name: string;
    email: string;
    phone: string;
    avatar_url: string;
    is_active: boolean;
    last_login_at?: string;
}

export interface BusinessTypeInfo {
    id: number;
    name: string;
    code: string;
    description: string;
    icon: string;
}

export interface BusinessInfo {
    id: number;
    name: string;
    code: string;
    business_type_id: number;
    business_type: BusinessTypeInfo;
    timezone: string;
    address: string;
    description: string;
    logo_url: string;
    primary_color: string;
    secondary_color: string;
    tertiary_color: string;
    quaternary_color: string;
    navbar_image_url: string;
    custom_domain: string;
    is_active: boolean;
    enable_delivery: boolean;
    enable_pickup: boolean;
    enable_reservations: boolean;
}

export interface LoginResponse {
    user: UserInfo;
    token: string;
    require_password_change: boolean;
    businesses: BusinessInfo[];
    scope: string;
    is_super_admin: boolean;
}

export interface LoginSuccessResponse {
    success: boolean;
    data: LoginResponse;
}

export interface RoleInfo {
    id: number;
    name: string;
    description: string;
}

export interface ResourcePermissions {
    resource: string;
    actions: string[];
    active: boolean;
}

export interface UserRolesPermissionsResponse {
    is_super: boolean;
    business_id: number;
    business_name: string;
    business_type_id: number;
    business_type_name: string;
    role: RoleInfo;
    resources: ResourcePermissions[];
}

export interface UserRolesPermissionsSuccessResponse {
    success: boolean;
    data: UserRolesPermissionsResponse;
}

export interface ChangePasswordResponse {
    success: boolean;
    message: string;
}

export interface GeneratePasswordResponse {
    success: boolean;
    email: string;
    password: string;
    message: string;
}

export interface BusinessTokenResponse {
    token: string;
}

export interface GenerateBusinessTokenSuccessResponse {
    success: boolean;
    data: BusinessTokenResponse;
    message: string;
}

// Mappers
export const mapLoginResponse = (res: any): LoginSuccessResponse => {
    return res as LoginSuccessResponse;
};

export const mapUserRolesPermissionsResponse = (res: any): UserRolesPermissionsSuccessResponse => {
    return res as UserRolesPermissionsSuccessResponse;
};

export const mapChangePasswordResponse = (res: any): ChangePasswordResponse => {
    return res as ChangePasswordResponse;
};

export const mapGeneratePasswordResponse = (res: any): GeneratePasswordResponse => {
    return res as GeneratePasswordResponse;
};

export const mapGenerateBusinessTokenResponse = (res: any): GenerateBusinessTokenSuccessResponse => {
    return res as GenerateBusinessTokenSuccessResponse;
};

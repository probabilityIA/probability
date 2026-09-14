export interface AccessNavItem {
    key: string;
    label: string;
    route: string;
    section: string;
    description: string;
}

export interface AccessBusiness {
    id: number;
    name: string;
}

export interface AccessRole {
    id: number;
    code: string;
    name: string;
}

export interface Access {
    is_super: boolean;
    business: AccessBusiness | null;
    role: AccessRole;
    subscription: { status: string };
    modules: string[];
    permissions: string[];
    navigation: AccessNavItem[];
}

export interface AccessResult {
    success: boolean;
    data?: Access;
    error?: string;
}

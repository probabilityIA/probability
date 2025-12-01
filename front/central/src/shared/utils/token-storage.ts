export interface BusinessColors {
    primary?: string;
    secondary?: string;
    tertiary?: string;
    quaternary?: string;
}

export interface BusinessData {
    id: number;
    name: string;
    code: string;
    logo_url?: string;
    is_active?: boolean;
    primary_color?: string;
    secondary_color?: string;
    tertiary_color?: string;
    quaternary_color?: string;
}

export interface UserData {
    userId: string;
    name: string;
    email: string;
    role: string;
    avatarUrl?: string;
    is_super_admin?: boolean;
    scope?: string;
}

const KEYS = {
    SESSION_TOKEN: 'session_token',
    BUSINESS_TOKEN: 'business_token',
    USER_DATA: 'user_data',
    BUSINESSES_DATA: 'businesses_data',
    ACTIVE_BUSINESS_ID: 'active_business_id',
    PERMISSIONS: 'permissions',
    BUSINESS_COLORS: 'business_colors',
};

export const TokenStorage = {
    getSessionToken: (): string | null => {
        if (typeof window === 'undefined') return null;
        return localStorage.getItem(KEYS.SESSION_TOKEN);
    },

    setSessionToken: (token: string) => {
        if (typeof window === 'undefined') return;
        localStorage.setItem(KEYS.SESSION_TOKEN, token);
    },

    getBusinessToken: (): string | null => {
        if (typeof window === 'undefined') return null;
        return localStorage.getItem(KEYS.BUSINESS_TOKEN);
    },

    setBusinessToken: (token: string) => {
        if (typeof window === 'undefined') return;
        localStorage.setItem(KEYS.BUSINESS_TOKEN, token);
    },

    getUser: (): UserData | null => {
        if (typeof window === 'undefined') return null;
        const data = localStorage.getItem(KEYS.USER_DATA);
        return data ? JSON.parse(data) : null;
    },

    setUser: (user: UserData) => {
        if (typeof window === 'undefined') return;
        localStorage.setItem(KEYS.USER_DATA, JSON.stringify(user));
    },

    getBusinessesData: (): BusinessData[] | null => {
        if (typeof window === 'undefined') return null;
        const data = localStorage.getItem(KEYS.BUSINESSES_DATA);
        return data ? JSON.parse(data) : null;
    },

    setBusinessesData: (businesses: BusinessData[]) => {
        if (typeof window === 'undefined') return;
        localStorage.setItem(KEYS.BUSINESSES_DATA, JSON.stringify(businesses));
    },

    setActiveBusiness: (id: number) => {
        if (typeof window === 'undefined') return;
        localStorage.setItem(KEYS.ACTIVE_BUSINESS_ID, id.toString());
    },

    setBusinessColors: (colors: BusinessColors) => {
        if (typeof window === 'undefined') return;
        localStorage.setItem(KEYS.BUSINESS_COLORS, JSON.stringify(colors));
    },

    getBusinessColors: (): BusinessColors | null => {
        if (typeof window === 'undefined') return null;
        const data = localStorage.getItem(KEYS.BUSINESS_COLORS);
        return data ? JSON.parse(data) : null;
    },

    removeUserPermissions: () => {
        if (typeof window === 'undefined') return;
        localStorage.removeItem(KEYS.PERMISSIONS);
    },

    clearSession: () => {
        if (typeof window === 'undefined') return;
        localStorage.removeItem(KEYS.SESSION_TOKEN);
        localStorage.removeItem(KEYS.BUSINESS_TOKEN);
        localStorage.removeItem(KEYS.USER_DATA);
        localStorage.removeItem(KEYS.BUSINESSES_DATA);
        localStorage.removeItem(KEYS.ACTIVE_BUSINESS_ID);
        localStorage.removeItem(KEYS.PERMISSIONS);
        localStorage.removeItem(KEYS.BUSINESS_COLORS);
    }
};

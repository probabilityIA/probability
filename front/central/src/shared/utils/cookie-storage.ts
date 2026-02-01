/**
 * Cookie Storage para manejar autenticación en contextos de iframe (Shopify)
 *
 * Usa cookies con SameSite=None; Secure para que funcionen en iframes de terceros
 * Fallback a localStorage cuando no es iframe
 */

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

export interface ResourcePermission {
    resource: string;
    actions: string[];
    active: boolean;
}

export interface UserPermissions {
    is_super: boolean;
    business_id: number;
    business_name: string;
    role_id: number;
    role_name: string;
    resources: ResourcePermission[];
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

/**
 * Detecta si estamos en un iframe
 */
function isInIframe(): boolean {
    if (typeof window === 'undefined') return false;
    try {
        return window.self !== window.top;
    } catch (e) {
        return true;
    }
}

/**
 * Detecta si estamos en un iframe de Shopify
 */
function isShopifyIframe(): boolean {
    if (typeof window === 'undefined') return false;
    try {
        // Verificar si el referrer o parent es de Shopify
        const referrer = document.referrer.toLowerCase();
        return (
            isInIframe() &&
            (referrer.includes('shopify.com') ||
             referrer.includes('myshopify.com'))
        );
    } catch (e) {
        return false;
    }
}

/**
 * Set cookie con SameSite=None para iframes de terceros
 */
function setCookie(name: string, value: string, days: number = 7): void {
    if (typeof window === 'undefined') return;

    const expires = new Date();
    expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000);

    // Si estamos en iframe, SIEMPRE usar SameSite=None; Secure
    const sameSite = isInIframe() ? 'None' : 'Lax';
    const secure = isInIframe() ? '; Secure' : '';

    document.cookie = `${name}=${value}; expires=${expires.toUTCString()}; path=/; SameSite=${sameSite}${secure}`;
}

/**
 * Get cookie
 */
function getCookie(name: string): string | null {
    if (typeof window === 'undefined') return null;

    const nameEQ = name + '=';
    const ca = document.cookie.split(';');

    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) === ' ') c = c.substring(1, c.length);
        if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
    }

    return null;
}

/**
 * Delete cookie
 */
function deleteCookie(name: string): void {
    if (typeof window === 'undefined') return;

    const sameSite = isInIframe() ? 'None' : 'Lax';
    const secure = isInIframe() ? '; Secure' : '';

    document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/; SameSite=${sameSite}${secure}`;
}

/**
 * Storage wrapper que usa cookies en iframe y localStorage fuera de iframe
 */
function setItem(key: string, value: string): void {
    if (typeof window === 'undefined') return;

    try {
        if (isInIframe()) {
            // En iframe, usar cookies
            setCookie(key, value);
        } else {
            // Fuera de iframe, usar localStorage
            localStorage.setItem(key, value);
        }
    } catch (e) {
        console.error(`Error setting ${key}:`, e);
        // Fallback a cookies si localStorage falla
        setCookie(key, value);
    }
}

/**
 * Get item del storage apropiado
 */
function getItem(key: string): string | null {
    if (typeof window === 'undefined') return null;

    try {
        if (isInIframe()) {
            // En iframe, leer de cookies
            return getCookie(key);
        } else {
            // Fuera de iframe, leer de localStorage
            return localStorage.getItem(key);
        }
    } catch (e) {
        console.error(`Error getting ${key}:`, e);
        // Fallback a cookies si localStorage falla
        return getCookie(key);
    }
}

/**
 * Remove item del storage apropiado
 */
function removeItem(key: string): void {
    if (typeof window === 'undefined') return;

    try {
        if (isInIframe()) {
            deleteCookie(key);
        } else {
            localStorage.removeItem(key);
        }
    } catch (e) {
        console.error(`Error removing ${key}:`, e);
        deleteCookie(key);
    }
}

export const CookieStorage = {
    // Utils
    isInIframe,
    isShopifyIframe,

    // Session Token
    getSessionToken: (): string | null => {
        return getItem(KEYS.SESSION_TOKEN);
    },

    setSessionToken: (token: string) => {
        setItem(KEYS.SESSION_TOKEN, token);
    },

    // Business Token
    getBusinessToken: (): string | null => {
        return getItem(KEYS.BUSINESS_TOKEN);
    },

    setBusinessToken: (token: string) => {
        setItem(KEYS.BUSINESS_TOKEN, token);
    },

    // User Data
    getUser: (): UserData | null => {
        const data = getItem(KEYS.USER_DATA);
        return data ? JSON.parse(data) : null;
    },

    setUser: (user: UserData) => {
        setItem(KEYS.USER_DATA, JSON.stringify(user));
    },

    // Businesses Data
    getBusinessesData: (): BusinessData[] | null => {
        const data = getItem(KEYS.BUSINESSES_DATA);
        return data ? JSON.parse(data) : null;
    },

    setBusinessesData: (businesses: BusinessData[]) => {
        setItem(KEYS.BUSINESSES_DATA, JSON.stringify(businesses));
    },

    // Active Business
    setActiveBusiness: (id: number) => {
        setItem(KEYS.ACTIVE_BUSINESS_ID, id.toString());
    },

    // Business Colors
    setBusinessColors: (colors: BusinessColors) => {
        setItem(KEYS.BUSINESS_COLORS, JSON.stringify(colors));
    },

    getBusinessColors: (): BusinessColors | null => {
        const data = getItem(KEYS.BUSINESS_COLORS);
        return data ? JSON.parse(data) : null;
    },

    // Permissions
    getPermissions: (): UserPermissions | null => {
        const data = getItem(KEYS.PERMISSIONS);
        return data ? JSON.parse(data) : null;
    },

    setPermissions: (permissions: UserPermissions) => {
        setItem(KEYS.PERMISSIONS, JSON.stringify(permissions));
    },

    removeUserPermissions: () => {
        removeItem(KEYS.PERMISSIONS);
    },

    // Clear Session
    clearSession: () => {
        removeItem(KEYS.SESSION_TOKEN);
        removeItem(KEYS.BUSINESS_TOKEN);
        removeItem(KEYS.USER_DATA);
        removeItem(KEYS.BUSINESSES_DATA);
        removeItem(KEYS.ACTIVE_BUSINESS_ID);
        removeItem(KEYS.PERMISSIONS);
        removeItem(KEYS.BUSINESS_COLORS);
    }
};

/**
 * CookieStorage Simplificado - SIEMPRE usa cookies
 *
 * Ventajas:
 * - Un solo code path (más simple)
 * - Funciona en iframes y páginas normales
 * - Comportamiento consistente
 *
 * Desventajas:
 * - Requiere HTTPS en desarrollo (cookies Secure)
 * - Safari/Firefox pueden bloquear cookies third-party
 */

export interface UserData {
    userId: string;
    name: string;
    email: string;
    role: string;
    avatarUrl?: string;
    is_super_admin?: boolean;
    scope?: string;
}

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
 * Set cookie con SameSite apropiado
 */
function setCookie(name: string, value: string, days: number = 7): void {
    if (typeof window === 'undefined') return;

    const expires = new Date();
    expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000);

    // SIEMPRE usar SameSite=None para máxima compatibilidad
    // (funciona en iframes Y en páginas normales)
    document.cookie = `${name}=${value}; expires=${expires.toUTCString()}; path=/; SameSite=None; Secure`;
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
    document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/; SameSite=None; Secure`;
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

export const SimpleCookieStorage = {
    // Utils
    isInIframe,

    // Session Token
    getSessionToken: (): string | null => getCookie(KEYS.SESSION_TOKEN),
    setSessionToken: (token: string) => setCookie(KEYS.SESSION_TOKEN, token),

    // Business Token
    getBusinessToken: (): string | null => getCookie(KEYS.BUSINESS_TOKEN),
    setBusinessToken: (token: string) => setCookie(KEYS.BUSINESS_TOKEN, token),

    // User Data
    getUser: (): UserData | null => {
        const data = getCookie(KEYS.USER_DATA);
        return data ? JSON.parse(decodeURIComponent(data)) : null;
    },
    setUser: (user: UserData) => {
        setCookie(KEYS.USER_DATA, encodeURIComponent(JSON.stringify(user)));
    },

    // Clear Session
    clearSession: () => {
        deleteCookie(KEYS.SESSION_TOKEN);
        deleteCookie(KEYS.BUSINESS_TOKEN);
        deleteCookie(KEYS.USER_DATA);
        deleteCookie(KEYS.BUSINESSES_DATA);
        deleteCookie(KEYS.ACTIVE_BUSINESS_ID);
        deleteCookie(KEYS.PERMISSIONS);
        deleteCookie(KEYS.BUSINESS_COLORS);
    }
};

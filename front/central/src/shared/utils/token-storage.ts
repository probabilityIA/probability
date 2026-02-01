/**
 * Token Storage - Wrapper que usa CookieStorage internamente
 *
 * DEPRECATED: Este archivo existe solo para compatibilidad.
 * Usa CookieStorage directamente en código nuevo.
 *
 * CookieStorage detecta automáticamente si estamos en iframe y usa:
 * - Cookies con SameSite=None; Secure en iframes (Shopify)
 * - localStorage en páginas normales
 */

import { CookieStorage } from './cookie-storage';

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

/**
 * @deprecated Use CookieStorage instead
 *
 * TokenStorage ahora usa CookieStorage internamente para soportar iframes de Shopify
 */
export const TokenStorage = {
    getSessionToken: () => CookieStorage.getSessionToken(),
    setSessionToken: (token: string) => CookieStorage.setSessionToken(token),
    getBusinessToken: () => CookieStorage.getBusinessToken(),
    setBusinessToken: (token: string) => CookieStorage.setBusinessToken(token),
    getUser: () => CookieStorage.getUser(),
    setUser: (user: UserData) => CookieStorage.setUser(user),
    getBusinessesData: () => CookieStorage.getBusinessesData(),
    setBusinessesData: (businesses: BusinessData[]) => CookieStorage.setBusinessesData(businesses),
    setActiveBusiness: (id: number) => CookieStorage.setActiveBusiness(id),
    setBusinessColors: (colors: BusinessColors) => CookieStorage.setBusinessColors(colors),
    getBusinessColors: () => CookieStorage.getBusinessColors(),
    getPermissions: () => CookieStorage.getPermissions(),
    setPermissions: (permissions: UserPermissions) => CookieStorage.setPermissions(permissions),
    removeUserPermissions: () => CookieStorage.removeUserPermissions(),
    clearSession: () => CookieStorage.clearSession(),
};

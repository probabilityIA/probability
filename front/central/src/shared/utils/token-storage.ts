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

import {
    CookieStorage,
    type BusinessColors,
    type BusinessData,
    type UserData,
    type UserPermissions
} from './cookie-storage';

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

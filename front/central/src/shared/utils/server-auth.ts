/**
 * Helper centralizado para autenticación en Server Actions
 *
 * Automáticamente obtiene el token de:
 * 1. Parámetro explícito (para iframes donde cookies están bloqueadas)
 * 2. Cookies (para navegador normal)
 *
 * Uso en Server Actions:
 * ```typescript
 * import { getAuthToken } from '@/shared/utils/server-auth';
 *
 * export async function myAction(params: any, token?: string | null) {
 *     const authToken = await getAuthToken(token);
 *     const repo = new MyRepository(authToken);
 *     // ...
 * }
 * ```
 */

import { cookies } from 'next/headers';

/**
 * Obtiene el token de autenticación de múltiples fuentes
 *
 * @param explicitToken Token explícito pasado como parámetro (para iframes)
 * @returns Token de autenticación o null
 */
export async function getAuthToken(explicitToken?: string | null): Promise<string | null> {
    // 1. Si hay token explícito, usarlo (iframe de Shopify)
    if (explicitToken) {
        console.log('[Auth] Usando token explícito (iframe)');
        return explicitToken;
    }

    // 2. Intentar leer de cookies (navegador normal)
    try {
        const cookieStore = await cookies();
        const token = cookieStore.get('session_token')?.value || null;

        if (token) {
            console.log('[Auth] Usando token de cookie (navegador normal)');
        } else {
            console.log('[Auth] ⚠️ No se encontró token en cookies ni parámetro');
        }

        return token;
    } catch (error) {
        console.error('[Auth] Error leyendo cookies:', error);
        return null;
    }
}

/**
 * Middleware helper para validar que existe un token
 * Lanza error si no hay token disponible
 */
export async function requireAuthToken(explicitToken?: string | null): Promise<string> {
    const token = await getAuthToken(explicitToken);

    if (!token) {
        throw new Error('Token de autorización requerido');
    }

    return token;
}

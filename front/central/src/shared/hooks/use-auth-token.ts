/**
 * Hook para obtener el token de autenticación en Client Components
 *
 * Uso en componentes que llaman Server Actions:
 * ```typescript
 * const token = useAuthToken();
 * const data = await myAction(params, token);
 * ```
 */

'use client';

import { useMemo } from 'react';
import { TokenStorage } from '../utils/token-storage';

/**
 * Detecta si estamos en iframe de Shopify
 */
function isShopifyIframe(): boolean {
    if (typeof window === 'undefined') return false;
    try {
        const referrer = document.referrer.toLowerCase();
        const isIframe = window.self !== window.top;
        return (
            isIframe &&
            (referrer.includes('shopify.com') ||
             referrer.includes('myshopify.com'))
        );
    } catch (e) {
        return false;
    }
}

/**
 * Hook que retorna el token si estamos en iframe de Shopify
 * En navegador normal, retorna null (se usarán cookies automáticamente)
 */
export function useAuthToken(): string | null {
    return useMemo(() => {
        // Solo obtener token si estamos en iframe de Shopify
        if (isShopifyIframe()) {
            const token = TokenStorage.getSessionToken();
            if (token) {
                console.log('🛍️ useAuthToken: Token obtenido para iframe de Shopify');
            }
            return token;
        }

        // En navegador normal, las cookies HttpOnly se manejan automáticamente
        return null;
    }, []);
}

/**
 * Hook que detecta si estamos en iframe de Shopify
 */
export function useIsShopifyIframe(): boolean {
    return useMemo(() => isShopifyIframe(), []);
}

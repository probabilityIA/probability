/**
 * Hook universal para API calls que funciona en cualquier contexto
 *
 * Automáticamente detecta:
 * - Iframe de Shopify: Usa fetch directo con token de sessionStorage
 * - Navegador normal: Usa Server Actions con cookies HttpOnly
 *
 * Uso:
 * ```typescript
 * const api = useUniversalApi();
 * const data = await api.get('/orders');
 * ```
 */

'use client';

import { useMemo } from 'react';
import { apiClient } from '../utils/api-client';

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
 * Hook que retorna el cliente API correcto según el contexto
 */
export function useUniversalApi() {
    const isIframe = useMemo(() => isShopifyIframe(), []);

    // En iframe de Shopify, siempre usar cliente directo
    // En navegador normal, también podemos usar cliente directo ya que las cookies se envían automáticamente
    return useMemo(() => ({
        client: apiClient,
        isShopifyIframe: isIframe,
        get: apiClient.get.bind(apiClient),
        post: apiClient.post.bind(apiClient),
        put: apiClient.put.bind(apiClient),
        delete: apiClient.delete.bind(apiClient),
        patch: apiClient.patch.bind(apiClient),
    }), [isIframe]);
}

/**
 * Hook para detectar si estamos en iframe de Shopify
 */
export function useIsShopifyIframe() {
    return useMemo(() => isShopifyIframe(), []);
}

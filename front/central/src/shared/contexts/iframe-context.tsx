/**
 * Context para detectar si estamos en iframe de Shopify
 * y proporcionar el cliente API correcto automáticamente
 */

'use client';

import { createContext, useContext, ReactNode, useMemo } from 'react';
import { apiClient } from '../utils/api-client';

interface IframeContextType {
    isShopifyIframe: boolean;
    apiClient: typeof apiClient;
    shouldUseClientFetch: boolean;
}

const IframeContext = createContext<IframeContextType>({
    isShopifyIframe: false,
    apiClient,
    shouldUseClientFetch: false,
});

/**
 * Detecta si estamos en iframe de Shopify
 */
function detectShopifyIframe(): boolean {
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

export function IframeProvider({ children }: { children: ReactNode }) {
    const contextValue = useMemo(() => {
        const isShopifyIframe = detectShopifyIframe();

        return {
            isShopifyIframe,
            apiClient,
            // En iframe de Shopify, SIEMPRE usar client fetch
            shouldUseClientFetch: isShopifyIframe,
        };
    }, []);

    return (
        <IframeContext.Provider value={contextValue}>
            {children}
        </IframeContext.Provider>
    );
}

/**
 * Hook para acceder al contexto de iframe
 */
export function useIframeContext() {
    return useContext(IframeContext);
}

/**
 * Hook simplificado para obtener el cliente API
 */
export function useApiClient() {
    const { apiClient } = useIframeContext();
    return apiClient;
}

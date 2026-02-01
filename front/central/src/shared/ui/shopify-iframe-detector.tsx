'use client';

import { useEffect, useState } from 'react';
import { CookieStorage } from '@/shared/utils';

interface ShopifyIframeDetectorProps {
    children: React.ReactNode;
    onShopifyDetected?: (shop: string) => void;
}

/**
 * Detecta si estamos en un iframe de Shopify y configura el ambiente apropiadamente
 */
export function ShopifyIframeDetector({ children, onShopifyDetected }: ShopifyIframeDetectorProps) {
    const [isReady, setIsReady] = useState(false);
    const [isShopifyIframe, setIsShopifyIframe] = useState(false);

    useEffect(() => {
        // Detectar si estamos en iframe de Shopify
        const inShopifyIframe = CookieStorage.isShopifyIframe();
        setIsShopifyIframe(inShopifyIframe);

        if (inShopifyIframe) {
            // Extraer shop parameter de la URL
            const urlParams = new URLSearchParams(window.location.search);
            const shop = urlParams.get('shop');

            if (shop) {
                console.log('🛍️ Shopify iframe detected:', shop);
                onShopifyDetected?.(shop);

                // Guardar shop en sessionStorage para acceso posterior
                sessionStorage.setItem('shopify_shop', shop);
            }
        }

        setIsReady(true);
    }, [onShopifyDetected]);

    // Mostrar loading mientras detectamos el contexto
    if (!isReady) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-50">
                <div className="text-center">
                    <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600"></div>
                    <p className="mt-4 text-gray-600">Inicializando...</p>
                </div>
            </div>
        );
    }

    // Debug info en desarrollo
    if (process.env.NODE_ENV === 'development' && isShopifyIframe) {
        console.log('🔍 Storage mode:', CookieStorage.isInIframe() ? 'Cookies (iframe)' : 'LocalStorage');
    }

    return <>{children}</>;
}

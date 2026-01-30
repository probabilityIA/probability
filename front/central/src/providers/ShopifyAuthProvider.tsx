'use client';

import { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { useSearchParams } from 'next/navigation';

interface ShopifyAuthContextType {
    isShopifyEmbedded: boolean;
    sessionToken: string | null;
    isLoading: boolean;
    shopOrigin: string | null;
}

const ShopifyAuthContext = createContext<ShopifyAuthContextType>({
    isShopifyEmbedded: false,
    sessionToken: null,
    isLoading: true,
    shopOrigin: null,
});

export const useShopifyAuth = () => useContext(ShopifyAuthContext);

export function ShopifyAuthProvider({ children }: { children: ReactNode }) {
    const [sessionToken, setSessionToken] = useState<string | null>(null);
    const [isShopifyEmbedded, setIsShopifyEmbedded] = useState(false);
    const [isLoading, setIsLoading] = useState(true);
    const searchParams = useSearchParams();
    const shopOrigin = searchParams.get('shop');
    const host = searchParams.get('host');

    useEffect(() => {
        // Detectar si estamos en un iframe de Shopify revisando los params
        if (shopOrigin || host) {
            setIsShopifyEmbedded(true);
            initializeShopifyAuth();
        } else {
            setIsLoading(false);
        }
    }, [shopOrigin, host]);

    const initializeShopifyAuth = async () => {
        try {
            // Esperar a que window.shopify esté disponible (inyectado por el CDN)
            if (typeof window !== 'undefined' && 'shopify' in window) {
                const appBridge = (window as any).shopify;
                if (appBridge && appBridge.idt) {
                    // Obtener el token inicial
                    const token = await appBridge.idt.getToken();
                    setSessionToken(token);

                    // Suscribirse a cambios de token (se renuevan cada minuto)
                    // Nota: La versión actual de App Bridge maneja esto automáticamente en su mayoría,
                    // pero podemos interceptarlo si es necesario. Por ahora un fetch simple basta.
                }
            }
        } catch (error) {
            console.error("Error inicializando Shopify Auth:", error);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <ShopifyAuthContext.Provider value={{ isShopifyEmbedded, sessionToken, isLoading, shopOrigin }}>
            {children}
        </ShopifyAuthContext.Provider>
    );
}

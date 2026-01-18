'use client';

import { useEffect } from 'react';

interface ShopifyAppBridgeProviderProps {
    children: React.ReactNode;
}

export function ShopifyAppBridgeProvider({ children }: ShopifyAppBridgeProviderProps) {
    useEffect(() => {
        if (typeof window !== 'undefined') {
            const urlParams = new URLSearchParams(window.location.search);
            const apiKey = process.env.NEXT_PUBLIC_SHOPIFY_API_KEY;

            // App Bridge v4 automatically initializes if the script is present and params are there.
            // We can inject the script if we are in an embedded contexts.

            if (urlParams.get('host') && apiKey && !document.getElementById('shopify-app-bridge')) {
                const script = document.createElement('script');
                script.id = 'shopify-app-bridge';
                script.src = `https://cdn.shopify.com/shopify-cloud/app-bridge.js?apiKey=${apiKey}`;
                script.setAttribute('data-api-key', apiKey);
                document.head.appendChild(script);
            }
        }
    }, []);

    return <>{children}</>;
}

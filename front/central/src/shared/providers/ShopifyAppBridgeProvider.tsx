'use client';

interface ShopifyAppBridgeProviderProps {
    children: React.ReactNode;
}

export function ShopifyAppBridgeProvider({ children }: ShopifyAppBridgeProviderProps) {
    // App Bridge v4 script is now loaded in the root layout via next/script.
    // The CDN script automatically initializes when loaded in an embedded context.
    return <>{children}</>;
}

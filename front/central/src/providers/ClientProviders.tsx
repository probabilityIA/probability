'use client';

import { ShopifyAuthProvider } from '@/providers/ShopifyAuthProvider';
import { ReactNode, Suspense } from 'react';

export function ClientProviders({ children }: { children: ReactNode }) {
    return (
        <Suspense fallback={null}>
            <ShopifyAuthProvider>
                {children}
            </ShopifyAuthProvider>
        </Suspense>
    );
}

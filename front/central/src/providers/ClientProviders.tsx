'use client';

import { ShopifyAuthProvider } from '@/providers/ShopifyAuthProvider';
import { ReactNode } from 'react';

export function ClientProviders({ children }: { children: ReactNode }) {
    return (
        <ShopifyAuthProvider>
            {children}
        </ShopifyAuthProvider>
    );
}

'use client';

import { ShopifyAuthProvider } from '@/providers/ShopifyAuthProvider';
import { ThemeProvider } from '@/shared/providers/theme-provider';
import { ReactNode, Suspense } from 'react';

export function ClientProviders({ children }: { children: ReactNode }) {
    return (
        <Suspense fallback={null}>
            <ShopifyAuthProvider>
                <ThemeProvider>
                    {children}
                </ThemeProvider>
            </ShopifyAuthProvider>
        </Suspense>
    );
}

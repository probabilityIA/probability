'use client';

import { ShopifyAuthProvider } from '@/providers/ShopifyAuthProvider';
import { ThemeProvider } from '@/shared/providers/theme-provider';
import { DarkModeProvider } from '@/shared/contexts/dark-mode-context';
import { ReactNode, Suspense } from 'react';

export function ClientProviders({ children }: { children: ReactNode }) {
    return (
        <Suspense fallback={null}>
            <ShopifyAuthProvider>
                <DarkModeProvider>
                    <ThemeProvider>
                        {children}
                    </ThemeProvider>
                </DarkModeProvider>
            </ShopifyAuthProvider>
        </Suspense>
    );
}

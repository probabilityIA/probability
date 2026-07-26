'use client';

import type { ReactNode } from 'react';
import { useSelectedBusiness } from './selected-business-context';

export function WalletBusinessProvider({ children }: { children: ReactNode }) {
    return <>{children}</>;
}

export function useWalletBusiness() {
    return useSelectedBusiness();
}

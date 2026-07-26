'use client';

import type { ReactNode } from 'react';
import { useSelectedBusiness } from './selected-business-context';

export function StorefrontBusinessProvider({ children }: { children: ReactNode }) {
    return <>{children}</>;
}

export function useStorefrontBusiness() {
    return useSelectedBusiness();
}

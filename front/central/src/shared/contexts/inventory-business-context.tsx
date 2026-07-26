'use client';

import type { ReactNode } from 'react';
import { useSelectedBusiness } from './selected-business-context';

export function InventoryBusinessProvider({ children }: { children: ReactNode }) {
    return <>{children}</>;
}

export function useInventoryBusiness() {
    return useSelectedBusiness();
}

'use client';

import type { ReactNode } from 'react';
import { useSelectedBusiness } from './selected-business-context';

export function OrdersBusinessProvider({ children }: { children: ReactNode }) {
    return <>{children}</>;
}

export function useOrdersBusiness() {
    return useSelectedBusiness();
}

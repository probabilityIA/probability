'use client';

import type { ReactNode } from 'react';
import { useSelectedBusiness } from './selected-business-context';

export function DeliveryBusinessProvider({ children }: { children: ReactNode }) {
    return <>{children}</>;
}

export function useDeliveryBusiness() {
    return useSelectedBusiness();
}

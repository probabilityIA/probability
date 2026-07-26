'use client';

import type { ReactNode } from 'react';
import { useSelectedBusiness } from './selected-business-context';

export function InvoicingBusinessProvider({ children }: { children: ReactNode }) {
    return <>{children}</>;
}

export function useInvoicingBusiness() {
    return useSelectedBusiness();
}

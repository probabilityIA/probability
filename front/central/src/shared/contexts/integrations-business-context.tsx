'use client';

import type { ReactNode } from 'react';
import { useSelectedBusiness } from './selected-business-context';

export function IntegrationsBusinessProvider({ children }: { children: ReactNode }) {
    return <>{children}</>;
}

export function useIntegrationsBusiness() {
    return useSelectedBusiness();
}

'use client';

import type { ReactNode } from 'react';
import { useSelectedBusiness } from './selected-business-context';

export function NotificationBusinessProvider({ children }: { children: ReactNode }) {
    return <>{children}</>;
}

export function useNotificationBusiness() {
    return useSelectedBusiness();
}

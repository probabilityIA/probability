'use client';

import { createContext, useContext, useState, type ReactNode } from 'react';

interface NotificationBusinessContextType {
    selectedBusinessId: number | null;
    setSelectedBusinessId: (id: number | null) => void;
}

const NotificationBusinessContext = createContext<NotificationBusinessContextType | null>(null);

export function NotificationBusinessProvider({ children }: { children: ReactNode }) {
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);

    return (
        <NotificationBusinessContext.Provider value={{ selectedBusinessId, setSelectedBusinessId }}>
            {children}
        </NotificationBusinessContext.Provider>
    );
}

export function useNotificationBusiness() {
    const ctx = useContext(NotificationBusinessContext);
    if (!ctx) {
        throw new Error('useNotificationBusiness must be used within NotificationBusinessProvider');
    }
    return ctx;
}

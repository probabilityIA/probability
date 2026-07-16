'use client';

import { createContext, useContext, useState, type ReactNode } from 'react';

interface IntegrationsBusinessContextType {
    selectedBusinessId: number | null;
    setSelectedBusinessId: (id: number | null) => void;
}

const IntegrationsBusinessContext = createContext<IntegrationsBusinessContextType | null>(null);

export function IntegrationsBusinessProvider({ children }: { children: ReactNode }) {
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);

    return (
        <IntegrationsBusinessContext.Provider value={{ selectedBusinessId, setSelectedBusinessId }}>
            {children}
        </IntegrationsBusinessContext.Provider>
    );
}

export function useIntegrationsBusiness() {
    const ctx = useContext(IntegrationsBusinessContext);
    if (!ctx) {
        throw new Error('useIntegrationsBusiness must be used within IntegrationsBusinessProvider');
    }
    return ctx;
}

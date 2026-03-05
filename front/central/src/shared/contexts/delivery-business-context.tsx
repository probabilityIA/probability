'use client';

import { createContext, useContext, useState, type ReactNode } from 'react';

interface DeliveryBusinessContextType {
    selectedBusinessId: number | null;
    setSelectedBusinessId: (id: number | null) => void;
}

const DeliveryBusinessContext = createContext<DeliveryBusinessContextType | null>(null);

export function DeliveryBusinessProvider({ children }: { children: ReactNode }) {
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);

    return (
        <DeliveryBusinessContext.Provider value={{ selectedBusinessId, setSelectedBusinessId }}>
            {children}
        </DeliveryBusinessContext.Provider>
    );
}

export function useDeliveryBusiness() {
    const ctx = useContext(DeliveryBusinessContext);
    if (!ctx) {
        throw new Error('useDeliveryBusiness must be used within DeliveryBusinessProvider');
    }
    return ctx;
}

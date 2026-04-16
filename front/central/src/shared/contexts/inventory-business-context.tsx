'use client';

import { createContext, useContext, useState, type ReactNode } from 'react';

interface InventoryBusinessContextType {
    selectedBusinessId: number | null;
    setSelectedBusinessId: (id: number | null) => void;
}

const InventoryBusinessContext = createContext<InventoryBusinessContextType | null>(null);

export function InventoryBusinessProvider({ children }: { children: ReactNode }) {
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);

    return (
        <InventoryBusinessContext.Provider value={{ selectedBusinessId, setSelectedBusinessId }}>
            {children}
        </InventoryBusinessContext.Provider>
    );
}

export function useInventoryBusiness() {
    const ctx = useContext(InventoryBusinessContext);
    if (!ctx) {
        throw new Error('useInventoryBusiness must be used within InventoryBusinessProvider');
    }
    return ctx;
}

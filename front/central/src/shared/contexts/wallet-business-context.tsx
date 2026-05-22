'use client';

import { createContext, useContext, useState, type ReactNode } from 'react';

interface WalletBusinessContextType {
    selectedBusinessId: number | null;
    setSelectedBusinessId: (id: number | null) => void;
}

const WalletBusinessContext = createContext<WalletBusinessContextType | null>(null);

export function WalletBusinessProvider({ children }: { children: ReactNode }) {
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);

    return (
        <WalletBusinessContext.Provider value={{ selectedBusinessId, setSelectedBusinessId }}>
            {children}
        </WalletBusinessContext.Provider>
    );
}

export function useWalletBusiness() {
    const ctx = useContext(WalletBusinessContext);
    if (!ctx) {
        throw new Error('useWalletBusiness must be used within WalletBusinessProvider');
    }
    return ctx;
}

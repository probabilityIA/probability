'use client';

import { createContext, useContext, useState, type ReactNode } from 'react';

interface InvoicingBusinessContextValue {
    selectedBusinessId: number | null;
    setSelectedBusinessId: (id: number | null) => void;
}

const InvoicingBusinessContext = createContext<InvoicingBusinessContextValue | undefined>(undefined);

export function InvoicingBusinessProvider({ children }: { children: ReactNode }) {
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    return (
        <InvoicingBusinessContext.Provider value={{ selectedBusinessId, setSelectedBusinessId }}>
            {children}
        </InvoicingBusinessContext.Provider>
    );
}

export function useInvoicingBusiness() {
    const ctx = useContext(InvoicingBusinessContext);
    if (!ctx) throw new Error('useInvoicingBusiness must be used within InvoicingBusinessProvider');
    return ctx;
}

'use client';

import { createContext, useContext, useState, type ReactNode } from 'react';

interface OrdersBusinessContextType {
    selectedBusinessId: number | null;
    setSelectedBusinessId: (id: number | null) => void;
}

const OrdersBusinessContext = createContext<OrdersBusinessContextType | null>(null);

export function OrdersBusinessProvider({ children }: { children: ReactNode }) {
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);

    return (
        <OrdersBusinessContext.Provider value={{ selectedBusinessId, setSelectedBusinessId }}>
            {children}
        </OrdersBusinessContext.Provider>
    );
}

export function useOrdersBusiness() {
    const ctx = useContext(OrdersBusinessContext);
    if (!ctx) {
        throw new Error('useOrdersBusiness must be used within OrdersBusinessProvider');
    }
    return ctx;
}

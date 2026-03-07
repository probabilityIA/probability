'use client';

import { createContext, useContext, useState, type ReactNode } from 'react';

interface StorefrontBusinessContextType {
    selectedBusinessId: number | null;
    setSelectedBusinessId: (id: number | null) => void;
}

const StorefrontBusinessContext = createContext<StorefrontBusinessContextType | null>(null);

export function StorefrontBusinessProvider({ children }: { children: ReactNode }) {
    const [selectedBusinessId, setSelectedBusinessIdState] = useState<number | null>(() => {
        if (typeof document !== 'undefined') {
            const match = document.cookie.match(/storefront_business_id=(\d+)/);
            return match ? parseInt(match[1], 10) : null;
        }
        return null;
    });

    const setSelectedBusinessId = (id: number | null) => {
        setSelectedBusinessIdState(id);
        if (id) {
            document.cookie = `storefront_business_id=${id}; path=/; max-age=86400`;
        } else {
            document.cookie = 'storefront_business_id=; path=/; max-age=0';
        }
    };

    return (
        <StorefrontBusinessContext.Provider value={{ selectedBusinessId, setSelectedBusinessId }}>
            {children}
        </StorefrontBusinessContext.Provider>
    );
}

export function useStorefrontBusiness() {
    const ctx = useContext(StorefrontBusinessContext);
    if (!ctx) {
        throw new Error('useStorefrontBusiness must be used within StorefrontBusinessProvider');
    }
    return ctx;
}

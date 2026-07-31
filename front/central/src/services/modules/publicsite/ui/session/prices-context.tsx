'use client';

import { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { getTiendaPricesAction } from '../../infra/actions';

const PricesContext = createContext<Record<string, number>>({});

export function TiendaPricesProvider({ slug, children }: { slug: string; children: ReactNode }) {
    const [prices, setPrices] = useState<Record<string, number>>({});

    useEffect(() => {
        let cancelled = false;
        getTiendaPricesAction(slug).then(p => {
            if (!cancelled && p) setPrices(p);
        });
        return () => { cancelled = true; };
    }, [slug]);

    return <PricesContext.Provider value={prices}>{children}</PricesContext.Provider>;
}

export function useCustomerPrice(productId: string, basePrice: number): { price: number; isCustom: boolean } {
    const prices = useContext(PricesContext);
    const custom = prices[productId];
    if (custom !== undefined && custom !== basePrice) {
        return { price: custom, isCustom: true };
    }
    return { price: basePrice, isCustom: false };
}

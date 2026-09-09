'use client';

import { createContext, useCallback, useContext, useEffect, useMemo, useState, ReactNode } from 'react';
import { StorefrontProduct } from '../../domain/types';
import { useStorefrontBusiness } from '@/shared/contexts/storefront-business-context';

export interface CartItem {
    product: StorefrontProduct;
    quantity: number;
}

interface CartContextValue {
    items: CartItem[];
    isOpen: boolean;
    totalItems: number;
    totalPrice: number;
    currency: string;
    open: () => void;
    close: () => void;
    toggle: () => void;
    addItem: (product: StorefrontProduct, quantity?: number) => void;
    incrementItem: (productId: string) => void;
    decrementItem: (productId: string) => void;
    removeItem: (productId: string) => void;
    clear: () => void;
    maxQuantity: (product: StorefrontProduct) => number;
}

const CartContext = createContext<CartContextValue | undefined>(undefined);

function maxQuantityFor(product: StorefrontProduct): number {
    return product.track_inventory ? Math.max(0, product.stock_quantity) : Infinity;
}

export function CartProvider({ children }: { children: ReactNode }) {
    const { selectedBusinessId } = useStorefrontBusiness();
    const storageKey = `storefront_cart_${selectedBusinessId ?? 'default'}`;
    const [items, setItems] = useState<CartItem[]>([]);
    const [isOpen, setIsOpen] = useState(false);
    const [hydrated, setHydrated] = useState(false);

    useEffect(() => {
        setHydrated(false);
        try {
            const raw = localStorage.getItem(storageKey);
            setItems(raw ? JSON.parse(raw) : []);
        } catch {
            setItems([]);
        }
        setHydrated(true);
    }, [storageKey]);

    useEffect(() => {
        if (!hydrated) return;
        try {
            localStorage.setItem(storageKey, JSON.stringify(items));
        } catch {}
    }, [items, storageKey, hydrated]);

    const addItem = useCallback((product: StorefrontProduct, quantity = 1) => {
        setItems(prev => {
            const cap = maxQuantityFor(product);
            const existing = prev.find(item => item.product.id === product.id);
            if (existing) {
                return prev.map(item =>
                    item.product.id === product.id
                        ? { ...item, quantity: Math.min(item.quantity + quantity, cap) }
                        : item
                );
            }
            if (cap <= 0) return prev;
            return [...prev, { product, quantity: Math.min(quantity, cap) }];
        });
    }, []);

    const incrementItem = useCallback((productId: string) => {
        setItems(prev =>
            prev.map(item =>
                item.product.id === productId
                    ? { ...item, quantity: Math.min(item.quantity + 1, maxQuantityFor(item.product)) }
                    : item
            )
        );
    }, []);

    const decrementItem = useCallback((productId: string) => {
        setItems(prev =>
            prev
                .map(item =>
                    item.product.id === productId
                        ? { ...item, quantity: item.quantity - 1 }
                        : item
                )
                .filter(item => item.quantity > 0)
        );
    }, []);

    const removeItem = useCallback((productId: string) => {
        setItems(prev => prev.filter(item => item.product.id !== productId));
    }, []);

    const clear = useCallback(() => setItems([]), []);

    const totalItems = useMemo(() => items.reduce((sum, item) => sum + item.quantity, 0), [items]);
    const totalPrice = useMemo(() => items.reduce((sum, item) => sum + item.product.price * item.quantity, 0), [items]);
    const currency = items[0]?.product.currency || 'COP';

    const value: CartContextValue = {
        items,
        isOpen,
        totalItems,
        totalPrice,
        currency,
        open: () => setIsOpen(true),
        close: () => setIsOpen(false),
        toggle: () => setIsOpen(prev => !prev),
        addItem,
        incrementItem,
        decrementItem,
        removeItem,
        clear,
        maxQuantity: maxQuantityFor,
    };

    return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
}

export function useCart() {
    const ctx = useContext(CartContext);
    if (!ctx) {
        throw new Error('useCart must be used within a CartProvider');
    }
    return ctx;
}

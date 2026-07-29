'use client';

import { createContext, useContext, useEffect, useState, useCallback, ReactNode } from 'react';
import { CartItem } from '../../domain/types';

interface CartContextValue {
    items: CartItem[];
    addItem: (item: Omit<CartItem, 'quantity'>, quantity?: number) => void;
    removeItem: (productId: string) => void;
    updateQuantity: (productId: string, quantity: number) => void;
    clear: () => void;
    total: number;
    count: number;
    isOpen: boolean;
    open: () => void;
    close: () => void;
}

const CartContext = createContext<CartContextValue | null>(null);

function storageKey(slug: string) {
    return `probability_cart_${slug}`;
}

export function CartProvider({ slug, children }: { slug: string; children: ReactNode }) {
    const [items, setItems] = useState<CartItem[]>([]);
    const [isOpen, setIsOpen] = useState(false);
    const [hydrated, setHydrated] = useState(false);

    useEffect(() => {
        try {
            const raw = window.localStorage.getItem(storageKey(slug));
            if (raw) setItems(JSON.parse(raw));
        } catch {
            // ignore malformed cart in storage
        }
        setHydrated(true);
    }, [slug]);

    useEffect(() => {
        if (!hydrated) return;
        window.localStorage.setItem(storageKey(slug), JSON.stringify(items));
    }, [items, slug, hydrated]);

    const addItem = useCallback((item: Omit<CartItem, 'quantity'>, quantity = 1) => {
        setItems((prev) => {
            const existing = prev.find((i) => i.product_id === item.product_id);
            if (existing) {
                return prev.map((i) =>
                    i.product_id === item.product_id ? { ...i, quantity: i.quantity + quantity } : i
                );
            }
            return [...prev, { ...item, quantity }];
        });
        setIsOpen(true);
    }, []);

    const removeItem = useCallback((productId: string) => {
        setItems((prev) => prev.filter((i) => i.product_id !== productId));
    }, []);

    const updateQuantity = useCallback((productId: string, quantity: number) => {
        if (quantity < 1) {
            setItems((prev) => prev.filter((i) => i.product_id !== productId));
            return;
        }
        setItems((prev) => prev.map((i) => (i.product_id === productId ? { ...i, quantity } : i)));
    }, []);

    const clear = useCallback(() => setItems([]), []);

    const total = items.reduce((sum, i) => sum + i.price * i.quantity, 0);
    const count = items.reduce((sum, i) => sum + i.quantity, 0);

    return (
        <CartContext.Provider
            value={{
                items,
                addItem,
                removeItem,
                updateQuantity,
                clear,
                total,
                count,
                isOpen,
                open: () => setIsOpen(true),
                close: () => setIsOpen(false),
            }}
        >
            {children}
        </CartContext.Provider>
    );
}

export function useCart() {
    const ctx = useContext(CartContext);
    if (!ctx) throw new Error('useCart debe usarse dentro de CartProvider');
    return ctx;
}

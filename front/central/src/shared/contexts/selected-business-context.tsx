'use client';

import { createContext, useContext, useEffect, useState, type ReactNode } from 'react';
import { usePermissions } from './permissions-context';

const STORAGE_KEY = 'selected_business_id';
const COOKIE_MAX_AGE = 60 * 60 * 24 * 30;

interface SelectedBusinessContextType {
    selectedBusinessId: number | null;
    setSelectedBusinessId: (id: number | null) => void;
}

const SelectedBusinessContext = createContext<SelectedBusinessContextType | null>(null);

function readStoredBusinessId(): number | null {
    if (typeof window === 'undefined') return null;

    const stored = window.localStorage.getItem(STORAGE_KEY);
    if (stored) {
        const parsed = parseInt(stored, 10);
        if (!Number.isNaN(parsed) && parsed > 0) return parsed;
    }

    const match = document.cookie.match(/selected_business_id=(\d+)/);
    return match ? parseInt(match[1], 10) : null;
}

function persistBusinessId(id: number | null) {
    if (typeof window === 'undefined') return;

    if (id) {
        window.localStorage.setItem(STORAGE_KEY, String(id));
        document.cookie = `selected_business_id=${id}; path=/; max-age=${COOKIE_MAX_AGE}`;
        document.cookie = `storefront_business_id=${id}; path=/; max-age=${COOKIE_MAX_AGE}`;
        return;
    }

    window.localStorage.removeItem(STORAGE_KEY);
    window.localStorage.removeItem('selected_business_name');
    window.localStorage.removeItem('selected_business_logo');
    document.cookie = 'selected_business_id=; path=/; max-age=0';
    document.cookie = 'storefront_business_id=; path=/; max-age=0';
}

export function SelectedBusinessProvider({ children }: { children: ReactNode }) {
    const { permissions, isSuperAdmin } = usePermissions();
    const [selectedBusinessId, setSelectedBusinessIdState] = useState<number | null>(null);
    const ownBusinessId = permissions?.business_id ?? null;

    useEffect(() => {
        if (isSuperAdmin) {
            setSelectedBusinessIdState(readStoredBusinessId());
            return;
        }

        if (!permissions) return;

        const own = ownBusinessId && ownBusinessId > 0 ? ownBusinessId : null;
        setSelectedBusinessIdState(own);
        if (own !== readStoredBusinessId()) persistBusinessId(own);
    }, [isSuperAdmin, permissions, ownBusinessId]);

    const setSelectedBusinessId = (id: number | null) => {
        if (!isSuperAdmin) return;
        setSelectedBusinessIdState(id);
        persistBusinessId(id);
    };

    return (
        <SelectedBusinessContext.Provider value={{ selectedBusinessId, setSelectedBusinessId }}>
            {children}
        </SelectedBusinessContext.Provider>
    );
}

export function useSelectedBusiness() {
    const ctx = useContext(SelectedBusinessContext);
    if (!ctx) {
        throw new Error('useSelectedBusiness must be used within SelectedBusinessProvider');
    }
    return ctx;
}

'use client';

import { useEffect, useRef, useState } from 'react';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useBusinessesSearch } from '@/services/auth/business/ui/hooks/useBusinessesSearch';
import { BusinessSimple } from '@/services/auth/business/domain/types';
import { applyBusinessTheme, resetTheme } from '@/shared/utils/apply-business-theme';

interface SuperAdminBusinessSelectorProps {
    value: number | null;
    onChange: (businessId: number | null) => void;
    variant?: 'navbar' | 'default';
    placeholder?: string;
}

export function SuperAdminBusinessSelector({
    value,
    onChange,
    variant = 'default',
    placeholder = 'Todos los negocios',
}: SuperAdminBusinessSelectorProps) {
    const { isSuperAdmin } = usePermissions();
    const {
        businesses,
        total,
        loading,
        search,
        setSearch,
        hasMore,
        loadMore,
        minSearchLength,
    } = useBusinessesSearch();

    const [open, setOpen] = useState(false);
    const [selected, setSelected] = useState<BusinessSimple | null>(null);
    const containerRef = useRef<HTMLDivElement>(null);
    const searchInputRef = useRef<HTMLInputElement>(null);

    useEffect(() => {
        if (value === null) {
            setSelected(null);
            return;
        }
        if (selected?.id === value) return;
        const match = businesses.find(b => b.id === value);
        if (match) setSelected(match);
    }, [value, businesses, selected]);

    useEffect(() => {
        const handler = (e: MouseEvent) => {
            if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
                setOpen(false);
            }
        };
        document.addEventListener('mousedown', handler);
        return () => document.removeEventListener('mousedown', handler);
    }, []);

    useEffect(() => {
        if (open) searchInputRef.current?.focus();
    }, [open]);

    if (!isSuperAdmin) {
        return null;
    }

    const handleSelect = (business: BusinessSimple | null) => {
        setSelected(business);
        setOpen(false);
        setSearch('');
        onChange(business ? business.id : null);

        if (business) {
            applyBusinessTheme({
                name: business.name,
                logo_url: business.logo_url,
                primary_color: business.primary_color,
                secondary_color: business.secondary_color,
                tertiary_color: business.tertiary_color,
                quaternary_color: business.quaternary_color,
            });
        } else {
            resetTheme();
        }
    };

    const trimmedSearch = search.trim();
    const showMinSearchHint = trimmedSearch.length > 0 && trimmedSearch.length < minSearchLength;
    const label = selected ? selected.name : (value ? `Negocio #${value}` : placeholder);

    const isNavbar = variant === 'navbar';

    const triggerClasses = isNavbar
        ? 'flex items-center justify-between gap-2 min-w-44 px-2 py-1.5 border border-purple-400 dark:border-purple-600 rounded-md text-sm font-medium focus:outline-none focus:ring-2 focus:ring-purple-600 bg-white dark:bg-gray-800 text-purple-900 dark:text-purple-200 cursor-pointer'
        : 'flex items-center justify-between gap-2 flex-1 max-w-xs px-3 py-2 border-2 border-purple-400 dark:border-purple-600 rounded-lg text-sm font-medium focus:outline-none focus:ring-2 focus:ring-purple-600 bg-white dark:bg-gray-800 text-purple-900 dark:text-purple-200 cursor-pointer';

    const wrapperClasses = isNavbar
        ? 'flex items-center gap-2 bg-purple-100 dark:bg-purple-900/30 border border-purple-300 dark:border-purple-700 rounded-lg px-3 py-1.5'
        : 'flex items-center gap-3 bg-purple-50 dark:bg-purple-900/20 border-2 border-purple-300 dark:border-purple-700 rounded-lg px-4 py-3';

    const badgeClasses = isNavbar
        ? 'px-2 py-0.5 text-xs font-bold text-white bg-purple-700 rounded select-none whitespace-nowrap'
        : 'px-2.5 py-1 text-xs font-bold text-white bg-purple-700 rounded-md select-none whitespace-nowrap';

    return (
        <div className={wrapperClasses}>
            <span className={badgeClasses}>SUPER ADMIN</span>
            <div ref={containerRef} className="relative">
                <button
                    type="button"
                    onClick={() => setOpen(prev => !prev)}
                    className={triggerClasses}
                >
                    <span className="truncate">{label}</span>
                    <svg className="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                    </svg>
                </button>

                {open && (
                    <div className="absolute right-0 mt-1 w-72 z-50 bg-white dark:bg-gray-800 border border-purple-300 dark:border-purple-700 rounded-lg shadow-lg overflow-hidden">
                        <div className="p-2 border-b border-purple-200 dark:border-purple-800">
                            <input
                                ref={searchInputRef}
                                type="text"
                                value={search}
                                onChange={(e) => setSearch(e.target.value)}
                                placeholder={`Buscar negocio (min. ${minSearchLength} letras)`}
                                className="w-full px-2 py-1.5 text-sm border border-purple-300 dark:border-purple-600 rounded-md bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-600"
                            />
                            {showMinSearchHint && (
                                <p className="mt-1 text-xs text-purple-600 dark:text-purple-400">
                                    {`Escribe al menos ${minSearchLength} letras para buscar`}
                                </p>
                            )}
                        </div>

                        <ul className="max-h-64 overflow-y-auto">
                            <li>
                                <button
                                    type="button"
                                    onClick={() => handleSelect(null)}
                                    className={`w-full text-left px-3 py-2 text-sm hover:bg-purple-50 dark:hover:bg-purple-900/30 ${value === null ? 'bg-purple-100 dark:bg-purple-900/40 font-semibold text-purple-900 dark:text-purple-200' : 'text-gray-800 dark:text-gray-200'}`}
                                >
                                    {placeholder}
                                </button>
                            </li>
                            {businesses.map((b) => (
                                <li key={b.id}>
                                    <button
                                        type="button"
                                        onClick={() => handleSelect(b)}
                                        className={`w-full text-left px-3 py-2 text-sm hover:bg-purple-50 dark:hover:bg-purple-900/30 ${value === b.id ? 'bg-purple-100 dark:bg-purple-900/40 font-semibold text-purple-900 dark:text-purple-200' : 'text-gray-800 dark:text-gray-200'}`}
                                    >
                                        {b.name}
                                    </button>
                                </li>
                            ))}
                            {!loading && businesses.length === 0 && (
                                <li className="px-3 py-3 text-sm text-gray-500 dark:text-gray-400 text-center">
                                    Sin resultados
                                </li>
                            )}
                            {loading && (
                                <li className="px-3 py-2 text-sm text-gray-500 dark:text-gray-400 text-center">
                                    Cargando...
                                </li>
                            )}
                        </ul>

                        {hasMore && !loading && (
                            <div className="p-2 border-t border-purple-200 dark:border-purple-800">
                                <button
                                    type="button"
                                    onClick={loadMore}
                                    className="w-full px-3 py-1.5 text-sm font-medium text-purple-700 dark:text-purple-300 hover:bg-purple-50 dark:hover:bg-purple-900/30 rounded-md"
                                >
                                    {`Cargar más (${businesses.length} de ${total})`}
                                </button>
                            </div>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
}

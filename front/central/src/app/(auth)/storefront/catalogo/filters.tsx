'use client';

import { useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { FunnelIcon, XMarkIcon } from '@heroicons/react/24/outline';
import { StorefrontFamily } from '@/services/modules/storefront/domain/types';

interface CatalogFiltersProps {
    categories: string[];
    families: StorefrontFamily[];
    selectedCategory: string;
    selectedFamilyId: number | null;
}

function Chip({ label, active, onClick }: { label: string; active: boolean; onClick: () => void }) {
    return (
        <button
            type="button"
            onClick={onClick}
            className={`px-3.5 py-1.5 rounded-full text-xs font-semibold whitespace-nowrap transition-colors ${
                active
                    ? 'bg-indigo-600 text-white'
                    : 'bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700'
            }`}
        >
            {label}
        </button>
    );
}

export function CatalogFilters({ categories, families, selectedCategory, selectedFamilyId }: CatalogFiltersProps) {
    const router = useRouter();
    const searchParams = useSearchParams();
    const [open, setOpen] = useState(false);

    const topLevelFamilies = families.filter(f => f.parent_family_id === null);
    const childrenByParent = new Map<number, StorefrontFamily[]>();
    families.forEach(f => {
        if (f.parent_family_id !== null) {
            const list = childrenByParent.get(f.parent_family_id) || [];
            list.push(f);
            childrenByParent.set(f.parent_family_id, list);
        }
    });

    const selectedFamily = families.find(f => f.id === selectedFamilyId) || null;
    const activeParentId = selectedFamily
        ? (selectedFamily.parent_family_id ?? selectedFamily.id)
        : null;
    const subFamilies = activeParentId !== null ? childrenByParent.get(activeParentId) || [] : [];

    const activeFilterCount = (selectedCategory ? 1 : 0) + (selectedFamilyId ? 1 : 0);

    const updateParams = (updates: Record<string, string | null>) => {
        const params = new URLSearchParams(searchParams);
        Object.entries(updates).forEach(([key, value]) => {
            if (value === null || value === '') params.delete(key);
            else params.set(key, value);
        });
        params.set('page', '1');
        router.push(`/storefront/catalogo?${params.toString()}`);
    };

    const clearAll = () => updateParams({ category: null, family_id: null });

    if (categories.length === 0 && topLevelFamilies.length === 0) {
        return null;
    }

    return (
        <div className="relative">
            <button
                type="button"
                onClick={() => setOpen(prev => !prev)}
                className="flex items-center gap-2 px-4 py-2 h-[42px] text-sm font-medium border border-gray-300 dark:border-gray-600 rounded-lg text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
            >
                <FunnelIcon className="w-4 h-4" />
                Filtros
                {activeFilterCount > 0 && (
                    <span className="flex items-center justify-center w-5 h-5 rounded-full bg-indigo-600 text-white text-[11px] font-bold">
                        {activeFilterCount}
                    </span>
                )}
            </button>

            {open && (
                <>
                    <div className="fixed inset-0 z-30" onClick={() => setOpen(false)} />
                    <div className="absolute left-0 mt-2 w-[22rem] max-h-[70vh] overflow-y-auto bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-xl z-40 p-4">
                        <div className="flex items-center justify-between mb-3">
                            <h3 className="text-sm font-bold text-gray-900 dark:text-white">Filtrar productos</h3>
                            <button onClick={() => setOpen(false)} className="text-gray-400 hover:text-gray-600">
                                <XMarkIcon className="w-4 h-4" />
                            </button>
                        </div>

                        {categories.length > 0 && (
                            <div className="mb-4">
                                <p className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-2">Categoria</p>
                                <div className="flex flex-wrap gap-2">
                                    <Chip label="Todas" active={!selectedCategory} onClick={() => updateParams({ category: null })} />
                                    {categories.map(category => (
                                        <Chip
                                            key={category}
                                            label={category}
                                            active={selectedCategory === category}
                                            onClick={() => updateParams({ category })}
                                        />
                                    ))}
                                </div>
                            </div>
                        )}

                        {topLevelFamilies.length > 0 && (
                            <div className="mb-2">
                                <p className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-2">Familia</p>
                                <div className="flex flex-wrap gap-2">
                                    <Chip label="Todas" active={!activeParentId} onClick={() => updateParams({ family_id: null })} />
                                    {topLevelFamilies.map(family => (
                                        <Chip
                                            key={family.id}
                                            label={family.name}
                                            active={activeParentId === family.id}
                                            onClick={() => updateParams({ family_id: String(family.id) })}
                                        />
                                    ))}
                                </div>
                            </div>
                        )}

                        {subFamilies.length > 0 && (
                            <div className="mb-2 pl-3 border-l-2 border-gray-100 dark:border-gray-700">
                                <p className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-2">Subfamilia</p>
                                <div className="flex flex-wrap gap-2">
                                    <Chip
                                        label="Todas"
                                        active={selectedFamilyId === activeParentId}
                                        onClick={() => updateParams({ family_id: String(activeParentId) })}
                                    />
                                    {subFamilies.map(sub => (
                                        <Chip
                                            key={sub.id}
                                            label={sub.name}
                                            active={selectedFamilyId === sub.id}
                                            onClick={() => updateParams({ family_id: String(sub.id) })}
                                        />
                                    ))}
                                </div>
                            </div>
                        )}

                        {activeFilterCount > 0 && (
                            <button
                                onClick={clearAll}
                                className="mt-3 text-xs font-medium text-gray-500 dark:text-gray-400 hover:text-red-500 transition-colors"
                            >
                                Limpiar filtros
                            </button>
                        )}
                    </div>
                </>
            )}
        </div>
    );
}

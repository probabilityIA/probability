'use client';

import { useState, useEffect } from 'react';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import { Spinner } from '@/shared/ui';

interface BusinessTargetSelectorProps {
    selectedIds: number[];
    onChange: (ids: number[]) => void;
}

export default function BusinessTargetSelector({ selectedIds, onChange }: BusinessTargetSelectorProps) {
    const { businesses, loading } = useBusinessesSimple();
    const [search, setSearch] = useState('');

    const filtered = businesses.filter(b =>
        b.name.toLowerCase().includes(search.toLowerCase())
    );

    const toggleBusiness = (id: number) => {
        if (selectedIds.includes(id)) {
            onChange(selectedIds.filter(sid => sid !== id));
        } else {
            onChange([...selectedIds, id]);
        }
    };

    const selectAll = () => {
        onChange(filtered.map(b => b.id));
    };

    const clearAll = () => {
        onChange([]);
    };

    if (loading) {
        return (
            <div className="flex justify-center p-4">
                <Spinner size="sm" />
            </div>
        );
    }

    return (
        <div className="space-y-2">
            <div className="flex items-center gap-2">
                <input
                    type="text"
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    placeholder="Buscar negocio..."
                    className="flex-1 px-3 py-1.5 border border-gray-300 dark:border-gray-600 rounded-lg text-sm text-gray-900 dark:text-white bg-white dark:bg-gray-700 placeholder-gray-500 dark:placeholder-white focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                />
                <button
                    type="button"
                    onClick={selectAll}
                    className="text-xs text-purple-600 hover:text-purple-700 whitespace-nowrap"
                >
                    Todos
                </button>
                <button
                    type="button"
                    onClick={clearAll}
                    className="text-xs text-gray-500 hover:text-gray-700 whitespace-nowrap"
                >
                    Ninguno
                </button>
            </div>

            <div className="max-h-40 overflow-y-auto border border-gray-200 dark:border-gray-700 rounded-lg divide-y divide-gray-100 dark:divide-gray-700">
                {filtered.length === 0 ? (
                    <p className="text-sm text-gray-400 text-center py-3">Sin resultados</p>
                ) : (
                    filtered.map(business => (
                        <label
                            key={business.id}
                            className="flex items-center gap-2 px-3 py-2 hover:bg-gray-50 dark:hover:bg-gray-700/50 cursor-pointer"
                        >
                            <input
                                type="checkbox"
                                checked={selectedIds.includes(business.id)}
                                onChange={() => toggleBusiness(business.id)}
                                className="w-4 h-4 rounded border-gray-300 text-purple-600 focus:ring-purple-500"
                            />
                            <span className="text-sm text-gray-700 dark:text-gray-300">{business.name}</span>
                        </label>
                    ))
                )}
            </div>

            {selectedIds.length > 0 && (
                <p className="text-xs text-gray-500">{selectedIds.length} negocio(s) seleccionado(s)</p>
            )}
        </div>
    );
}

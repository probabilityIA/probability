'use client';

import { useState, useEffect, useCallback } from 'react';
import { IntegrationCategory, IntegrationType } from '../../domain/types';
import { getActiveIntegrationTypesAction } from '../../infra/actions';
import { Spinner } from '@/shared/ui';
import { getActionError } from '@/shared/utils/action-result';

interface ProviderSelectorProps {
    category: IntegrationCategory;
    onSelect: (provider: IntegrationType) => void;
    onBack: () => void;
}

export function ProviderSelector({ category, onSelect, onBack }: ProviderSelectorProps) {
    const [providers, setProviders] = useState<IntegrationType[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchProviders = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getActiveIntegrationTypesAction();
            if (response.success && response.data) {
                const filtered = response.data.filter(
                    (provider) => provider.category?.id === category.id || provider.category_id === category.id
                );
                setProviders(filtered);
            } else {
                setError('Error al cargar proveedores');
            }
        } catch (err) {
            setError(getActionError(err, 'Error al cargar proveedores'));
        } finally {
            setLoading(false);
        }
    }, [category]);

    useEffect(() => {
        fetchProviders();
    }, [fetchProviders]);

    if (loading) {
        return (
            <div className="flex justify-center p-10">
                <Spinner />
            </div>
        );
    }

    if (error) {
        return (
            <div className="p-4">
                <div className="text-center text-red-600">{error}</div>
                <div className="mt-3 text-center">
                    <button
                        onClick={onBack}
                        className="text-sm text-blue-600 hover:text-blue-800"
                    >
                        &larr; Volver a categorias
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className="w-[42rem] max-w-[88vw] p-2">
            <div className="mb-4">
                <button
                    onClick={onBack}
                    className="mb-2 flex items-center gap-1.5 text-sm text-blue-600 hover:text-blue-800"
                >
                    <span>&larr;</span>
                    <span>Volver a categorias</span>
                </button>

                <h2 className="text-base font-bold text-gray-900 dark:text-white">
                    {category.name}
                </h2>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    Selecciona el proveedor que deseas integrar
                </p>
            </div>

            <div className="grid grid-cols-1 gap-2.5 md:grid-cols-2">
                {providers.map((provider) => {
                    const isDev = provider.in_development === true;
                    return (
                        <button
                            key={provider.id}
                            onClick={() => !isDev && onSelect(provider)}
                            disabled={isDev}
                            className={`group relative rounded-lg border-2 p-3 text-left transition-all ${
                                isDev
                                    ? 'cursor-not-allowed border-gray-200 opacity-60 grayscale dark:border-gray-700'
                                    : 'border-gray-200 hover:border-blue-500 hover:shadow-md dark:border-gray-700'
                            }`}
                        >
                            {isDev && (
                                <span className="absolute right-2 top-2 rounded-full bg-amber-100 px-2 py-0.5 text-[10px] font-semibold text-amber-800">
                                    Proximamente
                                </span>
                            )}

                            <div className="flex items-center gap-3">
                                <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center">
                                    {provider.image_url ? (
                                        <img
                                            src={provider.image_url}
                                            alt={provider.name}
                                            className="max-h-full max-w-full object-contain"
                                        />
                                    ) : (
                                        <span className="text-2xl">{provider.icon}</span>
                                    )}
                                </div>

                                <div className="min-w-0 flex-1">
                                    <h3 className={`text-sm font-semibold ${
                                        isDev
                                            ? 'text-gray-500 dark:text-gray-400'
                                            : 'text-gray-900 group-hover:text-blue-600 dark:text-white'
                                    }`}>
                                        {provider.name}
                                    </h3>
                                    {provider.description && (
                                        <p className="truncate text-xs text-gray-500 dark:text-gray-400">
                                            {provider.description}
                                        </p>
                                    )}
                                </div>
                            </div>
                        </button>
                    );
                })}
            </div>

            {providers.length === 0 && (
                <div className="py-8 text-center text-sm text-gray-500 dark:text-gray-400">
                    No hay proveedores disponibles para esta categoria
                </div>
            )}
        </div>
    );
}

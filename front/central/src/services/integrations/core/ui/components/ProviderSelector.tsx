'use client';

import { useState, useEffect } from 'react';
import { IntegrationCategory, IntegrationType } from '../../domain/types';
import { getActiveIntegrationTypesAction } from '../../infra/actions';
import { Spinner } from '@/shared/ui';

interface ProviderSelectorProps {
    category: IntegrationCategory;
    onSelect: (provider: IntegrationType) => void;
    onBack: () => void;
}

export function ProviderSelector({ category, onSelect, onBack }: ProviderSelectorProps) {
    const [providers, setProviders] = useState<IntegrationType[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        fetchProviders();
    }, [category]);

    const fetchProviders = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getActiveIntegrationTypesAction();
            console.log('[ProviderSelector] Response from API:', response);
            console.log('[ProviderSelector] Selected category:', category);
            if (response.success && response.data) {
                console.log('[ProviderSelector] All providers:', response.data);
                // Filter providers by category
                // El backend devuelve category como objeto, no category_id
                const filtered = response.data.filter(
                    (provider) => {
                        console.log(`[ProviderSelector] Provider ${provider.name}:`, {
                            category: provider.category,
                            category_id: provider.category_id,
                            matches: provider.category?.id === category.id || provider.category_id === category.id
                        });
                        return provider.category?.id === category.id || provider.category_id === category.id;
                    }
                );
                console.log('[ProviderSelector] Filtered providers:', filtered);
                setProviders(filtered);
            } else {
                setError('Error al cargar proveedores');
            }
        } catch (err: any) {
            setError(err.message || 'Error al cargar proveedores');
        } finally {
            setLoading(false);
        }
    };

    if (loading) {
        return (
            <div className="p-12 flex justify-center">
                <Spinner />
            </div>
        );
    }

    if (error) {
        return (
            <div className="p-6">
                <div className="text-red-600 text-center">{error}</div>
                <div className="mt-4 text-center">
                    <button
                        onClick={onBack}
                        className="text-blue-600 hover:text-blue-800"
                    >
                        ← Volver a categorías
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className="p-6">
            {/* Header with back button */}
            <div className="mb-6">
                <button
                    onClick={onBack}
                    className="text-blue-600 hover:text-blue-800 mb-4 flex items-center gap-2"
                >
                    <span>←</span>
                    <span>Volver a categorías</span>
                </button>

                <h2 className="text-2xl font-bold text-gray-900 mb-2">
                    {category.name}
                </h2>
                <p className="text-gray-600">
                    Selecciona el proveedor que deseas integrar
                </p>
            </div>

            {/* Providers Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {providers.map((provider) => {
                    const isDev = provider.in_development === true;
                    return (
                        <button
                            key={provider.id}
                            onClick={() => !isDev && onSelect(provider)}
                            disabled={isDev}
                            className={`p-6 border-2 rounded-lg transition-all text-left group relative ${
                                isDev
                                    ? 'border-gray-200 opacity-60 grayscale cursor-not-allowed'
                                    : 'border-gray-200 hover:border-blue-500 hover:shadow-md'
                            }`}
                        >
                            {/* Coming Soon Badge */}
                            {isDev && (
                                <span className="absolute top-2 right-2 bg-amber-100 text-amber-800 text-xs font-semibold px-2 py-1 rounded-full">
                                    Proximamente
                                </span>
                            )}

                            {/* Provider Logo */}
                            {provider.image_url && (
                                <div className="mb-4 h-16 flex items-center justify-center">
                                    <img
                                        src={provider.image_url}
                                        alt={provider.name}
                                        className="max-h-full max-w-full object-contain"
                                    />
                                </div>
                            )}

                            {/* Provider Icon (fallback if no image) */}
                            {!provider.image_url && provider.icon && (
                                <div className="mb-4 text-3xl text-center">
                                    {provider.icon}
                                </div>
                            )}

                            {/* Provider Name */}
                            <h3 className={`text-lg font-semibold mb-2 text-center ${
                                isDev
                                    ? 'text-gray-500'
                                    : 'text-gray-900 group-hover:text-blue-600'
                            }`}>
                                {provider.name}
                            </h3>

                            {/* Description */}
                            {provider.description && (
                                <p className="text-sm text-gray-600 text-center">
                                    {provider.description}
                                </p>
                            )}
                        </button>
                    );
                })}
            </div>

            {providers.length === 0 && (
                <div className="text-center py-12 text-gray-500">
                    No hay proveedores disponibles para esta categoría
                </div>
            )}
        </div>
    );
}

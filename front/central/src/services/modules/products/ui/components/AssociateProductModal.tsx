'use client';

import { useState, useEffect, useCallback } from 'react';
import { Product, ProductFamily } from '../../domain/types';
import { getProductsAction, updateProductAction } from '../../infra/actions';
import { XMarkIcon, MagnifyingGlassIcon } from '@heroicons/react/24/outline';

interface AssociateProductModalProps {
    businessId?: number;
    family: ProductFamily;
    onAssociated: () => void;
    onClose: () => void;
}

export default function AssociateProductModal({ businessId, family, onAssociated, onClose }: AssociateProductModalProps) {
    const [searchTerm, setSearchTerm] = useState('');
    const [results, setResults] = useState<Product[]>([]);
    const [searching, setSearching] = useState(false);
    const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
    const [variantLabel, setVariantLabel] = useState('');
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const search = useCallback(async (term: string) => {
        if (!term) {
            setResults([]);
            return;
        }
        setSearching(true);
        try {
            const res = await getProductsAction({ business_id: businessId, name: term, page: 1, page_size: 10 });
            setResults(res.success ? res.data : []);
        } finally {
            setSearching(false);
        }
    }, [businessId]);

    useEffect(() => {
        const timer = setTimeout(() => search(searchTerm), 500);
        return () => clearTimeout(timer);
    }, [searchTerm, search]);

    const handleSelect = (product: Product) => {
        setSelectedProduct(product);
        setVariantLabel(product.variant_label || '');
        setError(null);
    };

    const handleConfirm = async () => {
        if (!selectedProduct || !variantLabel.trim()) return;
        setSaving(true);
        setError(null);
        try {
            const res = await updateProductAction(
                selectedProduct.id,
                {
                    family_id: family.id,
                    variant_label: variantLabel.trim(),
                    variant_attributes: { variante: variantLabel.trim() },
                },
                businessId
            );
            if (res.success === false) {
                setError(res.message || 'Error al asociar el producto');
                return;
            }
            onAssociated();
            onClose();
        } catch (e: any) {
            setError(e.message || 'Error inesperado al asociar el producto');
        } finally {
            setSaving(false);
        }
    };

    return (
        <div className="fixed inset-0 z-[60] flex items-center justify-center p-4">
            <div className="absolute inset-0 bg-black/50" onClick={onClose} />
            <div className="relative bg-white dark:bg-gray-800 rounded-xl shadow-2xl w-full max-w-lg max-h-[85vh] flex flex-col">
                <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex-shrink-0">
                    <div>
                        <h3 className="text-base font-semibold text-gray-900 dark:text-white">Asociar producto</h3>
                        <p className="text-xs text-gray-500 mt-0.5">a la familia &quot;{family.name}&quot;</p>
                    </div>
                    <button
                        onClick={onClose}
                        className="p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                    >
                        <XMarkIcon className="w-5 h-5" />
                    </button>
                </div>

                <div className="overflow-auto flex-1 p-5 flex flex-col gap-4">
                    {!selectedProduct ? (
                        <>
                            <div className="relative">
                                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                    <MagnifyingGlassIcon className="h-4 w-4 text-gray-400" />
                                </div>
                                <input
                                    type="text"
                                    autoFocus
                                    placeholder="Buscar producto por nombre o SKU..."
                                    value={searchTerm}
                                    onChange={(e) => setSearchTerm(e.target.value)}
                                    className="w-full pl-9 pr-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder:text-gray-400 focus:ring-2 focus:ring-[var(--color-primary)] focus:border-transparent"
                                />
                            </div>

                            <div className="border border-gray-200 dark:border-gray-700 rounded-lg divide-y divide-gray-100 dark:divide-gray-700 max-h-72 overflow-auto">
                                {searching ? (
                                    <div className="p-4 text-center text-sm text-gray-500 dark:text-gray-400">Buscando...</div>
                                ) : results.length > 0 ? (
                                    results.map((product) => (
                                        <button
                                            key={product.id}
                                            type="button"
                                            onClick={() => handleSelect(product)}
                                            className="w-full px-4 py-2.5 text-left hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors flex items-center justify-between gap-2"
                                        >
                                            <div>
                                                <div className="text-sm font-medium text-gray-900 dark:text-white">{product.name}</div>
                                                <div className="text-xs text-gray-500 dark:text-gray-400">SKU: {product.sku}</div>
                                            </div>
                                            {product.family && (
                                                <span className="shrink-0 px-2 py-0.5 rounded-full text-xs font-medium bg-amber-100 text-amber-700">
                                                    Ya en: {product.family.name}
                                                </span>
                                            )}
                                        </button>
                                    ))
                                ) : searchTerm ? (
                                    <div className="p-4 text-center text-sm text-gray-400">No se encontraron productos.</div>
                                ) : (
                                    <div className="p-4 text-center text-sm text-gray-400">Escribe para buscar un producto.</div>
                                )}
                            </div>
                        </>
                    ) : (
                        <>
                            <div className="p-3 rounded-lg bg-gray-50 dark:bg-gray-700/50 border border-gray-200 dark:border-gray-700">
                                <div className="text-sm font-medium text-gray-900 dark:text-white">{selectedProduct.name}</div>
                                <div className="text-xs text-gray-500 dark:text-gray-400">SKU: {selectedProduct.sku}</div>
                                {selectedProduct.family && (
                                    <div className="text-xs text-amber-600 mt-1">
                                        Se reasignará desde &quot;{selectedProduct.family.name}&quot; a &quot;{family.name}&quot;
                                    </div>
                                )}
                            </div>

                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">
                                    Variante dentro de la familia
                                </label>
                                <input
                                    type="text"
                                    autoFocus
                                    placeholder="Ej: Talla M, Rojo..."
                                    value={variantLabel}
                                    onChange={(e) => setVariantLabel(e.target.value)}
                                    className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder:text-gray-400 focus:ring-2 focus:ring-[var(--color-primary)] focus:border-transparent"
                                />
                            </div>

                            {error && (
                                <div className="px-3 py-2 rounded-lg bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 text-xs">
                                    {error}
                                </div>
                            )}

                            <div className="flex items-center gap-2">
                                <button
                                    type="button"
                                    onClick={() => { setSelectedProduct(null); setError(null); }}
                                    className="px-3 py-2 text-sm rounded-lg border border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                                >
                                    Volver
                                </button>
                                <button
                                    type="button"
                                    disabled={!variantLabel.trim() || saving}
                                    onClick={handleConfirm}
                                    className="flex-1 px-3 py-2 text-sm rounded-lg bg-[var(--color-primary)] text-white font-medium disabled:opacity-50 hover:opacity-90 transition-opacity"
                                >
                                    {saving ? 'Asociando...' : 'Asociar'}
                                </button>
                            </div>
                        </>
                    )}
                </div>
            </div>
        </div>
    );
}

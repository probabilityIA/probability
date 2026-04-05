'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import { AdjustStockDTO } from '../../domain/types';
import { adjustStockAction } from '../../infra/actions';
import { getProductsAction } from '@/services/modules/products/infra/actions';
import { Button, Alert, Input } from '@/shared/ui';

interface ProductOption {
    id: string;
    name: string;
    sku: string;
}

interface AdjustStockModalProps {
    warehouseId: number;
    businessId?: number;
    productId?: string;
    onSuccess: () => void;
    onClose: () => void;
}

export default function AdjustStockModal({ warehouseId, businessId, productId, onSuccess, onClose }: AdjustStockModalProps) {
    const [selectedProduct, setSelectedProduct] = useState<ProductOption | null>(null);
    const [searchTerm, setSearchTerm] = useState('');
    const [searchResults, setSearchResults] = useState<ProductOption[]>([]);
    const [showDropdown, setShowDropdown] = useState(false);
    const [searchLoading, setSearchLoading] = useState(false);
    const dropdownRef = useRef<HTMLDivElement>(null);

    const [quantity, setQuantity] = useState(0);
    const [reason, setReason] = useState('');
    const [notes, setNotes] = useState('');

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    // Si viene un productId pre-seleccionado, cargarlo
    useEffect(() => {
        if (productId) {
            setSelectedProduct({ id: productId, name: '', sku: '' });
            // Intentar obtener datos del producto
            getProductsAction({ business_id: businessId, page: 1, page_size: 1 }).then(() => {
                // El productId ya está seteado, no necesitamos más
            });
        }
    }, [productId, businessId]);

    // Buscar productos con debounce
    const searchProducts = useCallback(async (term: string) => {
        if (term.length < 2) {
            setSearchResults([]);
            return;
        }
        setSearchLoading(true);
        try {
            const response = await getProductsAction({
                business_id: businessId,
                name: term,
                page: 1,
                page_size: 10,
            });
            if (response.success && response.data) {
                setSearchResults(response.data.map(p => ({ id: p.id, name: p.name, sku: p.sku })));
            }
        } catch {
            // silently fail
        } finally {
            setSearchLoading(false);
        }
    }, [businessId]);

    useEffect(() => {
        const timer = setTimeout(() => {
            if (searchTerm.length >= 2) searchProducts(searchTerm);
        }, 400);
        return () => clearTimeout(timer);
    }, [searchTerm, searchProducts]);

    // Cerrar dropdown al hacer click fuera
    useEffect(() => {
        const handleClickOutside = (e: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
                setShowDropdown(false);
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const handleSelectProduct = (product: ProductOption) => {
        setSelectedProduct(product);
        setSearchTerm('');
        setShowDropdown(false);
        setSearchResults([]);
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedProduct) return;

        setLoading(true);
        setError(null);
        setSuccess(null);

        const dto: AdjustStockDTO = {
            product_id: selectedProduct.id,
            warehouse_id: warehouseId,
            quantity,
            reason: reason.trim(),
            notes: notes.trim() || undefined,
        };
        const result = await adjustStockAction(dto, businessId);
        if (!result.success) {
            setError(result.error);
        } else {
            setSuccess('Stock ajustado exitosamente');
            setTimeout(() => onSuccess(), 800);
        }
        setLoading(false);
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-md max-h-[90vh] overflow-y-auto">
                <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700">
                    <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Ajustar stock</h2>
                    <button onClick={onClose} className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 text-xl leading-none">
                        &times;
                    </button>
                </div>
                <form onSubmit={handleSubmit} className="p-6 space-y-4">
                    {error && (
                        <Alert type="error" onClose={() => setError(null)}>
                            {error}
                        </Alert>
                    )}
                    {success && (
                        <Alert type="success" onClose={() => setSuccess(null)}>
                            {success}
                        </Alert>
                    )}

                    {/* Producto - Buscador */}
                    <div ref={dropdownRef} className="relative">
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                            Producto <span className="text-red-500">*</span>
                        </label>
                        {selectedProduct && !productId ? (
                            <div className="flex items-center justify-between px-3 py-2 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
                                <div>
                                    <span className="text-sm font-medium text-gray-900 dark:text-white">{selectedProduct.name}</span>
                                    {selectedProduct.sku && (
                                        <span className="ml-2 text-xs text-gray-500 dark:text-gray-400">SKU: {selectedProduct.sku}</span>
                                    )}
                                </div>
                                <button
                                    type="button"
                                    onClick={() => setSelectedProduct(null)}
                                    className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 ml-2"
                                >
                                    &times;
                                </button>
                            </div>
                        ) : productId ? (
                            <div className="px-3 py-2 bg-gray-100 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg text-sm text-gray-700 dark:text-gray-300">
                                {productId}
                            </div>
                        ) : (
                            <>
                                <Input
                                    type="text"
                                    value={searchTerm}
                                    onChange={(e) => {
                                        setSearchTerm(e.target.value);
                                        setShowDropdown(true);
                                    }}
                                    onFocus={() => { if (searchTerm.length >= 2) setShowDropdown(true); }}
                                    placeholder="Buscar por nombre o SKU..."
                                />
                                {showDropdown && (searchTerm.length >= 2) && (
                                    <div className="absolute z-20 w-full mt-1 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-xl max-h-48 overflow-auto">
                                        {searchLoading ? (
                                            <div className="px-4 py-3 text-sm text-gray-500 dark:text-gray-400 text-center">Buscando...</div>
                                        ) : searchResults.length > 0 ? (
                                            <ul>
                                                {searchResults.map((p) => (
                                                    <li key={p.id}>
                                                        <button
                                                            type="button"
                                                            onClick={() => handleSelectProduct(p)}
                                                            className="w-full px-4 py-2.5 text-left hover:bg-blue-50 dark:hover:bg-blue-900/20 flex items-center justify-between"
                                                        >
                                                            <span className="text-sm font-medium text-gray-900 dark:text-white truncate">{p.name}</span>
                                                            <span className="text-xs text-gray-500 dark:text-gray-400 ml-2 shrink-0">{p.sku}</span>
                                                        </button>
                                                    </li>
                                                ))}
                                            </ul>
                                        ) : (
                                            <div className="px-4 py-3 text-sm text-gray-500 dark:text-gray-400 text-center">No se encontraron productos</div>
                                        )}
                                    </div>
                                )}
                            </>
                        )}
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                            Cantidad <span className="text-red-500">*</span>
                        </label>
                        <Input
                            type="number"
                            value={quantity.toString()}
                            onChange={(e) => setQuantity(parseInt(e.target.value) || 0)}
                            placeholder="Positivo para agregar, negativo para quitar"
                            required
                        />
                        <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                            Positivo para agregar stock, negativo para quitar.
                        </p>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                            Razon <span className="text-red-500">*</span>
                        </label>
                        <Input
                            type="text"
                            value={reason}
                            onChange={(e) => setReason(e.target.value)}
                            placeholder="Conteo fisico, correccion, etc."
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Notas</label>
                        <textarea
                            value={notes}
                            onChange={(e) => setNotes(e.target.value)}
                            placeholder="Notas adicionales..."
                            rows={3}
                            className="w-full px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-transparent resize-none"
                        />
                    </div>

                    <div className="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
                        <Button type="button" variant="outline" onClick={onClose} disabled={loading}>
                            Cancelar
                        </Button>
                        <Button type="submit" variant="primary" disabled={loading || !selectedProduct}>
                            {loading ? 'Ajustando...' : 'Ajustar stock'}
                        </Button>
                    </div>
                </form>
            </div>
        </div>
    );
}

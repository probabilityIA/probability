'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import { TransferStockDTO } from '../../domain/types';
import { transferStockAction } from '../../infra/actions';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { getProductsAction } from '@/services/modules/products/infra/actions';
import { Warehouse } from '@/services/modules/warehouses/domain/types';
import { Button, Alert, Input } from '@/shared/ui';
import { getActionError } from '@/shared/utils/action-result';

interface ProductOption {
    id: string;
    name: string;
    sku: string;
}

interface TransferStockModalProps {
    fromWarehouseId: number;
    businessId?: number;
    productId?: string;
    onSuccess: () => void;
    onClose: () => void;
}

export default function TransferStockModal({ fromWarehouseId, businessId, productId, onSuccess, onClose }: TransferStockModalProps) {
    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
    const [selectedProduct, setSelectedProduct] = useState<ProductOption | null>(null);
    const [searchTerm, setSearchTerm] = useState('');
    const [searchResults, setSearchResults] = useState<ProductOption[]>([]);
    const [showDropdown, setShowDropdown] = useState(false);
    const [searchLoading, setSearchLoading] = useState(false);
    const dropdownRef = useRef<HTMLDivElement>(null);

    const [toWarehouseId, setToWarehouseId] = useState(0);
    const [quantity, setQuantity] = useState(1);
    const [reason, setReason] = useState('');
    const [notes, setNotes] = useState('');

    const [loading, setLoading] = useState(false);
    const [loadingWarehouses, setLoadingWarehouses] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    // Cargar bodegas
    useEffect(() => {
        const loadWarehouses = async () => {
            try {
                const response = await getWarehousesAction({
                    page: 1,
                    page_size: 100,
                    is_active: true,
                    business_id: businessId,
                });
                setWarehouses((response.data || []).filter(w => w.id !== fromWarehouseId));
            } catch {
                setError('Error al cargar bodegas');
            } finally {
                setLoadingWarehouses(false);
            }
        };
        loadWarehouses();
    }, [businessId, fromWarehouseId]);

    // Si viene productId pre-seleccionado
    useEffect(() => {
        if (productId) {
            setSelectedProduct({ id: productId, name: '', sku: '' });
        }
    }, [productId]);

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

        try {
            const dto: TransferStockDTO = {
                product_id: selectedProduct.id,
                from_warehouse_id: fromWarehouseId,
                to_warehouse_id: toWarehouseId,
                quantity,
                reason: reason.trim() || undefined,
                notes: notes.trim() || undefined,
            };
            await transferStockAction(dto, businessId);
            setSuccess('Transferencia realizada exitosamente');
            setTimeout(() => onSuccess(), 800);
        } catch (err: any) {
            setError(getActionError(err, 'Error al transferir stock'));
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-md max-h-[90vh] overflow-y-auto">
                <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700">
                    <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Transferir stock</h2>
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
                            Bodega destino <span className="text-red-500">*</span>
                        </label>
                        {loadingWarehouses ? (
                            <p className="text-sm text-gray-500 dark:text-gray-400">Cargando bodegas...</p>
                        ) : (
                            <select
                                value={toWarehouseId || ''}
                                onChange={(e) => setToWarehouseId(Number(e.target.value))}
                                className="w-full px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-transparent"
                                required
                            >
                                <option value="">-- Selecciona bodega destino --</option>
                                {warehouses.map(w => (
                                    <option key={w.id} value={w.id}>{w.name} ({w.code})</option>
                                ))}
                            </select>
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
                            min="1"
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Razon</label>
                        <Input
                            type="text"
                            value={reason}
                            onChange={(e) => setReason(e.target.value)}
                            placeholder="Reabastecimiento, redistribucion, etc."
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
                            {loading ? 'Transfiriendo...' : 'Transferir'}
                        </Button>
                    </div>
                </form>
            </div>
        </div>
    );
}

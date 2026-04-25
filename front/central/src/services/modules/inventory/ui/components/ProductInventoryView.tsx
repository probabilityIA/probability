'use client';

import { useState, useEffect, useCallback } from 'react';
import { getProductInventoryAction } from '../../infra/actions';
import { getProductsAction } from '@/services/modules/products/infra/actions';
import { InventoryLevel } from '../../domain/types';
import { Spinner } from '@/shared/ui';
import { AdjustmentsHorizontalIcon, XMarkIcon, ChevronRightIcon } from '@heroicons/react/24/outline';

interface ProductRow {
    id: string;
    name: string;
    sku: string;
    family?: string;
    familyId?: number;
    variantLabel?: string;
    variantAttributes?: any;
}

interface WarehouseSummary {
    warehouse_id: number;
    warehouse_name: string;
    warehouse_code: string;
    quantity: number;
    reserved_qty: number;
    available_qty: number;
    min_stock: number | null;
    max_stock: number | null;
    locations: number;
    product_id: string;
}

function aggregateByWarehouse(levels: InventoryLevel[]): WarehouseSummary[] {
    const map = new Map<number, WarehouseSummary>();
    for (const l of levels) {
        const existing = map.get(l.warehouse_id);
        if (existing) {
            existing.quantity += l.quantity;
            existing.reserved_qty += l.reserved_qty;
            existing.available_qty += l.available_qty;
            existing.locations += 1;
            if (l.min_stock != null) existing.min_stock = Math.min(existing.min_stock ?? l.min_stock, l.min_stock);
            if (l.max_stock != null) existing.max_stock = (existing.max_stock ?? 0) + l.max_stock;
        } else {
            map.set(l.warehouse_id, {
                warehouse_id: l.warehouse_id,
                warehouse_name: l.warehouse_name || String(l.warehouse_id),
                warehouse_code: l.warehouse_code || '',
                quantity: l.quantity,
                reserved_qty: l.reserved_qty,
                available_qty: l.available_qty,
                min_stock: l.min_stock ?? null,
                max_stock: l.max_stock ?? null,
                locations: 1,
                product_id: l.product_id,
            });
        }
    }
    return Array.from(map.values()).sort((a, b) => b.quantity - a.quantity);
}

interface ProductInventoryViewProps {
    businessId?: number;
    onAdjust?: (productId: string, warehouseId: number) => void;
    onRefreshRef?: (ref: () => void) => void;
}

const MODAL_PAGE_SIZE = 10;

export default function ProductInventoryView({ businessId, onAdjust, onRefreshRef }: ProductInventoryViewProps) {
    const [products, setProducts] = useState<ProductRow[]>([]);
    const [loadingProducts, setLoadingProducts] = useState(false);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [total, setTotal] = useState(0);
    const pageSize = 20;

    const [nameInput, setNameInput] = useState('');
    const [skuInput, setSkuInput] = useState('');
    const [nameFilter, setNameFilter] = useState('');
    const [skuFilter, setSkuFilter] = useState('');

    const [stockCounts, setStockCounts] = useState<Record<string, number | null>>({});

    const [selectedProduct, setSelectedProduct] = useState<ProductRow | null>(null);
    const [warehouses, setWarehouses] = useState<WarehouseSummary[]>([]);
    const [loadingLevels, setLoadingLevels] = useState(false);
    const [modalPage, setModalPage] = useState(1);

    const fetchProducts = useCallback(async () => {
        setLoadingProducts(true);
        try {
            const params: Record<string, any> = { page, page_size: pageSize };
            if (businessId) params.business_id = businessId;
            if (nameFilter) params.name = nameFilter;
            if (skuFilter) params.sku = skuFilter;
            const response = await getProductsAction(params);
            if (response.success && response.data) {
                const rows = response.data.map((p) => ({
                    id: p.id,
                    name: p.name,
                    sku: p.sku,
                    family: (p as any).family?.name,
                    familyId: p.family_id,
                    variantLabel: p.variant_label,
                    variantAttributes: p.variant_attributes
                }));
                setProducts(rows);
                setTotal((response as any).total ?? rows.length);
                setTotalPages((response as any).total_pages ?? 1);
                setStockCounts({});
                const entries = await Promise.all(
                    rows.map(async (p) => {
                        try {
                            const data = await getProductInventoryAction(p.id, businessId);
                            const levels = Array.isArray(data) ? data : [];
                            const uniqueWarehouses = new Set(levels.filter((l) => l.quantity > 0).map((l) => l.warehouse_id)).size;
                            return [p.id, uniqueWarehouses] as const;
                        } catch {
                            return [p.id, 0] as const;
                        }
                    })
                );
                setStockCounts(Object.fromEntries(entries));
            }
        } catch {
            setProducts([]);
        } finally {
            setLoadingProducts(false);
        }
    }, [businessId, page, pageSize, nameFilter, skuFilter]);

    useEffect(() => { fetchProducts(); }, [fetchProducts]);

    useEffect(() => {
        onRefreshRef?.(() => { fetchProducts(); });
    }, [fetchProducts, onRefreshRef]);

    const openModal = async (p: ProductRow) => {
        setSelectedProduct(p);
        setWarehouses([]);
        setModalPage(1);
        setLoadingLevels(true);
        try {
            const data = await getProductInventoryAction(p.id, businessId);
            const levels = Array.isArray(data) ? data : [];
            setWarehouses(aggregateByWarehouse(levels));
        } catch {
            setWarehouses([]);
        } finally {
            setLoadingLevels(false);
        }
    };

    const closeModal = () => {
        setSelectedProduct(null);
        setWarehouses([]);
    };

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        setNameFilter(nameInput);
        setSkuFilter(skuInput);
        setPage(1);
    };

    const handleClear = () => {
        setNameInput(''); setSkuInput('');
        setNameFilter(''); setSkuFilter('');
        setPage(1);
    };

    const isLowStock = (w: WarehouseSummary) => {
        if (w.min_stock != null) return w.available_qty <= w.min_stock;
        return false;
    };

    const getVariantDisplay = (p: ProductRow) => {
        if (!p.familyId) {
            return <span className="text-gray-300 dark:text-gray-600">&mdash;</span>;
        }

        if (p.variantLabel) {
            return <span className="text-sm text-gray-700 dark:text-gray-300">{p.variantLabel}</span>;
        }

        if (p.variantAttributes && typeof p.variantAttributes === 'object') {
            const attrs = Object.values(p.variantAttributes)
                .filter(Boolean)
                .map(String)
                .join(' - ');
            return <span className="text-sm text-gray-700 dark:text-gray-300">{attrs || '-'}</span>;
        }

        return <span className="text-sm text-gray-700 dark:text-gray-300">Variante</span>;
    };

    const renderStockBadge = (productId: string) => {
        const count = stockCounts[productId];
        if (count === undefined) {
            return <span className="inline-block w-4 h-4 rounded-full border-2 border-gray-300 border-t-gray-500 animate-spin" />;
        }
        if (count === 0) {
            return <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400">Sin stock</span>;
        }
        return <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 dark:bg-green-900 text-green-700 dark:text-green-300">{count} {count === 1 ? 'bodega' : 'bodegas'}</span>;
    };

    const modalTotalPages = Math.max(1, Math.ceil(warehouses.length / MODAL_PAGE_SIZE));
    const pagedWarehouses = warehouses.slice((modalPage - 1) * MODAL_PAGE_SIZE, modalPage * MODAL_PAGE_SIZE);

    return (
        <>
            <div className="space-y-4">
                <form onSubmit={handleSearch} className="flex gap-2 flex-wrap items-end">
                    <div className="flex-1 min-w-[160px]">
                        <label className="block text-xs text-gray-500 dark:text-gray-400 mb-1">Nombre</label>
                        <input
                            type="text"
                            value={nameInput}
                            onChange={(e) => { setNameInput(e.target.value); setSkuInput(''); }}
                            placeholder="Buscar por nombre..."
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                        />
                    </div>
                    <div className="flex-1 min-w-[160px]">
                        <label className="block text-xs text-gray-500 dark:text-gray-400 mb-1">SKU</label>
                        <input
                            type="text"
                            value={skuInput}
                            onChange={(e) => { setSkuInput(e.target.value); setNameInput(''); }}
                            placeholder="Buscar por SKU..."
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                        />
                    </div>
                    <button type="submit" className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-700 transition-colors">
                        Buscar
                    </button>
                    {(nameFilter || skuFilter) && (
                        <button type="button" onClick={handleClear} className="px-4 py-2 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-200 rounded-lg text-sm hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors">
                            Limpiar
                        </button>
                    )}
                </form>

                <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
                    <table className="table w-full">
                        <thead>
                            <tr>
                                <th className="text-left">Producto</th>
                                <th className="text-left">Variante</th>
                                <th className="text-left">Familia</th>
                                <th className="text-left">SKU</th>
                                <th className="text-center">Stock</th>
                                <th className="text-center">Acciones</th>
                            </tr>
                        </thead>
                        <tbody>
                            {loadingProducts ? (
                                <tr>
                                    <td colSpan={6} className="px-6 py-12 text-center text-gray-500 dark:text-gray-400">
                                        <div className="flex justify-center items-center gap-3">
                                            <div className="spinner"></div>
                                            <span>Cargando...</span>
                                        </div>
                                    </td>
                                </tr>
                            ) : products.length === 0 ? (
                                <tr>
                                    <td colSpan={6} className="px-6 py-12 text-center text-gray-500 dark:text-gray-400">
                                        No se encontraron productos
                                    </td>
                                </tr>
                            ) : (
                                products.map((p) => {
                                    const hasStock = (stockCounts[p.id] ?? 0) > 0;
                                    return (
                                        <tr key={p.id} className="bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
                                            <td className="text-left font-medium text-gray-900 dark:text-white">{p.name}</td>
                                            <td className="text-left text-sm text-gray-500 dark:text-gray-400">{getVariantDisplay(p)}</td>
                                            <td className="text-left text-sm text-gray-500 dark:text-gray-400">{p.family ?? <span className="text-gray-300 dark:text-gray-600">&mdash;</span>}</td>
                                            <td className="text-left text-sm text-gray-500 dark:text-gray-400 font-mono">{p.sku}</td>
                                            <td className="text-center">{renderStockBadge(p.id)}</td>
                                            <td className="text-center">
                                                {hasStock && (
                                                    <button
                                                        onClick={() => openModal(p)}
                                                        className="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 dark:hover:bg-blue-900/30 rounded-md transition-colors"
                                                        title="Ver detalle de bodegas"
                                                    >
                                                        <ChevronRightIcon className="w-4 h-4" />
                                                    </button>
                                                )}
                                            </td>
                                        </tr>
                                    );
                                })
                            )}
                        </tbody>
                    </table>

                    {!loadingProducts && total > 0 && (
                        <div className="pagination-alt border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-700">
                            <div className="flex items-center justify-center gap-3 w-full flex-wrap py-1">
                                <button onClick={() => setPage((p) => Math.max(1, p - 1))} disabled={page === 1 || loadingProducts} className="pagination-button">
                                    &larr; Anterior
                                </button>
                                <span className="pagination-info">
                                    Pagina {page} de {totalPages} ({total} registros totales)
                                </span>
                                <button onClick={() => setPage((p) => Math.min(totalPages, p + 1))} disabled={page === totalPages || loadingProducts} className="pagination-button">
                                    Siguiente &rarr;
                                </button>
                            </div>
                        </div>
                    )}
                </div>
            </div>

            {selectedProduct && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
                    <div className="absolute inset-0 bg-black/50" onClick={closeModal} />
                    <div className="relative bg-white dark:bg-gray-800 rounded-xl shadow-2xl w-full max-w-4xl max-h-[85vh] flex flex-col">
                        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex-shrink-0">
                            <div>
                                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">{selectedProduct.name}</h2>
                                <p className="text-sm text-gray-500 dark:text-gray-400 font-mono mt-0.5">{selectedProduct.sku}</p>
                            </div>
                            <button
                                onClick={closeModal}
                                className="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                            >
                                <XMarkIcon className="w-5 h-5" />
                            </button>
                        </div>

                        <div className="overflow-auto flex-1">
                            {loadingLevels ? (
                                <div className="flex justify-center items-center py-16">
                                    <Spinner size="lg" />
                                </div>
                            ) : warehouses.length === 0 ? (
                                <p className="text-sm text-gray-500 dark:text-gray-400 text-center py-16">
                                    Este producto no tiene stock en ninguna bodega.
                                </p>
                            ) : (
                                <table className="table w-full">
                                    <thead>
                                        <tr>
                                            <th className="text-left">Bodega</th>
                                            <th className="text-center">Total</th>
                                            <th className="text-center">Reservado</th>
                                            <th className="text-center">Disponible</th>
                                            <th className="text-center">Min / Max</th>
                                            <th className="text-center">Estado</th>
                                            {onAdjust && <th className="text-center">Acciones</th>}
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {pagedWarehouses.map((w) => (
                                            <tr key={w.warehouse_id} className="border-t border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                                                <td>
                                                    <span className="text-sm font-medium text-gray-900 dark:text-white">{w.warehouse_name}</span>
                                                    <div className="flex items-center gap-2 mt-0.5">
                                                        {w.warehouse_code && <span className="text-xs text-gray-500 dark:text-gray-400 font-mono">{w.warehouse_code}</span>}
                                                        {w.locations > 1 && (
                                                            <span className="text-xs text-gray-400 dark:text-gray-500">{w.locations} ubicaciones</span>
                                                        )}
                                                    </div>
                                                </td>
                                                <td className="text-center text-sm font-medium text-gray-900 dark:text-white">{w.quantity}</td>
                                                <td className="text-center text-sm">
                                                    <span className={w.reserved_qty > 0 ? 'text-orange-600 font-medium' : 'text-gray-400 dark:text-gray-500'}>{w.reserved_qty}</span>
                                                </td>
                                                <td className="text-center text-sm font-semibold text-gray-900 dark:text-white">{w.available_qty}</td>
                                                <td className="text-center text-xs text-gray-500 dark:text-gray-400">
                                                    {w.min_stock != null ? w.min_stock : '—'} / {w.max_stock != null ? w.max_stock : '—'}
                                                </td>
                                                <td className="text-center">
                                                    {isLowStock(w) ? (
                                                        <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200">Stock bajo</span>
                                                    ) : (
                                                        <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200">OK</span>
                                                    )}
                                                </td>
                                                {onAdjust && (
                                                    <td className="text-center">
                                                        <button
                                                            onClick={() => { onAdjust(w.product_id, w.warehouse_id); closeModal(); }}
                                                            className="p-1.5 text-gray-500 hover:text-blue-600 hover:bg-blue-50 dark:hover:bg-blue-900/30 rounded-md transition-colors"
                                                            title="Ajustar stock"
                                                        >
                                                            <AdjustmentsHorizontalIcon className="w-4 h-4" />
                                                        </button>
                                                    </td>
                                                )}
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            )}
                        </div>

                        {!loadingLevels && warehouses.length > MODAL_PAGE_SIZE && (
                            <div className="flex items-center justify-center gap-3 px-6 py-3 border-t border-gray-200 dark:border-gray-700 flex-shrink-0">
                                <button
                                    onClick={() => setModalPage((p) => Math.max(1, p - 1))}
                                    disabled={modalPage === 1}
                                    className="pagination-button"
                                >
                                    &larr; Anterior
                                </button>
                                <span className="pagination-info">
                                    Pagina {modalPage} de {modalTotalPages} ({warehouses.length} bodegas)
                                </span>
                                <button
                                    onClick={() => setModalPage((p) => Math.min(modalTotalPages, p + 1))}
                                    disabled={modalPage === modalTotalPages}
                                    className="pagination-button"
                                >
                                    Siguiente &rarr;
                                </button>
                            </div>
                        )}
                    </div>
                </div>
            )}
        </>
    );
}

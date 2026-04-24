'use client';

import { useState, useEffect, useCallback } from 'react';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { Warehouse } from '@/services/modules/warehouses/domain/types';
import { Spinner } from '@/shared/ui';
import { XMarkIcon, ChevronRightIcon } from '@heroicons/react/24/outline';
import InventoryLevelList from './InventoryLevelList';

interface WarehouseInventoryViewProps {
    businessId?: number;
    onAdjust?: (productId: string, warehouseId: number) => void;
    onRefreshRef?: (ref: () => void) => void;
}

export default function WarehouseInventoryView({ businessId, onAdjust, onRefreshRef }: WarehouseInventoryViewProps) {
    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
    const [loading, setLoading] = useState(true);
    const [selectedWarehouse, setSelectedWarehouse] = useState<Warehouse | null>(null);
    const [refreshInner, setRefreshInner] = useState<(() => void) | null>(null);

    const fetchWarehouses = useCallback(async () => {
        setLoading(true);
        try {
            const res = await getWarehousesAction({ page: 1, page_size: 100, is_active: true, business_id: businessId });
            setWarehouses(res.data || []);
        } catch {
            setWarehouses([]);
        } finally {
            setLoading(false);
        }
    }, [businessId]);

    useEffect(() => { fetchWarehouses(); }, [fetchWarehouses]);

    useEffect(() => {
        onRefreshRef?.(() => { fetchWarehouses(); refreshInner?.(); });
    }, [fetchWarehouses, onRefreshRef, refreshInner]);

    const handleInnerRefreshRef = useCallback((ref: () => void) => {
        setRefreshInner(() => ref);
    }, []);

    const handleAdjust = selectedWarehouse && onAdjust
        ? (productId: string) => { onAdjust(productId, selectedWarehouse.id); setSelectedWarehouse(null); }
        : undefined;

    return (
        <>
            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
                <table className="table w-full">
                    <thead>
                        <tr>
                            <th className="text-left">Bodega</th>
                            <th className="text-left">Codigo</th>
                            <th className="text-left">Ciudad</th>
                            <th className="text-center w-12"></th>
                        </tr>
                    </thead>
                    <tbody>
                        {loading ? (
                            <tr>
                                <td colSpan={4} className="px-6 py-12 text-center text-gray-500 dark:text-gray-400">
                                    <div className="flex justify-center items-center gap-3">
                                        <div className="spinner"></div>
                                        <span>Cargando...</span>
                                    </div>
                                </td>
                            </tr>
                        ) : warehouses.length === 0 ? (
                            <tr>
                                <td colSpan={4} className="px-6 py-12 text-center text-gray-500 dark:text-gray-400">
                                    No hay bodegas activas. Crea una en el modulo de Bodegas.
                                </td>
                            </tr>
                        ) : (
                            warehouses.map((w) => (
                                <tr key={w.id} className="bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
                                    <td>
                                        <span className="font-medium text-gray-900 dark:text-white">{w.name}</span>
                                        {w.is_default && (
                                            <span className="ml-2 inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300">Principal</span>
                                        )}
                                    </td>
                                    <td className="text-sm text-gray-500 dark:text-gray-400 font-mono">{w.code}</td>
                                    <td className="text-sm text-gray-500 dark:text-gray-400">{w.city || <span className="text-gray-300 dark:text-gray-600">&mdash;</span>}</td>
                                    <td className="text-center">
                                        <button
                                            onClick={() => setSelectedWarehouse(w)}
                                            className="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 dark:hover:bg-blue-900/30 rounded-md transition-colors"
                                            title="Ver inventario"
                                        >
                                            <ChevronRightIcon className="w-4 h-4" />
                                        </button>
                                    </td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>

            {selectedWarehouse && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
                    <div className="absolute inset-0 bg-black/50" onClick={() => setSelectedWarehouse(null)} />
                    <div className="relative bg-white dark:bg-gray-800 rounded-xl shadow-2xl w-full max-w-5xl max-h-[90vh] flex flex-col">
                        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex-shrink-0">
                            <div>
                                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">{selectedWarehouse.name}</h2>
                                <p className="text-sm text-gray-500 dark:text-gray-400 font-mono mt-0.5">{selectedWarehouse.code}{selectedWarehouse.city ? ` — ${selectedWarehouse.city}` : ''}</p>
                            </div>
                            <button
                                onClick={() => setSelectedWarehouse(null)}
                                className="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                            >
                                <XMarkIcon className="w-5 h-5" />
                            </button>
                        </div>
                        <div className="overflow-auto flex-1 p-4">
                            <InventoryLevelList
                                warehouseId={selectedWarehouse.id}
                                selectedBusinessId={businessId}
                                onRefreshRef={handleInnerRefreshRef}
                                onAdjust={handleAdjust}
                            />
                        </div>
                    </div>
                </div>
            )}
        </>
    );
}

'use client';

import { useState, useEffect, useCallback } from 'react';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { Warehouse } from '@/services/modules/warehouses/domain/types';
import { Spinner } from '@/shared/ui';
import { getWarehouseInventoryAction } from '../../infra/actions';

interface WarehouseInventoryViewProps {
    businessId?: number;
}

interface WarehouseStats {
    products: number;
    totalQty: number;
    loading: boolean;
}

export default function WarehouseInventoryView({ businessId }: WarehouseInventoryViewProps) {
    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
    const [loading, setLoading] = useState(true);
    const [stats, setStats] = useState<Record<number, WarehouseStats>>({});

    const fetchWarehouses = useCallback(async () => {
        setLoading(true);
        try {
            const res = await getWarehousesAction({ page: 1, page_size: 100, is_active: true, business_id: businessId });
            const list: Warehouse[] = res.data || [];
            setWarehouses(list);
            setStats(Object.fromEntries(list.map((w) => [w.id, { products: 0, totalQty: 0, loading: true }])));

            await Promise.all(list.map(async (w) => {
                try {
                    const inv = await getWarehouseInventoryAction(w.id, { page: 1, page_size: 500, business_id: businessId });
                    const levels = inv.data ?? [];
                    const uniqueProducts = new Set(levels.map((l) => l.product_id)).size;
                    const totalQty = levels.reduce((s, l) => s + l.quantity, 0);
                    setStats((prev) => ({ ...prev, [w.id]: { products: uniqueProducts, totalQty, loading: false } }));
                } catch {
                    setStats((prev) => ({ ...prev, [w.id]: { products: 0, totalQty: 0, loading: false } }));
                }
            }));
        } catch {
            setWarehouses([]);
        } finally {
            setLoading(false);
        }
    }, [businessId]);

    useEffect(() => { fetchWarehouses(); }, [fetchWarehouses]);

    return (
        <>
            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
                <table className="table w-full">
                    <thead>
                        <tr>
                            <th className="text-left">Bodega</th>
                            <th className="text-left">Codigo</th>
                            <th className="text-left">Ciudad</th>
                            <th className="text-center">Productos</th>
                            <th className="text-center">Inventario total</th>
                        </tr>
                    </thead>
                    <tbody>
                        {loading ? (
                            <tr>
                                <td colSpan={6} className="px-6 py-12 text-center text-gray-500 dark:text-gray-400">
                                    <div className="flex justify-center items-center gap-3">
                                        <div className="spinner"></div>
                                        <span>Cargando...</span>
                                    </div>
                                </td>
                            </tr>
                        ) : warehouses.length === 0 ? (
                            <tr>
                                <td colSpan={6} className="px-6 py-12 text-center text-gray-500 dark:text-gray-400">
                                    No hay bodegas activas. Crea una en el modulo de Bodegas.
                                </td>
                            </tr>
                        ) : (
                            warehouses.map((w) => {
                                const s = stats[w.id];
                                return (
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
                                            {s?.loading ? (
                                                <span className="inline-block w-3 h-3 rounded-full border-2 border-gray-300 border-t-gray-500 animate-spin" />
                                            ) : (
                                                <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-indigo-100 dark:bg-indigo-900 text-indigo-700 dark:text-indigo-300">
                                                    {s?.products ?? 0}
                                                </span>
                                            )}
                                        </td>
                                        <td className="text-center">
                                            {s?.loading ? (
                                                <span className="inline-block w-3 h-3 rounded-full border-2 border-gray-300 border-t-gray-500 animate-spin" />
                                            ) : (
                                                <span className="text-sm font-semibold text-gray-900 dark:text-white">
                                                    {s?.totalQty ?? 0}
                                                </span>
                                            )}
                                        </td>
                                    </tr>
                                );
                            })
                        )}
                    </tbody>
                </table>
            </div>
        </>
    );
}

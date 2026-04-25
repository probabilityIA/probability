'use client';

import { useState, useEffect, useCallback } from 'react';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { Warehouse } from '@/services/modules/warehouses/domain/types';
import { getMovementsAction } from '../../infra/actions';
import { ChevronRightIcon, XMarkIcon } from '@heroicons/react/24/outline';
import { Spinner } from '@/shared/ui';
import MovementsInlineTable from './MovementsInlineTable';

interface Props {
    businessId?: number;
}

interface WarehouseWithCount extends Warehouse {
    movCount: number;
}

export default function MovementsByWarehouseView({ businessId }: Props) {
    const [warehouses, setWarehouses] = useState<WarehouseWithCount[]>([]);
    const [loading, setLoading] = useState(true);
    const [selected, setSelected] = useState<WarehouseWithCount | null>(null);

    const load = useCallback(async () => {
        setLoading(true);
        try {
            const res = await getWarehousesAction({ page: 1, page_size: 100, is_active: true, business_id: businessId });
            const all: Warehouse[] = res.data ?? [];

            const counts = await Promise.all(
                all.map((w) =>
                    getMovementsAction({ warehouse_id: w.id, page: 1, page_size: 1, business_id: businessId })
                        .then((r) => r.total ?? 0)
                        .catch(() => 0)
                )
            );

            const withMovements: WarehouseWithCount[] = all
                .map((w, i) => ({ ...w, movCount: counts[i] }))
                .filter((w) => w.movCount > 0)
                .sort((a, b) => b.movCount - a.movCount);

            setWarehouses(withMovements);
        } finally {
            setLoading(false);
        }
    }, [businessId]);

    useEffect(() => { load(); }, [load]);

    return (
        <>
            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
                <table className="table w-full">
                    <thead>
                        <tr>
                            <th className="text-left">Bodega</th>
                            <th className="text-left">Codigo</th>
                            <th className="text-left">Ciudad</th>
                            <th className="text-center">Movimientos</th>
                            <th className="text-center w-12"></th>
                        </tr>
                    </thead>
                    <tbody>
                        {loading ? (
                            <tr><td colSpan={5} className="py-12 text-center"><div className="flex justify-center items-center gap-3"><div className="spinner"></div><span className="text-sm text-gray-500">Cargando...</span></div></td></tr>
                        ) : warehouses.length === 0 ? (
                            <tr><td colSpan={5} className="py-12 text-center text-sm text-gray-500">Sin movimientos en ninguna bodega.</td></tr>
                        ) : (
                            warehouses.map((w) => (
                                <tr key={w.id} className="hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
                                    <td>
                                        <span className="font-medium text-gray-900 dark:text-white">{w.name}</span>
                                        {w.is_default && <span className="ml-2 inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-700">Principal</span>}
                                    </td>
                                    <td className="text-sm text-gray-500 dark:text-gray-400 font-mono">{w.code}</td>
                                    <td className="text-sm text-gray-500 dark:text-gray-400">{w.city || <span className="text-gray-300">&mdash;</span>}</td>
                                    <td className="text-center">
                                        <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-700">{w.movCount}</span>
                                    </td>
                                    <td className="text-center">
                                        <button onClick={() => setSelected(w)} className="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors" title="Ver movimientos">
                                            <ChevronRightIcon className="w-4 h-4" />
                                        </button>
                                    </td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>

            {selected && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
                    <div className="absolute inset-0 bg-black/50" onClick={() => setSelected(null)} />
                    <div className="relative bg-white dark:bg-gray-800 rounded-xl shadow-2xl w-full max-w-5xl max-h-[90vh] flex flex-col">
                        <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex-shrink-0">
                            <div>
                                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">{selected.name}</h2>
                                <p className="text-sm text-gray-500 font-mono mt-0.5">{selected.code}{selected.city ? ` — ${selected.city}` : ''}</p>
                            </div>
                            <button onClick={() => setSelected(null)} className="p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 transition-colors"><XMarkIcon className="w-5 h-5" /></button>
                        </div>
                        <div className="overflow-auto flex-1 p-5">
                            <MovementsInlineTable warehouseId={selected.id} businessId={businessId} />
                        </div>
                    </div>
                </div>
            )}
        </>
    );
}

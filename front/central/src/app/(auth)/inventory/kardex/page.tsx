'use client';

import { useEffect, useState } from 'react';
import { DocumentArrowDownIcon } from '@heroicons/react/24/outline';
import { Alert, Table, Spinner, Button } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useInventoryBusiness } from '@/shared/contexts/inventory-business-context';
import { useKardex } from '@/services/modules/inventory/ui/hooks/useAudit';
import { KardexEntry } from '@/services/modules/inventory/domain/audit-types';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { Warehouse } from '@/services/modules/warehouses/domain/types';

export default function KardexPage() {
    const { isSuperAdmin } = usePermissions();
    const { selectedBusinessId } = useInventoryBusiness();
    const businessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;

    const { data, loading, error, load } = useKardex();
    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);

    const [form, setForm] = useState({ product_id: '', warehouse_id: 0, from: '', to: '' });

    const requiresBusiness = isSuperAdmin && selectedBusinessId === null;

    useEffect(() => {
        if (requiresBusiness) return;
        (async () => {
            try {
                const r = await getWarehousesAction({ page: 1, page_size: 100, is_active: true, business_id: businessId });
                setWarehouses(r.data || []);
            } catch { setWarehouses([]); }
        })();
    }, [businessId, requiresBusiness]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!form.product_id || !form.warehouse_id) return;
        await load({
            product_id: form.product_id,
            warehouse_id: Number(form.warehouse_id),
            from: form.from ? `${form.from}T00:00:00Z` : undefined,
            to: form.to ? `${form.to}T23:59:59Z` : undefined,
            business_id: businessId,
        });
    };

    const columns = [
        { key: 'movement_id', label: '#', align: 'center' as const },
        { key: 'type', label: 'Tipo' },
        { key: 'qty', label: 'Cantidad', align: 'center' as const },
        { key: 'prev', label: 'Saldo previo', align: 'center' as const },
        { key: 'new', label: 'Nuevo saldo', align: 'center' as const },
        { key: 'running', label: 'Acumulado', align: 'center' as const },
        { key: 'ref', label: 'Referencia' },
    ];

    const renderRow = (e: KardexEntry) => ({
        movement_id: <span className="text-xs text-gray-500 dark:text-gray-400 font-mono">#{e.movement_id}</span>,
        type: (
            <div>
                <span className="text-sm font-medium text-gray-900 dark:text-white">{e.movement_type_name}</span>
                <span className="block text-xs text-gray-500 dark:text-gray-400 font-mono">{e.movement_type_code}</span>
            </div>
        ),
        qty: <span className={`text-sm font-semibold ${e.quantity > 0 ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}`}>{e.quantity > 0 ? `+${e.quantity}` : e.quantity}</span>,
        prev: <span className="text-sm text-gray-600 dark:text-gray-300">{e.previous_qty}</span>,
        new: <span className="text-sm text-gray-600 dark:text-gray-300">{e.new_qty}</span>,
        running: <span className="text-sm font-bold text-purple-700 dark:text-purple-300">{e.running_balance}</span>,
        ref: (
            <div className="text-xs text-gray-600 dark:text-gray-300">
                {e.reference_type && <span className="font-mono">{e.reference_type}:{e.reference_id}</span>}
                {e.reason && <p className="text-gray-500 dark:text-gray-400 truncate max-w-[200px]">{e.reason}</p>}
            </div>
        ),
    });

    const downloadCSV = () => {
        if (!data) return;
        const headers = ['movement_id', 'type_code', 'type_name', 'quantity', 'previous', 'new', 'running', 'reason'];
        const rows = data.entries.map((e) => [e.movement_id, e.movement_type_code, e.movement_type_name, e.quantity, e.previous_qty, e.new_qty, e.running_balance, e.reason].join(','));
        const csv = [headers.join(','), ...rows].join('\n');
        const blob = new Blob([csv], { type: 'text/csv' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `kardex-${data.product_id}-${data.warehouse_id}.csv`;
        a.click();
    };

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="space-y-4">                {requiresBusiness ? (
                    <Alert type="info">Selecciona un negocio.</Alert>
                ) : (
                    <>
                        <form onSubmit={handleSubmit} className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4">
                            <div className="grid grid-cols-1 md:grid-cols-5 gap-3 items-end">
                                <div>
                                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Producto (SKU) *</label>
                                    <input required value={form.product_id} onChange={(e) => setForm({ ...form, product_id: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                                </div>
                                <div>
                                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Bodega *</label>
                                    <select required value={form.warehouse_id} onChange={(e) => setForm({ ...form, warehouse_id: Number(e.target.value) })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                                        <option value={0}>—</option>
                                        {warehouses.map((w) => <option key={w.id} value={w.id}>{w.name}</option>)}
                                    </select>
                                </div>
                                <div>
                                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Desde</label>
                                    <input type="date" value={form.from} onChange={(e) => setForm({ ...form, from: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                                </div>
                                <div>
                                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Hasta</label>
                                    <input type="date" value={form.to} onChange={(e) => setForm({ ...form, to: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                                </div>
                                <Button type="submit" variant="primary" disabled={loading}>{loading ? 'Cargando...' : 'Consultar'}</Button>
                            </div>
                        </form>

                        {error && <Alert type="error">{error}</Alert>}

                        {data && (
                            <>
                                <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
                                    <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-4 py-3">
                                        <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide">Total entradas</p>
                                        <p className="text-2xl font-bold text-green-600 dark:text-green-400">+{data.total_in}</p>
                                    </div>
                                    <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-4 py-3">
                                        <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide">Total salidas</p>
                                        <p className="text-2xl font-bold text-red-600 dark:text-red-400">-{data.total_out}</p>
                                    </div>
                                    <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-4 py-3">
                                        <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide">Saldo final</p>
                                        <p className="text-2xl font-bold text-purple-700 dark:text-purple-300">{data.final_balance}</p>
                                    </div>
                                </div>

                                <div className="flex justify-end">
                                    <button onClick={downloadCSV} className="inline-flex items-center gap-2 px-4 py-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-200 text-sm font-medium rounded-md">
                                        <DocumentArrowDownIcon className="w-4 h-4" /> Exportar CSV
                                    </button>
                                </div>

                                <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                                    <Table
                                        columns={columns}
                                        data={data.entries.map(renderRow)}
                                        keyExtractor={(_, i) => String(data.entries[i]?.movement_id || i)}
                                        emptyMessage="Sin movimientos en el rango"
                                    />
                                </div>
                            </>
                        )}

                        {loading && !data && <div className="flex justify-center p-8"><Spinner size="lg" /></div>}
                    </>
                )}
            </div>
        </div>
    );
}

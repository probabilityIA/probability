'use client';

import { useEffect, useState } from 'react';
import { ChartBarIcon, BoltIcon } from '@heroicons/react/24/outline';
import { Alert, Table, Spinner, Button } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useInventoryBusiness } from '@/shared/contexts/inventory-business-context';
import { runSlottingAction, listVelocitiesAction } from '@/services/modules/inventory/infra/actions/operations';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { ProductVelocity } from '@/services/modules/inventory/domain/operations-types';
import { Warehouse } from '@/services/modules/warehouses/domain/types';

export default function SlottingPage() {
    const { isSuperAdmin } = usePermissions();
    const { selectedBusinessId } = useInventoryBusiness();
    const businessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;

    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
    const [warehouseId, setWarehouseId] = useState<number>(0);
    const [period, setPeriod] = useState('30d');
    const [rank, setRank] = useState('');
    const [velocities, setVelocities] = useState<ProductVelocity[]>([]);
    const [loading, setLoading] = useState(false);
    const [running, setRunning] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [runMessage, setRunMessage] = useState<string | null>(null);

    const requiresBusiness = isSuperAdmin && selectedBusinessId === null;

    useEffect(() => {
        if (requiresBusiness) return;
        (async () => {
            try {
                const r = await getWarehousesAction({ page: 1, page_size: 100, is_active: true, business_id: businessId });
                setWarehouses(r.data || []);
                if (r.data && r.data[0]) setWarehouseId(r.data[0].id);
            } catch { setWarehouses([]); }
        })();
    }, [businessId, requiresBusiness]);

    const load = async () => {
        if (!warehouseId) return;
        setLoading(true);
        setError(null);
        try {
            const r = await listVelocitiesAction({ warehouse_id: warehouseId, period, rank: rank || undefined, limit: 100, business_id: businessId });
            setVelocities(r.data || []);
        } catch (e: any) { setError(e.message); }
        finally { setLoading(false); }
    };

    useEffect(() => { if (warehouseId) load(); }, [warehouseId, period, rank]);

    const handleRun = async () => {
        if (!warehouseId) return;
        setRunning(true);
        setRunMessage(null);
        try {
            const r = await runSlottingAction({ warehouse_id: warehouseId, period }, businessId);
            if (r.success) setRunMessage(`Slotting computado · ${r.data?.total_scanned || 0} productos escaneados`);
            else setRunMessage(r.error || 'Error');
            load();
        } finally { setRunning(false); }
    };

    const rankStyles: Record<string, string> = {
        A: 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200',
        B: 'bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200',
        C: 'bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200',
    };

    const columns = [
        { key: 'product', label: 'Producto' },
        { key: 'rank', label: 'Rank ABC', align: 'center' as const },
        { key: 'units', label: 'Unidades movidas', align: 'center' as const },
        { key: 'period', label: 'Período', align: 'center' as const },
        { key: 'computed', label: 'Computado', align: 'center' as const },
    ];

    const renderRow = (v: ProductVelocity) => ({
        product: <span className="text-sm font-mono text-gray-900 dark:text-white">{v.product_id}</span>,
        rank: <span className={`inline-flex items-center justify-center w-8 h-8 rounded-full font-bold text-sm ${rankStyles[v.rank] || 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300'}`}>{v.rank}</span>,
        units: <span className="text-sm font-bold text-gray-900 dark:text-white">{v.units_moved}</span>,
        period: <span className="text-xs text-gray-600 dark:text-gray-300">{v.period}</span>,
        computed: <span className="text-xs text-gray-500">{v.computed_at ? new Date(v.computed_at).toLocaleString() : '—'}</span>,
    });

    const stats = {
        A: velocities.filter((v) => v.rank === 'A').length,
        B: velocities.filter((v) => v.rank === 'B').length,
        C: velocities.filter((v) => v.rank === 'C').length,
    };

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="space-y-4">
                <div className="flex items-center justify-between flex-wrap gap-3">                    {!requiresBusiness && (
                        <button onClick={handleRun} disabled={running || !warehouseId} className="px-5 py-3 btn-business-primary text-white font-bold rounded-lg shadow-lg hover:shadow-2xl transition-all duration-300 transform hover:scale-110 hover:-translate-y-1 flex items-center gap-2 disabled:opacity-60">
                            <BoltIcon className="w-4 h-4" /> {running ? 'Computando...' : 'Ejecutar slotting'}
                        </button>
                    )}
                </div>

                {runMessage && <Alert type="info">{runMessage}</Alert>}

                {requiresBusiness ? <Alert type="info">Selecciona un negocio.</Alert> : (
                    <>
                        <div className="flex gap-3 items-end flex-wrap">
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Bodega</label>
                                <select value={warehouseId} onChange={(e) => setWarehouseId(Number(e.target.value))} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                                    <option value={0}>—</option>
                                    {warehouses.map((w) => <option key={w.id} value={w.id}>{w.name}</option>)}
                                </select>
                            </div>
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Período</label>
                                <select value={period} onChange={(e) => setPeriod(e.target.value)} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                                    <option value="7d">7 días</option>
                                    <option value="30d">30 días</option>
                                    <option value="90d">90 días</option>
                                    <option value="180d">180 días</option>
                                    <option value="365d">365 días</option>
                                </select>
                            </div>
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Rank</label>
                                <select value={rank} onChange={(e) => setRank(e.target.value)} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                                    <option value="">Todos</option>
                                    <option value="A">A (alta rotación)</option>
                                    <option value="B">B (media)</option>
                                    <option value="C">C (baja)</option>
                                </select>
                            </div>
                        </div>

                        <div className="grid grid-cols-3 gap-3">
                            <div className="bg-white dark:bg-gray-800 border-l-4 border-green-500 rounded-lg px-4 py-3">
                                <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide">Rank A</p>
                                <p className="text-2xl font-bold text-green-600 dark:text-green-400">{stats.A}</p>
                                <p className="text-xs text-gray-400">Alta rotación (80% del volumen)</p>
                            </div>
                            <div className="bg-white dark:bg-gray-800 border-l-4 border-yellow-500 rounded-lg px-4 py-3">
                                <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide">Rank B</p>
                                <p className="text-2xl font-bold text-yellow-600 dark:text-yellow-400">{stats.B}</p>
                                <p className="text-xs text-gray-400">Media rotación (15%)</p>
                            </div>
                            <div className="bg-white dark:bg-gray-800 border-l-4 border-red-500 rounded-lg px-4 py-3">
                                <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide">Rank C</p>
                                <p className="text-2xl font-bold text-red-600 dark:text-red-400">{stats.C}</p>
                                <p className="text-xs text-gray-400">Baja rotación (5%)</p>
                            </div>
                        </div>

                        {error && <Alert type="error">{error}</Alert>}

                        {loading ? <div className="flex justify-center p-8"><Spinner size="lg" /></div> : (
                            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                                <Table columns={columns} data={velocities.map(renderRow)} keyExtractor={(_, i) => String(velocities[i]?.id || i)} emptyMessage="Sin datos. Ejecuta el slotting primero." loading={loading} />
                            </div>
                        )}
                    </>
                )}
            </div>
        </div>
    );
}

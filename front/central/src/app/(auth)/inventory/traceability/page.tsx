'use client';

import { useState } from 'react';
import { PlusIcon, PencilIcon, TrashIcon, ArrowsRightLeftIcon } from '@heroicons/react/24/outline';
import { format } from 'date-fns';
import { Alert, Table, Spinner, Button } from '@/shared/ui';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useInventoryBusiness } from '@/shared/contexts/inventory-business-context';
import { useLots } from '@/services/modules/inventory/ui/hooks/useLots';
import { useSerials } from '@/services/modules/inventory/ui/hooks/useSerials';
import { useUoMs, useUoMConverter } from '@/services/modules/inventory/ui/hooks/useUoMs';
import { deleteLotAction } from '@/services/modules/inventory/infra/actions/traceability';
import { InventoryLot, InventorySerial, UnitOfMeasure } from '@/services/modules/inventory/domain/traceability-types';
import LotFormModal from '@/services/modules/inventory/ui/components/LotFormModal';
import SerialFormModal from '@/services/modules/inventory/ui/components/SerialFormModal';

type Tab = 'lots' | 'serials' | 'uoms';

export default function TraceabilityPage() {
    const { isSuperAdmin } = usePermissions();
    const { selectedBusinessId } = useInventoryBusiness();
    const businessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;
    const requiresBusiness = isSuperAdmin && selectedBusinessId === null;

    const [tab, setTab] = useState<Tab>('lots');

    const lots = useLots({ business_id: businessId });
    const serials = useSerials({ business_id: businessId });
    const uoms = useUoMs();

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="space-y-4">                {requiresBusiness ? <Alert type="info">Selecciona un negocio.</Alert> : (
                    <>
                        <div className="flex gap-2 border-b border-gray-200 dark:border-gray-700">
                            <TabButton active={tab === 'lots'} onClick={() => setTab('lots')} label="Lotes" count={lots.total} />
                            <TabButton active={tab === 'serials'} onClick={() => setTab('serials')} label="Series" count={serials.total} />
                            <TabButton active={tab === 'uoms'} onClick={() => setTab('uoms')} label="Unidades de medida" count={uoms.uoms.length} />
                        </div>

                        {tab === 'lots' && <LotsSection hook={lots} businessId={businessId} />}
                        {tab === 'serials' && <SerialsSection hook={serials} businessId={businessId} />}
                        {tab === 'uoms' && <UoMsSection uomsHook={uoms} businessId={businessId} />}
                    </>
                )}
            </div>
        </div>
    );
}

function TabButton({ active, onClick, label, count }: { active: boolean; onClick: () => void; label: string; count: number }) {
    return (
        <button onClick={onClick} className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${active ? 'border-business-primary text-business-primary' : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'}`}>
            {label} <span className="ml-1 text-xs text-gray-400">({count})</span>
        </button>
    );
}

function LotsSection({ hook, businessId }: { hook: ReturnType<typeof useLots>; businessId?: number }) {
    const { lots, loading, error, page, pageSize, total, totalPages, filters, setPage, setPageSize, setFilters, refresh } = hook;
    const [modalOpen, setModalOpen] = useState(false);
    const [editing, setEditing] = useState<InventoryLot | null>(null);
    const [deleting, setDeleting] = useState<InventoryLot | null>(null);

    const handleDelete = async () => {
        if (!deleting) return;
        await deleteLotAction(deleting.id, businessId);
        setDeleting(null);
        refresh();
    };

    const formatDate = (iso: string | null) => (iso ? format(new Date(iso), 'dd/MM/yyyy') : '—');
    const daysUntil = (iso: string | null) => {
        if (!iso) return null;
        return Math.ceil((new Date(iso).getTime() - Date.now()) / (1000 * 60 * 60 * 24));
    };

    const statusStyles: Record<string, string> = {
        active: 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200',
        expired: 'bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200',
        recalled: 'bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200',
        blocked: 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300',
    };

    const columns = [
        { key: 'product', label: 'Producto' },
        { key: 'lot_code', label: 'Lote' },
        { key: 'manufacture', label: 'Fabricación', align: 'center' as const },
        { key: 'expiration', label: 'Vence', align: 'center' as const },
        { key: 'days', label: 'Días', align: 'center' as const },
        { key: 'status', label: 'Estado', align: 'center' as const },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];

    const renderRow = (l: InventoryLot) => {
        const d = daysUntil(l.expiration_date);
        const daysClass = d === null ? 'text-gray-400 dark:text-gray-500' : d < 0 ? 'text-red-600 dark:text-red-400 font-semibold' : d < 30 ? 'text-yellow-600 dark:text-yellow-400 font-semibold' : 'text-gray-700 dark:text-gray-200';
        return {
            product: <span className="text-sm text-gray-700 dark:text-gray-200 font-mono">{l.product_id}</span>,
            lot_code: <span className="font-medium text-gray-900 dark:text-white">{l.lot_code}</span>,
            manufacture: <span className="text-sm text-gray-600 dark:text-gray-300">{formatDate(l.manufacture_date)}</span>,
            expiration: <span className="text-sm text-gray-600 dark:text-gray-300">{formatDate(l.expiration_date)}</span>,
            days: <span className={`text-sm ${daysClass}`}>{d === null ? '—' : d < 0 ? `Vencido ${Math.abs(d)}d` : `${d}d`}</span>,
            status: (
                <div className="flex justify-center">
                    <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${statusStyles[l.status] || statusStyles.blocked}`}>{l.status}</span>
                </div>
            ),
            actions: (
                <div className="flex justify-end gap-2">
                    <button onClick={() => { setEditing(l); setModalOpen(true); }} className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md"><PencilIcon className="w-4 h-4" /></button>
                    <button onClick={() => setDeleting(l)} className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md"><TrashIcon className="w-4 h-4" /></button>
                </div>
            ),
        };
    };

    return (
        <div className="space-y-4">
            <div className="flex items-end justify-between flex-wrap gap-3">
                <div className="flex gap-3 items-end flex-wrap">
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Producto (SKU)</label>
                        <input type="text" value={filters.product_id || ''} onChange={(e) => setFilters({ product_id: e.target.value || undefined })} placeholder="Filtrar por SKU..." className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Estado</label>
                        <select value={filters.status || ''} onChange={(e) => setFilters({ status: e.target.value || undefined })} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white">
                            <option value="">Todos</option>
                            <option value="active">Activo</option>
                            <option value="expired">Vencido</option>
                            <option value="recalled">Retirado</option>
                            <option value="blocked">Bloqueado</option>
                        </select>
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Vence en (días)</label>
                        <input type="number" min={0} value={filters.expiring_in_days ?? ''} onChange={(e) => setFilters({ expiring_in_days: e.target.value ? Number(e.target.value) : undefined })} placeholder="Ej: 30" className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white w-32" />
                    </div>
                </div>
                <button onClick={() => { setEditing(null); setModalOpen(true); }} className="px-5 py-3 btn-business-primary text-white font-bold rounded-lg shadow-lg hover:shadow-2xl transition-all duration-300 transform hover:scale-105 flex items-center gap-2">
                    <PlusIcon className="w-4 h-4" /> Nuevo lote
                </button>
            </div>

            {error && <Alert type="error">{error}</Alert>}
            {loading && lots.length === 0 ? <div className="flex justify-center p-8"><Spinner size="lg" /></div> : (
                <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                    <Table columns={columns} data={lots.map(renderRow)} keyExtractor={(_, i) => String(lots[i]?.id || i)} emptyMessage="Sin lotes" loading={loading} pagination={{ currentPage: page, totalPages, totalItems: total, itemsPerPage: pageSize, onPageChange: setPage, onItemsPerPageChange: setPageSize }} />
                </div>
            )}

            {modalOpen && <LotFormModal businessId={businessId} lot={editing} onClose={() => { setModalOpen(false); setEditing(null); }} onSuccess={() => { setModalOpen(false); setEditing(null); refresh(); }} />}
            {deleting && <ConfirmModal isOpen={true} onClose={() => setDeleting(null)} onConfirm={handleDelete} title="Eliminar lote" message={`Se eliminará el lote "${deleting.lot_code}". Acción irreversible.`} confirmText="Eliminar" type="danger" />}
        </div>
    );
}

function SerialsSection({ hook, businessId }: { hook: ReturnType<typeof useSerials>; businessId?: number }) {
    const { serials, loading, error, page, pageSize, total, totalPages, filters, setPage, setPageSize, setFilters, refresh } = hook;
    const [modalOpen, setModalOpen] = useState(false);
    const [editing, setEditing] = useState<InventorySerial | null>(null);

    const columns = [
        { key: 'product', label: 'Producto' },
        { key: 'serial', label: 'Serie' },
        { key: 'lot', label: 'Lote', align: 'center' as const },
        { key: 'location', label: 'Ubicación', align: 'center' as const },
        { key: 'received', label: 'Recibido', align: 'center' as const },
        { key: 'sold', label: 'Vendido', align: 'center' as const },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];

    const renderRow = (s: InventorySerial) => ({
        product: <span className="text-sm text-gray-700 dark:text-gray-200 font-mono">{s.product_id}</span>,
        serial: <span className="font-medium text-gray-900 dark:text-white">{s.serial_number}</span>,
        lot: <span className="text-sm text-gray-600 dark:text-gray-300">{s.lot_id || '—'}</span>,
        location: <span className="text-sm text-gray-600 dark:text-gray-300">{s.current_location_id || '—'}</span>,
        received: <span className="text-sm text-gray-500 dark:text-gray-400">{s.received_at ? new Date(s.received_at).toLocaleDateString() : '—'}</span>,
        sold: <span className="text-sm text-gray-500 dark:text-gray-400">{s.sold_at ? new Date(s.sold_at).toLocaleDateString() : '—'}</span>,
        actions: (
            <div className="flex justify-end gap-2">
                <button onClick={() => { setEditing(s); setModalOpen(true); }} className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md"><PencilIcon className="w-4 h-4" /></button>
            </div>
        ),
    });

    return (
        <div className="space-y-4">
            <div className="flex items-end justify-between flex-wrap gap-3">
                <div className="flex gap-3 items-end flex-wrap">
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Producto (SKU)</label>
                        <input type="text" value={filters.product_id || ''} onChange={(e) => setFilters({ product_id: e.target.value || undefined })} placeholder="Filtrar..." className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Lote (ID)</label>
                        <input type="number" value={filters.lot_id ?? ''} onChange={(e) => setFilters({ lot_id: e.target.value ? Number(e.target.value) : undefined })} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white w-32" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Ubicación (ID)</label>
                        <input type="number" value={filters.location_id ?? ''} onChange={(e) => setFilters({ location_id: e.target.value ? Number(e.target.value) : undefined })} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white w-32" />
                    </div>
                </div>
                <button onClick={() => { setEditing(null); setModalOpen(true); }} className="px-5 py-3 btn-business-primary text-white font-bold rounded-lg shadow-lg hover:shadow-2xl transition-all duration-300 transform hover:scale-105 flex items-center gap-2">
                    <PlusIcon className="w-4 h-4" /> Nueva serie
                </button>
            </div>

            {error && <Alert type="error">{error}</Alert>}
            {loading && serials.length === 0 ? <div className="flex justify-center p-8"><Spinner size="lg" /></div> : (
                <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                    <Table columns={columns} data={serials.map(renderRow)} keyExtractor={(_, i) => String(serials[i]?.id || i)} emptyMessage="Sin series" loading={loading} pagination={{ currentPage: page, totalPages, totalItems: total, itemsPerPage: pageSize, onPageChange: setPage, onItemsPerPageChange: setPageSize }} />
                </div>
            )}

            {modalOpen && <SerialFormModal businessId={businessId} serial={editing} onClose={() => { setModalOpen(false); setEditing(null); }} onSuccess={() => { setModalOpen(false); setEditing(null); refresh(); }} />}
        </div>
    );
}

function UoMsSection({ uomsHook, businessId }: { uomsHook: ReturnType<typeof useUoMs>; businessId?: number }) {
    const { uoms, loading, error } = uomsHook;
    const { result, loading: converting, error: convErr, convert } = useUoMConverter(businessId);
    const [convForm, setConvForm] = useState({ product_id: '', from_uom_code: '', to_uom_code: '', quantity: 1 });

    const typeStyles: Record<string, string> = {
        unit: 'bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200',
        weight: 'bg-purple-100 dark:bg-purple-900 text-purple-800 dark:text-purple-200',
        volume: 'bg-cyan-100 dark:bg-cyan-900 text-cyan-800 dark:text-cyan-200',
        length: 'bg-emerald-100 dark:bg-emerald-900 text-emerald-800 dark:text-emerald-200',
        package: 'bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200',
    };

    const columns = [
        { key: 'code', label: 'Código' },
        { key: 'name', label: 'Nombre' },
        { key: 'type', label: 'Tipo', align: 'center' as const },
        { key: 'status', label: 'Estado', align: 'center' as const },
    ];

    const renderRow = (u: UnitOfMeasure) => ({
        code: <span className="font-mono font-medium text-gray-900 dark:text-white">{u.code}</span>,
        name: <span className="text-sm text-gray-700 dark:text-gray-200">{u.name}</span>,
        type: (
            <div className="flex justify-center">
                <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${typeStyles[u.type] || 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300'}`}>{u.type}</span>
            </div>
        ),
        status: (
            <div className="flex justify-center">
                <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${u.is_active ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200' : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300'}`}>{u.is_active ? 'Activo' : 'Inactivo'}</span>
            </div>
        ),
    });

    const handleConvert = async (e: React.FormEvent) => {
        e.preventDefault();
        await convert({ ...convForm, quantity: Number(convForm.quantity) });
    };

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
            <div className="lg:col-span-2">
                {error && <Alert type="error">{error}</Alert>}
                {loading && uoms.length === 0 ? <div className="flex justify-center p-8"><Spinner size="lg" /></div> : (
                    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                        <Table columns={columns} data={uoms.map(renderRow)} keyExtractor={(_, i) => String(uoms[i]?.id || i)} emptyMessage="Sin unidades" loading={loading} />
                    </div>
                )}
            </div>

            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4">
                <h2 className="text-base font-semibold text-gray-900 dark:text-white flex items-center gap-2 mb-4">
                    <ArrowsRightLeftIcon className="w-5 h-5 text-purple-600" /> Conversor
                </h2>
                <form onSubmit={handleConvert} className="space-y-3">
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Producto (SKU)</label>
                        <input required value={convForm.product_id} onChange={(e) => setConvForm({ ...convForm, product_id: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                    </div>
                    <div className="grid grid-cols-2 gap-2">
                        <div>
                            <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">De</label>
                            <select required value={convForm.from_uom_code} onChange={(e) => setConvForm({ ...convForm, from_uom_code: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                                <option value="">—</option>
                                {uoms.map((u) => <option key={u.code} value={u.code}>{u.code}</option>)}
                            </select>
                        </div>
                        <div>
                            <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">A</label>
                            <select required value={convForm.to_uom_code} onChange={(e) => setConvForm({ ...convForm, to_uom_code: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                                <option value="">—</option>
                                {uoms.map((u) => <option key={u.code} value={u.code}>{u.code}</option>)}
                            </select>
                        </div>
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Cantidad</label>
                        <input required type="number" min={0} step="0.001" value={convForm.quantity} onChange={(e) => setConvForm({ ...convForm, quantity: Number(e.target.value) })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                    </div>
                    <Button type="submit" variant="primary" disabled={converting}>{converting ? 'Convirtiendo...' : 'Convertir'}</Button>
                </form>

                {convErr && <div className="mt-3"><Alert type="error">{convErr}</Alert></div>}
                {result && (
                    <div className="mt-4 p-3 rounded-md bg-purple-50 dark:bg-purple-900/20 border border-purple-200 dark:border-purple-800">
                        <p className="text-sm text-purple-900 dark:text-purple-100">
                            <span className="font-semibold">{result.input_quantity} {result.from_uom_code}</span> ={' '}
                            <span className="font-semibold">{result.converted_quantity} {result.to_uom_code}</span>
                        </p>
                        <p className="text-xs text-purple-700 dark:text-purple-300 mt-1">
                            Base: {result.base_unit_quantity} {result.base_uom_code}
                        </p>
                    </div>
                )}
            </div>
        </div>
    );
}

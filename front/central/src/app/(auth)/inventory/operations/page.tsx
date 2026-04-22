'use client';

import { useState } from 'react';
import { PlusIcon, PencilIcon, TrashIcon, CheckIcon, BoltIcon, CheckCircleIcon, XCircleIcon, UserPlusIcon, PlayIcon } from '@heroicons/react/24/outline';
import { Alert, Table, Spinner, Button } from '@/shared/ui';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { FormModal } from '@/shared/ui/form-modal';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useInventoryBusiness } from '@/shared/contexts/inventory-business-context';
import { usePutawayRules, usePutawaySuggestions, useReplenishmentTasks, useCrossDockLinks } from '@/services/modules/inventory/ui/hooks/useOperations';
import {
    deletePutawayRuleAction,
    confirmPutawayAction,
    assignReplenishmentAction,
    completeReplenishmentAction,
    cancelReplenishmentAction,
    detectReplenishmentAction,
    createCrossDockLinkAction,
    executeCrossDockAction,
} from '@/services/modules/inventory/infra/actions/operations';
import { PutawayRule, PutawaySuggestion, ReplenishmentTask, CrossDockLink } from '@/services/modules/inventory/domain/operations-types';
import PutawayRuleFormModal from '@/services/modules/inventory/ui/components/PutawayRuleFormModal';

type Tab = 'putaway' | 'replenishment' | 'crossdock';

export default function OperationsPage() {
    const { isSuperAdmin } = usePermissions();
    const { selectedBusinessId } = useInventoryBusiness();
    const businessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;
    const requiresBusiness = isSuperAdmin && selectedBusinessId === null;

    const [tab, setTab] = useState<Tab>('putaway');

    const rules = usePutawayRules({ business_id: businessId });
    const sugs = usePutawaySuggestions({ business_id: businessId });
    const repl = useReplenishmentTasks({ business_id: businessId });
    const cross = useCrossDockLinks({ business_id: businessId });

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="space-y-4">                {requiresBusiness ? <Alert type="info">Selecciona un negocio.</Alert> : (
                    <>
                        <div className="flex gap-2 border-b border-gray-200 dark:border-gray-700">
                            <TabButton active={tab === 'putaway'} onClick={() => setTab('putaway')} label="Put-away" count={rules.total + sugs.total} />
                            <TabButton active={tab === 'replenishment'} onClick={() => setTab('replenishment')} label="Reposición" count={repl.total} />
                            <TabButton active={tab === 'crossdock'} onClick={() => setTab('crossdock')} label="Cross-dock" count={cross.total} />
                        </div>

                        {tab === 'putaway' && <PutawaySection rules={rules} sugs={sugs} businessId={businessId} />}
                        {tab === 'replenishment' && <ReplenishmentSection hook={repl} businessId={businessId} />}
                        {tab === 'crossdock' && <CrossDockSection hook={cross} businessId={businessId} />}
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

function PutawaySection({ rules, sugs, businessId }: { rules: ReturnType<typeof usePutawayRules>; sugs: ReturnType<typeof usePutawaySuggestions>; businessId?: number }) {
    const [sub, setSub] = useState<'rules' | 'suggestions'>('rules');
    const [modalOpen, setModalOpen] = useState(false);
    const [editing, setEditing] = useState<PutawayRule | null>(null);
    const [deleting, setDeleting] = useState<PutawayRule | null>(null);
    const [confirming, setConfirming] = useState<PutawaySuggestion | null>(null);
    const [confirmLocation, setConfirmLocation] = useState(0);

    const handleDelete = async () => {
        if (!deleting) return;
        await deletePutawayRuleAction(deleting.id, businessId);
        setDeleting(null);
        rules.refresh();
    };

    const handleConfirmSugg = async () => {
        if (!confirming || !confirmLocation) return;
        await confirmPutawayAction(confirming.id, { actual_location_id: confirmLocation }, businessId);
        setConfirming(null);
        setConfirmLocation(0);
        sugs.refresh();
    };

    const statusStyles: Record<string, string> = {
        pending: 'bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200',
        confirmed: 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200',
        canceled: 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300',
    };

    const rulesCols = [
        { key: 'product', label: 'Producto / Cat' },
        { key: 'zone', label: 'Zona destino', align: 'center' as const },
        { key: 'priority', label: 'Prioridad', align: 'center' as const },
        { key: 'strategy', label: 'Estrategia', align: 'center' as const },
        { key: 'status', label: 'Estado', align: 'center' as const },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];
    const rulesRow = (r: PutawayRule) => ({
        product: <span className="text-sm text-gray-700 dark:text-gray-200">{r.product_id ? <span className="font-mono">{r.product_id}</span> : r.category_id ? `Categoría #${r.category_id}` : <span className="text-gray-400">Global</span>}</span>,
        zone: <span className="text-sm font-medium text-gray-900 dark:text-white">#{r.target_zone_id}</span>,
        priority: <span className="text-sm text-gray-600 dark:text-gray-300">{r.priority}</span>,
        strategy: <span className="text-xs inline-flex items-center px-2 py-0.5 rounded-full bg-purple-100 dark:bg-purple-900 text-purple-800 dark:text-purple-200">{r.strategy}</span>,
        status: <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${r.is_active ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200' : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300'}`}>{r.is_active ? 'Activa' : 'Inactiva'}</span>,
        actions: (
            <div className="flex justify-end gap-2">
                <button onClick={() => { setEditing(r); setModalOpen(true); }} className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md"><PencilIcon className="w-4 h-4" /></button>
                <button onClick={() => setDeleting(r)} className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md"><TrashIcon className="w-4 h-4" /></button>
            </div>
        ),
    });

    const sugsCols = [
        { key: 'id', label: '#', align: 'center' as const },
        { key: 'product', label: 'Producto' },
        { key: 'qty', label: 'Cantidad', align: 'center' as const },
        { key: 'location', label: 'Ubicación sugerida', align: 'center' as const },
        { key: 'reason', label: 'Razón' },
        { key: 'status', label: 'Estado', align: 'center' as const },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];
    const sugsRow = (s: PutawaySuggestion) => ({
        id: <span className="text-xs text-gray-500 font-mono">#{s.id}</span>,
        product: <span className="text-sm font-mono text-gray-700 dark:text-gray-200">{s.product_id}</span>,
        qty: <span className="text-sm font-medium text-gray-900 dark:text-white">{s.quantity}</span>,
        location: <span className="text-sm text-gray-600 dark:text-gray-300">#{s.recommended_location_id}</span>,
        reason: <span className="text-xs text-gray-500 dark:text-gray-400">{s.reason}</span>,
        status: <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${statusStyles[s.status] || statusStyles.canceled}`}>{s.status}</span>,
        actions: (
            <div className="flex justify-end">
                {s.status === 'pending' && (
                    <button onClick={() => { setConfirming(s); setConfirmLocation(s.recommended_location_id); }} className="p-2 bg-green-500 hover:bg-green-600 text-white rounded-md" title="Confirmar">
                        <CheckIcon className="w-4 h-4" />
                    </button>
                )}
            </div>
        ),
    });

    return (
        <div className="space-y-4">
            <div className="flex items-center justify-between flex-wrap gap-3">
                <div className="flex gap-2">
                    <button onClick={() => setSub('rules')} className={`px-3 py-1.5 text-sm font-medium rounded-md ${sub === 'rules' ? 'bg-purple-100 dark:bg-purple-900/40 text-purple-700 dark:text-purple-200' : 'text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700'}`}>Reglas ({rules.total})</button>
                    <button onClick={() => setSub('suggestions')} className={`px-3 py-1.5 text-sm font-medium rounded-md ${sub === 'suggestions' ? 'bg-purple-100 dark:bg-purple-900/40 text-purple-700 dark:text-purple-200' : 'text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700'}`}>Sugerencias ({sugs.total})</button>
                </div>
                {sub === 'rules' && (
                    <button onClick={() => { setEditing(null); setModalOpen(true); }} className="px-5 py-3 btn-business-primary text-white font-bold rounded-lg shadow-lg hover:shadow-2xl transition-all duration-300 transform hover:scale-105 flex items-center gap-2">
                        <PlusIcon className="w-4 h-4" /> Nueva regla
                    </button>
                )}
            </div>

            {sub === 'rules' ? (
                <>
                    {rules.error && <Alert type="error">{rules.error}</Alert>}
                    {rules.loading && rules.rules.length === 0 ? <div className="flex justify-center p-8"><Spinner size="lg" /></div> : (
                        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                            <Table columns={rulesCols} data={rules.rules.map(rulesRow)} keyExtractor={(_, i) => String(rules.rules[i]?.id || i)} emptyMessage="Sin reglas" loading={rules.loading} pagination={{ currentPage: rules.page, totalPages: rules.totalPages, totalItems: rules.total, itemsPerPage: rules.pageSize, onPageChange: rules.setPage, onItemsPerPageChange: rules.setPageSize }} />
                        </div>
                    )}
                </>
            ) : (
                <>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Estado</label>
                        <select value={sugs.status} onChange={(e) => sugs.setStatus(e.target.value)} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                            <option value="">Todos</option>
                            <option value="pending">Pendiente</option>
                            <option value="confirmed">Confirmada</option>
                            <option value="canceled">Cancelada</option>
                        </select>
                    </div>
                    {sugs.error && <Alert type="error">{sugs.error}</Alert>}
                    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                        <Table columns={sugsCols} data={sugs.suggestions.map(sugsRow)} keyExtractor={(_, i) => String(sugs.suggestions[i]?.id || i)} emptyMessage="Sin sugerencias" loading={sugs.loading} pagination={{ currentPage: sugs.page, totalPages: sugs.totalPages, totalItems: sugs.total, itemsPerPage: sugs.pageSize, onPageChange: sugs.setPage, onItemsPerPageChange: sugs.setPageSize }} />
                    </div>
                </>
            )}

            {modalOpen && <PutawayRuleFormModal businessId={businessId} rule={editing} onClose={() => { setModalOpen(false); setEditing(null); }} onSuccess={() => { setModalOpen(false); setEditing(null); rules.refresh(); }} />}
            {deleting && <ConfirmModal isOpen={true} onClose={() => setDeleting(null)} onConfirm={handleDelete} title="Eliminar regla" message={`Se eliminará la regla #${deleting.id}. Acción irreversible.`} confirmText="Eliminar" type="danger" />}

            {confirming && (
                <FormModal isOpen={true} onClose={() => setConfirming(null)} title={`Confirmar put-away #${confirming.id}`}>
                    <div className="p-6 space-y-3">
                        <p className="text-sm text-gray-600 dark:text-gray-300">Producto <span className="font-mono">{confirming.product_id}</span> · qty <span className="font-bold">{confirming.quantity}</span></p>
                        <div>
                            <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Ubicación real (ID)</label>
                            <input type="number" value={confirmLocation} onChange={(e) => setConfirmLocation(Number(e.target.value))} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                            <p className="text-xs text-gray-500 mt-1">Sugerida: #{confirming.recommended_location_id}</p>
                        </div>
                        <div className="flex justify-end gap-2 pt-4 border-t border-gray-200 dark:border-gray-700">
                            <Button type="button" variant="outline" onClick={() => setConfirming(null)}>Cancelar</Button>
                            <Button type="button" variant="primary" onClick={handleConfirmSugg} disabled={!confirmLocation}>Confirmar</Button>
                        </div>
                    </div>
                </FormModal>
            )}
        </div>
    );
}

function ReplenishmentSection({ hook, businessId }: { hook: ReturnType<typeof useReplenishmentTasks>; businessId?: number }) {
    const { tasks, loading, error, page, pageSize, total, totalPages, filters, setPage, setPageSize, setFilters, refresh } = hook;
    const [detecting, setDetecting] = useState(false);
    const [detectResult, setDetectResult] = useState<string | null>(null);
    const [assigning, setAssigning] = useState<ReplenishmentTask | null>(null);
    const [userId, setUserId] = useState('');

    const handleDetect = async () => {
        setDetecting(true);
        setDetectResult(null);
        try {
            const r = await detectReplenishmentAction(businessId);
            if (r.success) setDetectResult(`Se crearon ${r.data?.created || 0} tareas nuevas`);
            else setDetectResult(r.error || 'Error en detección');
            refresh();
        } finally { setDetecting(false); }
    };

    const handleComplete = async (id: number) => {
        await completeReplenishmentAction(id, {}, businessId);
        refresh();
    };

    const handleCancel = async (id: number) => {
        if (!window.confirm('¿Cancelar esta tarea de reposición?')) return;
        await cancelReplenishmentAction(id, 'Cancelada por usuario', businessId);
        refresh();
    };

    const handleAssign = async () => {
        if (!assigning || !userId) return;
        await assignReplenishmentAction(assigning.id, { user_id: Number(userId) }, businessId);
        setAssigning(null);
        setUserId('');
        refresh();
    };

    const statusStyles: Record<string, string> = {
        pending: 'bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200',
        in_progress: 'bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200',
        completed: 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200',
        canceled: 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300',
    };

    const columns = [
        { key: 'id', label: '#', align: 'center' as const },
        { key: 'product', label: 'Producto' },
        { key: 'warehouse', label: 'Bodega', align: 'center' as const },
        { key: 'qty', label: 'Cantidad', align: 'center' as const },
        { key: 'route', label: 'Ruta', align: 'center' as const },
        { key: 'trigger', label: 'Origen', align: 'center' as const },
        { key: 'assignee', label: 'Asignado', align: 'center' as const },
        { key: 'status', label: 'Estado', align: 'center' as const },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];

    const renderRow = (t: ReplenishmentTask) => ({
        id: <span className="text-xs text-gray-500 font-mono">#{t.id}</span>,
        product: <span className="text-sm font-mono text-gray-700 dark:text-gray-200">{t.product_id}</span>,
        warehouse: <span className="text-sm text-gray-600 dark:text-gray-300">#{t.warehouse_id}</span>,
        qty: <span className="text-sm font-bold text-gray-900 dark:text-white">{t.quantity}</span>,
        route: <span className="text-xs text-gray-500 dark:text-gray-400">{t.from_location_id ? `#${t.from_location_id}` : '—'} → {t.to_location_id ? `#${t.to_location_id}` : '—'}</span>,
        trigger: <span className="text-xs inline-flex items-center px-2 py-0.5 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-200">{t.triggered_by}</span>,
        assignee: <span className="text-xs text-gray-500">{t.assigned_to_id ? `User #${t.assigned_to_id}` : '—'}</span>,
        status: <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${statusStyles[t.status] || statusStyles.canceled}`}>{t.status}</span>,
        actions: (
            <div className="flex justify-end gap-1">
                {t.status === 'pending' && <button onClick={() => setAssigning(t)} className="p-2 bg-blue-500 hover:bg-blue-600 text-white rounded-md" title="Asignar"><UserPlusIcon className="w-4 h-4" /></button>}
                {t.status === 'in_progress' && <button onClick={() => handleComplete(t.id)} className="p-2 bg-green-500 hover:bg-green-600 text-white rounded-md" title="Completar"><CheckCircleIcon className="w-4 h-4" /></button>}
                {(t.status === 'pending' || t.status === 'in_progress') && <button onClick={() => handleCancel(t.id)} className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md" title="Cancelar"><XCircleIcon className="w-4 h-4" /></button>}
            </div>
        ),
    });

    return (
        <div className="space-y-4">
            {detectResult && <Alert type="info">{detectResult}</Alert>}

            <div className="flex items-end justify-between flex-wrap gap-3">
                <div className="flex gap-3 items-end flex-wrap">
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Bodega (ID)</label>
                        <input type="number" value={filters.warehouse_id ?? ''} onChange={(e) => setFilters({ warehouse_id: e.target.value ? Number(e.target.value) : undefined })} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm w-32" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Estado</label>
                        <select value={filters.status || ''} onChange={(e) => setFilters({ status: e.target.value || undefined })} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                            <option value="">Todos</option>
                            <option value="pending">Pendientes</option>
                            <option value="in_progress">En progreso</option>
                            <option value="completed">Completadas</option>
                            <option value="canceled">Canceladas</option>
                        </select>
                    </div>
                </div>
                <button onClick={handleDetect} disabled={detecting} className="px-5 py-3 btn-business-primary text-white font-bold rounded-lg shadow-lg hover:shadow-2xl transition-all duration-300 transform hover:scale-105 flex items-center gap-2 disabled:opacity-60">
                    <BoltIcon className="w-4 h-4" /> {detecting ? 'Escaneando...' : 'Detectar reposiciones'}
                </button>
            </div>

            {error && <Alert type="error">{error}</Alert>}
            {loading && tasks.length === 0 ? <div className="flex justify-center p-8"><Spinner size="lg" /></div> : (
                <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                    <Table columns={columns} data={tasks.map(renderRow)} keyExtractor={(_, i) => String(tasks[i]?.id || i)} emptyMessage="Sin tareas de reposición" loading={loading} pagination={{ currentPage: page, totalPages, totalItems: total, itemsPerPage: pageSize, onPageChange: setPage, onItemsPerPageChange: setPageSize }} />
                </div>
            )}

            {assigning && (
                <FormModal isOpen={true} onClose={() => { setAssigning(null); setUserId(''); }} title={`Asignar tarea #${assigning.id}`}>
                    <div className="p-6 space-y-3">
                        <p className="text-sm text-gray-600 dark:text-gray-300">Producto <span className="font-mono">{assigning.product_id}</span> · qty {assigning.quantity}</p>
                        <div>
                            <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Usuario (ID)</label>
                            <input type="number" value={userId} onChange={(e) => setUserId(e.target.value)} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                        </div>
                        <div className="flex justify-end gap-2 pt-4 border-t border-gray-200 dark:border-gray-700">
                            <Button type="button" variant="outline" onClick={() => { setAssigning(null); setUserId(''); }}>Cancelar</Button>
                            <Button type="button" variant="primary" onClick={handleAssign} disabled={!userId}>Asignar</Button>
                        </div>
                    </div>
                </FormModal>
            )}
        </div>
    );
}

function CrossDockSection({ hook, businessId }: { hook: ReturnType<typeof useCrossDockLinks>; businessId?: number }) {
    const { links, loading, error, page, pageSize, total, totalPages, filters, setPage, setPageSize, setFilters, refresh } = hook;
    const [modalOpen, setModalOpen] = useState(false);
    const [form, setForm] = useState({ inbound_shipment_id: '', outbound_order_id: '', product_id: '', quantity: 1 });
    const [submitting, setSubmitting] = useState(false);
    const [formError, setFormError] = useState<string | null>(null);

    const submit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSubmitting(true);
        setFormError(null);
        try {
            const r = await createCrossDockLinkAction({
                inbound_shipment_id: form.inbound_shipment_id ? Number(form.inbound_shipment_id) : null,
                outbound_order_id: form.outbound_order_id,
                product_id: form.product_id,
                quantity: Number(form.quantity),
            }, businessId);
            if (!r.success) { setFormError(r.error || 'Error'); return; }
            setModalOpen(false);
            setForm({ inbound_shipment_id: '', outbound_order_id: '', product_id: '', quantity: 1 });
            refresh();
        } finally { setSubmitting(false); }
    };

    const handleExecute = async (id: number) => {
        await executeCrossDockAction(id, businessId);
        refresh();
    };

    const statusStyles: Record<string, string> = {
        pending: 'bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200',
        executed: 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200',
        canceled: 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300',
    };

    const columns = [
        { key: 'id', label: '#', align: 'center' as const },
        { key: 'inbound', label: 'Inbound', align: 'center' as const },
        { key: 'outbound', label: 'Orden salida' },
        { key: 'product', label: 'Producto' },
        { key: 'qty', label: 'Cantidad', align: 'center' as const },
        { key: 'status', label: 'Estado', align: 'center' as const },
        { key: 'executed_at', label: 'Ejecutada', align: 'center' as const },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];

    const renderRow = (l: CrossDockLink) => ({
        id: <span className="text-xs text-gray-500 font-mono">#{l.id}</span>,
        inbound: <span className="text-xs text-gray-600 dark:text-gray-300">{l.inbound_shipment_id ? `#${l.inbound_shipment_id}` : '—'}</span>,
        outbound: <span className="text-sm font-mono text-gray-900 dark:text-white">{l.outbound_order_id}</span>,
        product: <span className="text-sm font-mono text-gray-700 dark:text-gray-200">{l.product_id}</span>,
        qty: <span className="text-sm font-bold text-gray-900 dark:text-white">{l.quantity}</span>,
        status: <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${statusStyles[l.status] || statusStyles.canceled}`}>{l.status}</span>,
        executed_at: <span className="text-xs text-gray-500">{l.executed_at ? new Date(l.executed_at).toLocaleString() : '—'}</span>,
        actions: (
            <div className="flex justify-end">
                {l.status === 'pending' && (
                    <button onClick={() => handleExecute(l.id)} className="p-2 bg-green-500 hover:bg-green-600 text-white rounded-md" title="Ejecutar"><PlayIcon className="w-4 h-4" /></button>
                )}
            </div>
        ),
    });

    return (
        <div className="space-y-4">
            <div className="flex items-end justify-between flex-wrap gap-3">
                <div className="flex gap-3 items-end flex-wrap">
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Orden salida</label>
                        <input value={filters.outbound_order_id || ''} onChange={(e) => setFilters({ outbound_order_id: e.target.value || undefined })} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Estado</label>
                        <select value={filters.status || ''} onChange={(e) => setFilters({ status: e.target.value || undefined })} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                            <option value="">Todos</option>
                            <option value="pending">Pendiente</option>
                            <option value="executed">Ejecutado</option>
                            <option value="canceled">Cancelado</option>
                        </select>
                    </div>
                </div>
                <button onClick={() => setModalOpen(true)} className="px-5 py-3 btn-business-primary text-white font-bold rounded-lg shadow-lg hover:shadow-2xl transition-all duration-300 transform hover:scale-105 flex items-center gap-2">
                    <PlusIcon className="w-4 h-4" /> Nuevo enlace
                </button>
            </div>

            {error && <Alert type="error">{error}</Alert>}
            {loading && links.length === 0 ? <div className="flex justify-center p-8"><Spinner size="lg" /></div> : (
                <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                    <Table columns={columns} data={links.map(renderRow)} keyExtractor={(_, i) => String(links[i]?.id || i)} emptyMessage="Sin enlaces cross-dock" loading={loading} pagination={{ currentPage: page, totalPages, totalItems: total, itemsPerPage: pageSize, onPageChange: setPage, onItemsPerPageChange: setPageSize }} />
                </div>
            )}

            {modalOpen && (
                <FormModal isOpen={true} onClose={() => setModalOpen(false)} title="Crear enlace cross-dock">
                    <form onSubmit={submit} className="p-6 space-y-4">
                        {formError && <Alert type="error">{formError}</Alert>}
                        <div className="grid grid-cols-2 gap-4">
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Inbound shipment (ID) [opcional]</label>
                                <input type="number" value={form.inbound_shipment_id} onChange={(e) => setForm({ ...form, inbound_shipment_id: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                            </div>
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Orden de salida *</label>
                                <input required value={form.outbound_order_id} onChange={(e) => setForm({ ...form, outbound_order_id: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                            </div>
                        </div>
                        <div className="grid grid-cols-2 gap-4">
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Producto (SKU) *</label>
                                <input required value={form.product_id} onChange={(e) => setForm({ ...form, product_id: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                            </div>
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Cantidad *</label>
                                <input required type="number" min={1} value={form.quantity} onChange={(e) => setForm({ ...form, quantity: Number(e.target.value) })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                            </div>
                        </div>
                        <div className="flex justify-end gap-2 pt-4 border-t border-gray-200 dark:border-gray-700">
                            <Button type="button" variant="outline" onClick={() => setModalOpen(false)} disabled={submitting}>Cancelar</Button>
                            <Button type="submit" variant="primary" disabled={submitting}>{submitting ? 'Creando...' : 'Crear'}</Button>
                        </div>
                    </form>
                </FormModal>
            )}
        </div>
    );
}

'use client';

import { useState } from 'react';
import { PlusIcon, PencilIcon, TrashIcon, PlayIcon, StopIcon, EyeIcon, CheckIcon, XMarkIcon } from '@heroicons/react/24/outline';
import { Alert, Table, Spinner, Button } from '@/shared/ui';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { FormModal } from '@/shared/ui/form-modal';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useInventoryBusiness } from '@/shared/contexts/inventory-business-context';
import { useCountPlans, useCountTasks, useCountLines, useDiscrepancies } from '@/services/modules/inventory/ui/hooks/useAudit';
import {
    deleteCountPlanAction,
    generateCountTaskAction,
    startCountTaskAction,
    finishCountTaskAction,
    submitCountLineAction,
    approveDiscrepancyAction,
    rejectDiscrepancyAction,
} from '@/services/modules/inventory/infra/actions/audit';
import { CycleCountPlan, CycleCountTask, CycleCountLine, InventoryDiscrepancy } from '@/services/modules/inventory/domain/audit-types';
import CountPlanFormModal from '@/services/modules/inventory/ui/components/CountPlanFormModal';

type Tab = 'plans' | 'tasks' | 'discrepancies';

export default function AuditPage() {
    const { isSuperAdmin } = usePermissions();
    const { selectedBusinessId } = useInventoryBusiness();
    const businessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;
    const requiresBusiness = isSuperAdmin && selectedBusinessId === null;

    const [tab, setTab] = useState<Tab>('plans');

    const plans = useCountPlans({ business_id: businessId });
    const tasks = useCountTasks({ business_id: businessId });
    const disc = useDiscrepancies({ business_id: businessId });

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="space-y-4">                {requiresBusiness ? (
                    <Alert type="info">Selecciona un negocio.</Alert>
                ) : (
                    <>
                        <div className="flex gap-2 border-b border-gray-200 dark:border-gray-700">
                            <TabButton active={tab === 'plans'} onClick={() => setTab('plans')} label="Planes" count={plans.total} />
                            <TabButton active={tab === 'tasks'} onClick={() => setTab('tasks')} label="Tareas" count={tasks.total} />
                            <TabButton active={tab === 'discrepancies'} onClick={() => setTab('discrepancies')} label="Discrepancias" count={disc.total} />
                        </div>

                        {tab === 'plans' && <PlansSection hook={plans} businessId={businessId} />}
                        {tab === 'tasks' && <TasksSection hook={tasks} businessId={businessId} />}
                        {tab === 'discrepancies' && <DiscrepanciesSection hook={disc} businessId={businessId} />}
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

function PlansSection({ hook, businessId }: { hook: ReturnType<typeof useCountPlans>; businessId?: number }) {
    const { plans, loading, error, page, pageSize, total, totalPages, setPage, setPageSize, refresh } = hook;
    const [modalOpen, setModalOpen] = useState(false);
    const [editing, setEditing] = useState<CycleCountPlan | null>(null);
    const [deleting, setDeleting] = useState<CycleCountPlan | null>(null);

    const handleDelete = async () => {
        if (!deleting) return;
        await deleteCountPlanAction(deleting.id, businessId);
        setDeleting(null);
        refresh();
    };

    const strategyStyles: Record<string, string> = {
        abc: 'bg-purple-100 dark:bg-purple-900 text-purple-800 dark:text-purple-200',
        zone: 'bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200',
        random: 'bg-cyan-100 dark:bg-cyan-900 text-cyan-800 dark:text-cyan-200',
        full: 'bg-orange-100 dark:bg-orange-900 text-orange-800 dark:text-orange-200',
    };

    const columns = [
        { key: 'name', label: 'Nombre' },
        { key: 'warehouse', label: 'Bodega', align: 'center' as const },
        { key: 'strategy', label: 'Estrategia', align: 'center' as const },
        { key: 'frequency', label: 'Frecuencia', align: 'center' as const },
        { key: 'next_run', label: 'Próxima', align: 'center' as const },
        { key: 'status', label: 'Estado', align: 'center' as const },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];

    const renderRow = (p: CycleCountPlan) => ({
        name: <span className="font-medium text-gray-900 dark:text-white">{p.name}</span>,
        warehouse: <span className="text-sm text-gray-600 dark:text-gray-300">#{p.warehouse_id}</span>,
        strategy: <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${strategyStyles[p.strategy] || 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300'}`}>{p.strategy}</span>,
        frequency: <span className="text-sm text-gray-600 dark:text-gray-300">{p.frequency_days}d</span>,
        next_run: <span className="text-xs text-gray-500">{p.next_run_at ? new Date(p.next_run_at).toLocaleDateString() : '—'}</span>,
        status: <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${p.is_active ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200' : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300'}`}>{p.is_active ? 'Activo' : 'Inactivo'}</span>,
        actions: (
            <div className="flex justify-end gap-2">
                <button onClick={() => { setEditing(p); setModalOpen(true); }} className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md"><PencilIcon className="w-4 h-4" /></button>
                <button onClick={() => setDeleting(p)} className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md"><TrashIcon className="w-4 h-4" /></button>
            </div>
        ),
    });

    return (
        <div className="space-y-4">
            <div className="flex justify-end">
                <button onClick={() => { setEditing(null); setModalOpen(true); }} className="px-5 py-3 btn-business-primary text-white font-bold rounded-lg shadow-lg hover:shadow-2xl transition-all duration-300 transform hover:scale-105 flex items-center gap-2">
                    <PlusIcon className="w-4 h-4" /> Nuevo plan
                </button>
            </div>

            {error && <Alert type="error">{error}</Alert>}
            {loading && plans.length === 0 ? <div className="flex justify-center p-8"><Spinner size="lg" /></div> : (
                <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                    <Table columns={columns} data={plans.map(renderRow)} keyExtractor={(_, i) => String(plans[i]?.id || i)} emptyMessage="Sin planes" loading={loading} pagination={{ currentPage: page, totalPages, totalItems: total, itemsPerPage: pageSize, onPageChange: setPage, onItemsPerPageChange: setPageSize }} />
                </div>
            )}

            {modalOpen && <CountPlanFormModal businessId={businessId} plan={editing} onClose={() => { setModalOpen(false); setEditing(null); }} onSuccess={() => { setModalOpen(false); setEditing(null); refresh(); }} />}
            {deleting && <ConfirmModal isOpen={true} onClose={() => setDeleting(null)} onConfirm={handleDelete} title="Eliminar plan" message={`Se eliminará el plan "${deleting.name}". Acción irreversible.`} confirmText="Eliminar" type="danger" />}
        </div>
    );
}

function TasksSection({ hook, businessId }: { hook: ReturnType<typeof useCountTasks>; businessId?: number }) {
    const { tasks, loading, error, page, pageSize, total, totalPages, filters, setPage, setPageSize, setFilters, refresh } = hook;
    const [genOpen, setGenOpen] = useState(false);
    const [genForm, setGenForm] = useState({ plan_id: 0, scope_type: '', scope_id: '' });
    const [linesTask, setLinesTask] = useState<CycleCountTask | null>(null);
    const [startTask, setStartTask] = useState<CycleCountTask | null>(null);
    const [startUser, setStartUser] = useState('');

    const handleGenerate = async (e: React.FormEvent) => {
        e.preventDefault();
        const r = await generateCountTaskAction({
            plan_id: Number(genForm.plan_id),
            scope_type: genForm.scope_type || undefined,
            scope_id: genForm.scope_id ? Number(genForm.scope_id) : null,
        }, businessId);
        if (r.success) {
            setGenOpen(false);
            setGenForm({ plan_id: 0, scope_type: '', scope_id: '' });
            refresh();
        }
    };

    const handleStart = async () => {
        if (!startTask || !startUser) return;
        await startCountTaskAction(startTask.id, { user_id: Number(startUser) }, businessId);
        setStartTask(null);
        setStartUser('');
        refresh();
    };

    const handleFinish = async (id: number) => {
        await finishCountTaskAction(id, businessId);
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
        { key: 'plan', label: 'Plan', align: 'center' as const },
        { key: 'warehouse', label: 'Bodega', align: 'center' as const },
        { key: 'scope', label: 'Alcance' },
        { key: 'assignee', label: 'Asignado', align: 'center' as const },
        { key: 'started', label: 'Inicio', align: 'center' as const },
        { key: 'finished', label: 'Fin', align: 'center' as const },
        { key: 'status', label: 'Estado', align: 'center' as const },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];

    const renderRow = (t: CycleCountTask) => ({
        id: <span className="text-xs text-gray-500 font-mono">#{t.id}</span>,
        plan: <span className="text-sm text-gray-600 dark:text-gray-300">#{t.plan_id}</span>,
        warehouse: <span className="text-sm text-gray-600 dark:text-gray-300">#{t.warehouse_id}</span>,
        scope: <span className="text-xs text-gray-500">{t.scope_type}{t.scope_id ? `:#${t.scope_id}` : ''}</span>,
        assignee: <span className="text-xs text-gray-500">{t.assigned_to_id ? `User #${t.assigned_to_id}` : '—'}</span>,
        started: <span className="text-xs text-gray-500">{t.started_at ? new Date(t.started_at).toLocaleString() : '—'}</span>,
        finished: <span className="text-xs text-gray-500">{t.finished_at ? new Date(t.finished_at).toLocaleString() : '—'}</span>,
        status: <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${statusStyles[t.status] || statusStyles.canceled}`}>{t.status}</span>,
        actions: (
            <div className="flex justify-end gap-1">
                <button onClick={() => setLinesTask(t)} className="p-2 bg-blue-500 hover:bg-blue-600 text-white rounded-md" title="Ver líneas"><EyeIcon className="w-4 h-4" /></button>
                {t.status === 'pending' && <button onClick={() => setStartTask(t)} className="p-2 bg-green-500 hover:bg-green-600 text-white rounded-md" title="Iniciar"><PlayIcon className="w-4 h-4" /></button>}
                {t.status === 'in_progress' && <button onClick={() => handleFinish(t.id)} className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md" title="Finalizar"><StopIcon className="w-4 h-4" /></button>}
            </div>
        ),
    });

    return (
        <div className="space-y-4">
            <div className="flex items-end justify-between flex-wrap gap-3">
                <div className="flex gap-3 items-end flex-wrap">
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Bodega (ID)</label>
                        <input type="number" value={filters.warehouse_id ?? ''} onChange={(e) => setFilters({ warehouse_id: e.target.value ? Number(e.target.value) : undefined })} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm w-32" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Plan (ID)</label>
                        <input type="number" value={filters.plan_id ?? ''} onChange={(e) => setFilters({ plan_id: e.target.value ? Number(e.target.value) : undefined })} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm w-32" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Estado</label>
                        <select value={filters.status || ''} onChange={(e) => setFilters({ status: e.target.value || undefined })} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                            <option value="">Todos</option>
                            <option value="pending">Pendientes</option>
                            <option value="in_progress">En progreso</option>
                            <option value="completed">Completadas</option>
                        </select>
                    </div>
                </div>
                <button onClick={() => setGenOpen(true)} className="px-5 py-3 btn-business-primary text-white font-bold rounded-lg shadow-lg hover:shadow-2xl transition-all duration-300 transform hover:scale-105 flex items-center gap-2">
                    <PlusIcon className="w-4 h-4" /> Generar tarea
                </button>
            </div>

            {error && <Alert type="error">{error}</Alert>}
            {loading && tasks.length === 0 ? <div className="flex justify-center p-8"><Spinner size="lg" /></div> : (
                <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                    <Table columns={columns} data={tasks.map(renderRow)} keyExtractor={(_, i) => String(tasks[i]?.id || i)} emptyMessage="Sin tareas" loading={loading} pagination={{ currentPage: page, totalPages, totalItems: total, itemsPerPage: pageSize, onPageChange: setPage, onItemsPerPageChange: setPageSize }} />
                </div>
            )}

            {genOpen && (
                <FormModal isOpen={true} onClose={() => setGenOpen(false)} title="Generar tarea de conteo">
                    <form onSubmit={handleGenerate} className="p-6 space-y-4">
                        <div className="grid grid-cols-3 gap-4">
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Plan (ID) *</label>
                                <input required type="number" min={1} value={genForm.plan_id} onChange={(e) => setGenForm({ ...genForm, plan_id: Number(e.target.value) })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                            </div>
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Tipo alcance</label>
                                <select value={genForm.scope_type} onChange={(e) => setGenForm({ ...genForm, scope_type: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                                    <option value="">(auto)</option>
                                    <option value="zone">Zona</option>
                                    <option value="aisle">Pasillo</option>
                                    <option value="rack">Rack</option>
                                </select>
                            </div>
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">ID alcance</label>
                                <input type="number" value={genForm.scope_id} onChange={(e) => setGenForm({ ...genForm, scope_id: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                            </div>
                        </div>
                        <div className="flex justify-end gap-2 pt-4 border-t border-gray-200 dark:border-gray-700">
                            <Button type="button" variant="outline" onClick={() => setGenOpen(false)}>Cancelar</Button>
                            <Button type="submit" variant="primary">Generar</Button>
                        </div>
                    </form>
                </FormModal>
            )}

            {startTask && (
                <FormModal isOpen={true} onClose={() => { setStartTask(null); setStartUser(''); }} title={`Iniciar tarea #${startTask.id}`}>
                    <div className="p-6 space-y-4">
                        <div>
                            <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Usuario responsable (ID)</label>
                            <input type="number" value={startUser} onChange={(e) => setStartUser(e.target.value)} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                        </div>
                        <div className="flex justify-end gap-2 pt-4 border-t border-gray-200 dark:border-gray-700">
                            <Button type="button" variant="outline" onClick={() => { setStartTask(null); setStartUser(''); }}>Cancelar</Button>
                            <Button type="button" variant="primary" onClick={handleStart} disabled={!startUser}>Iniciar</Button>
                        </div>
                    </div>
                </FormModal>
            )}

            {linesTask && <CountLinesModal task={linesTask} businessId={businessId} onClose={() => { setLinesTask(null); refresh(); }} />}
        </div>
    );
}

function CountLinesModal({ task, businessId, onClose }: { task: CycleCountTask; businessId?: number; onClose: () => void }) {
    const { lines, loading, error, refresh } = useCountLines({ task_id: task.id, business_id: businessId, page: 1, page_size: 100 });
    const [submitting, setSubmitting] = useState<number | null>(null);
    const [counts, setCounts] = useState<Record<number, string>>({});

    const submit = async (line: CycleCountLine) => {
        const val = counts[line.id];
        if (val === undefined || val === '') return;
        setSubmitting(line.id);
        try {
            await submitCountLineAction(line.id, { counted_qty: Number(val) }, businessId);
            refresh();
        } finally { setSubmitting(null); }
    };

    return (
        <FormModal isOpen={true} onClose={onClose} title={`Líneas de conteo · Task #${task.id}`} size="xl">
            <div className="p-6 space-y-3">
                {error && <Alert type="error">{error}</Alert>}
                {loading ? <div className="flex justify-center p-4"><Spinner /></div> : (
                    <div className="max-h-[60vh] overflow-y-auto">
                        <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700 text-sm">
                            <thead className="bg-gray-50 dark:bg-gray-900/60 sticky top-0">
                                <tr className="text-xs uppercase tracking-wider text-gray-500 dark:text-gray-400">
                                    <th className="px-3 py-2 text-left">Producto</th>
                                    <th className="px-3 py-2 text-center">Ubicación</th>
                                    <th className="px-3 py-2 text-center">Lote</th>
                                    <th className="px-3 py-2 text-center">Esperado</th>
                                    <th className="px-3 py-2 text-center">Contado</th>
                                    <th className="px-3 py-2 text-center">Variación</th>
                                    <th className="px-3 py-2 text-center">Estado</th>
                                    <th className="px-3 py-2 text-right">Acción</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-100 dark:divide-gray-700/50">
                                {lines.map((l) => (
                                    <tr key={l.id}>
                                        <td className="px-3 py-2 font-mono text-xs">{l.product_id}</td>
                                        <td className="px-3 py-2 text-center text-xs text-gray-500">{l.location_id ? `#${l.location_id}` : '—'}</td>
                                        <td className="px-3 py-2 text-center text-xs text-gray-500">{l.lot_id ? `#${l.lot_id}` : '—'}</td>
                                        <td className="px-3 py-2 text-center font-medium">{l.expected_qty}</td>
                                        <td className="px-3 py-2 text-center">
                                            {l.counted_qty !== null ? l.counted_qty : (
                                                <input type="number" value={counts[l.id] ?? ''} onChange={(e) => setCounts({ ...counts, [l.id]: e.target.value })} className="w-20 px-2 py-1 border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                                            )}
                                        </td>
                                        <td className={`px-3 py-2 text-center font-semibold ${l.variance === 0 ? 'text-gray-400' : l.variance > 0 ? 'text-green-600' : 'text-red-600'}`}>{l.variance !== null && l.counted_qty !== null ? (l.variance > 0 ? `+${l.variance}` : l.variance) : '—'}</td>
                                        <td className="px-3 py-2 text-center">
                                            <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${l.status === 'counted' ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200' : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300'}`}>{l.status}</span>
                                        </td>
                                        <td className="px-3 py-2 text-right">
                                            {l.counted_qty === null && (
                                                <button onClick={() => submit(l)} disabled={submitting === l.id || !counts[l.id]} className="px-3 py-1 bg-purple-600 hover:bg-purple-700 text-white text-xs rounded disabled:opacity-50">
                                                    {submitting === l.id ? '...' : 'Enviar'}
                                                </button>
                                            )}
                                        </td>
                                    </tr>
                                ))}
                                {lines.length === 0 && <tr><td colSpan={8} className="text-center py-6 text-gray-400 text-sm">Sin líneas</td></tr>}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </FormModal>
    );
}

function DiscrepanciesSection({ hook, businessId }: { hook: ReturnType<typeof useDiscrepancies>; businessId?: number }) {
    const { discrepancies, loading, error, page, pageSize, total, totalPages, filters, setPage, setPageSize, setFilters, refresh } = hook;
    const [reviewing, setReviewing] = useState<{ disc: InventoryDiscrepancy; action: 'approve' | 'reject' } | null>(null);
    const [notes, setNotes] = useState('');

    const handleReview = async () => {
        if (!reviewing) return;
        const { disc, action } = reviewing;
        if (action === 'approve') await approveDiscrepancyAction(disc.id, { notes }, businessId);
        else await rejectDiscrepancyAction(disc.id, { reason: notes }, businessId);
        setReviewing(null);
        setNotes('');
        refresh();
    };

    const statusStyles: Record<string, string> = {
        open: 'bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200',
        approved: 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200',
        rejected: 'bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200',
    };

    const columns = [
        { key: 'id', label: '#', align: 'center' as const },
        { key: 'task', label: 'Tarea', align: 'center' as const },
        { key: 'line', label: 'Línea', align: 'center' as const },
        { key: 'resolution', label: 'Mov. resol.', align: 'center' as const },
        { key: 'reviewer', label: 'Revisor', align: 'center' as const },
        { key: 'reviewed_at', label: 'Revisado', align: 'center' as const },
        { key: 'notes', label: 'Notas' },
        { key: 'status', label: 'Estado', align: 'center' as const },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];

    const renderRow = (d: InventoryDiscrepancy) => ({
        id: <span className="text-xs text-gray-500 font-mono">#{d.id}</span>,
        task: <span className="text-sm text-gray-600 dark:text-gray-300">#{d.task_id}</span>,
        line: <span className="text-sm text-gray-600 dark:text-gray-300">#{d.line_id}</span>,
        resolution: <span className="text-xs text-gray-500">{d.resolution_movement_id ? `Mov #${d.resolution_movement_id}` : '—'}</span>,
        reviewer: <span className="text-xs text-gray-500">{d.reviewed_by_id ? `User #${d.reviewed_by_id}` : '—'}</span>,
        reviewed_at: <span className="text-xs text-gray-500">{d.reviewed_at ? new Date(d.reviewed_at).toLocaleString() : '—'}</span>,
        notes: <span className="text-xs text-gray-600 dark:text-gray-300 truncate block max-w-xs">{d.notes || '—'}</span>,
        status: <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${statusStyles[d.status] || statusStyles.rejected}`}>{d.status}</span>,
        actions: (
            <div className="flex justify-end gap-1">
                {d.status === 'open' && (
                    <>
                        <button onClick={() => { setReviewing({ disc: d, action: 'approve' }); setNotes(''); }} className="p-2 bg-green-500 hover:bg-green-600 text-white rounded-md" title="Aprobar"><CheckIcon className="w-4 h-4" /></button>
                        <button onClick={() => { setReviewing({ disc: d, action: 'reject' }); setNotes(''); }} className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md" title="Rechazar"><XMarkIcon className="w-4 h-4" /></button>
                    </>
                )}
            </div>
        ),
    });

    return (
        <div className="space-y-4">
            <div className="flex gap-3 items-end flex-wrap">
                <div>
                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Tarea (ID)</label>
                    <input type="number" value={filters.task_id ?? ''} onChange={(e) => setFilters({ task_id: e.target.value ? Number(e.target.value) : undefined })} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm w-32" />
                </div>
                <div>
                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Estado</label>
                    <select value={filters.status || ''} onChange={(e) => setFilters({ status: e.target.value || undefined })} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                        <option value="">Todos</option>
                        <option value="open">Abiertas</option>
                        <option value="approved">Aprobadas</option>
                        <option value="rejected">Rechazadas</option>
                    </select>
                </div>
            </div>

            {error && <Alert type="error">{error}</Alert>}
            {loading && discrepancies.length === 0 ? <div className="flex justify-center p-8"><Spinner size="lg" /></div> : (
                <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                    <Table columns={columns} data={discrepancies.map(renderRow)} keyExtractor={(_, i) => String(discrepancies[i]?.id || i)} emptyMessage="Sin discrepancias" loading={loading} pagination={{ currentPage: page, totalPages, totalItems: total, itemsPerPage: pageSize, onPageChange: setPage, onItemsPerPageChange: setPageSize }} />
                </div>
            )}

            {reviewing && (
                <FormModal isOpen={true} onClose={() => setReviewing(null)} title={`${reviewing.action === 'approve' ? 'Aprobar' : 'Rechazar'} discrepancia #${reviewing.disc.id}`}>
                    <div className="p-6 space-y-4">
                        <p className="text-sm text-gray-600 dark:text-gray-300">
                            {reviewing.action === 'approve' ? 'Al aprobar se creará un movimiento de ajuste.' : 'Al rechazar se descarta la variación.'}
                        </p>
                        <div>
                            <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">{reviewing.action === 'approve' ? 'Notas' : 'Razón del rechazo'}</label>
                            <textarea value={notes} onChange={(e) => setNotes(e.target.value)} rows={3} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                        </div>
                        <div className="flex justify-end gap-2 pt-4 border-t border-gray-200 dark:border-gray-700">
                            <Button type="button" variant="outline" onClick={() => setReviewing(null)}>Cancelar</Button>
                            <button onClick={handleReview} className={`px-4 py-2 text-sm text-white rounded-md ${reviewing.action === 'approve' ? 'bg-green-600 hover:bg-green-700' : 'bg-red-600 hover:bg-red-700'}`}>{reviewing.action === 'approve' ? 'Aprobar' : 'Rechazar'}</button>
                        </div>
                    </div>
                </FormModal>
            )}
        </div>
    );
}

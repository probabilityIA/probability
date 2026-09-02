'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Modal } from '@/shared/ui';
import {
    Sprint,
    SprintStatus,
    SPRINT_STATUSES,
    SPRINT_STATUS_META,
    SPRINT_STATUS_RANK,
} from '@/services/modules/sprints/domain/types';
import {
    createSprintAction,
    updateSprintAction,
    changeSprintStatusAction,
} from '@/services/modules/sprints/infra/actions';
import { Ticket, ListTicketsParams, STATUS_META } from '../../domain/types';
import { listTicketsAction, changeTicketSprintAction } from '../../infra/actions';
import { PriorityBadge, TypeBadge } from './TicketBadges';

export const BACKLOG_PAGE_SIZE = 100;
const BACKLOG_KEY = 'backlog';
const DONE_STATUSES = ['resolved', 'closed', 'wont_fix'];

interface Bucket {
    tickets: Ticket[];
    total: number;
    loading: boolean;
    error: string | null;
}

type BucketMap = Record<string, Bucket>;

interface TicketsBacklogProps {
    sprints: Sprint[];
    sprintsLoading: boolean;
    sprintsError: string | null;
    canManage: boolean;
    baseParams: ListTicketsParams;
    onOpenTicket: (ticket: Ticket) => void;
    onSprintsChanged: () => void | Promise<void>;
    getAvatarUrl?: (ticket: Ticket) => string;
}

interface SprintFormState {
    name: string;
    goal: string;
    start_date: string;
    end_date: string;
    status: SprintStatus;
}

const emptyForm = (): SprintFormState => ({
    name: '',
    goal: '',
    start_date: '',
    end_date: '',
    status: 'planned',
});

const sprintKey = (id: number) => `sprint-${id}`;

const emptyBucket = (): Bucket => ({ tickets: [], total: 0, loading: true, error: null });

const formatDate = (value?: string) => {
    if (!value) return '';
    const d = new Date(value);
    if (Number.isNaN(d.getTime())) return value;
    return d.toLocaleDateString(undefined, { year: 'numeric', month: '2-digit', day: '2-digit' });
};

const toInputDate = (value?: string) => {
    if (!value) return '';
    const d = new Date(value);
    if (Number.isNaN(d.getTime())) return value.slice(0, 10);
    const pad = (n: number) => String(n).padStart(2, '0');
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
};

export default function TicketsBacklog({
    sprints,
    sprintsLoading,
    sprintsError,
    canManage,
    baseParams,
    onOpenTicket,
    onSprintsChanged,
    getAvatarUrl,
}: TicketsBacklogProps) {
    const [buckets, setBuckets] = useState<BucketMap>({});
    const [expanded, setExpanded] = useState<Record<string, boolean>>({});
    const [dragId, setDragId] = useState<number | null>(null);
    const [overKey, setOverKey] = useState<string | null>(null);
    const [movingId, setMovingId] = useState<number | null>(null);
    const [error, setError] = useState<string | null>(null);

    const [showForm, setShowForm] = useState(false);
    const [editing, setEditing] = useState<Sprint | null>(null);
    const [form, setForm] = useState<SprintFormState>(emptyForm());
    const [saving, setSaving] = useState(false);
    const [formError, setFormError] = useState<string | null>(null);
    const [statusBusyId, setStatusBusyId] = useState<number | null>(null);

    const paramsKey = useMemo(() => JSON.stringify(baseParams), [baseParams]);
    const paramsRef = useRef<ListTicketsParams>(baseParams);
    paramsRef.current = baseParams;

    const orderedSprints = useMemo(() => {
        return [...sprints].sort((a, b) => {
            const ra = SPRINT_STATUS_RANK[a.status] ?? 9;
            const rb = SPRINT_STATUS_RANK[b.status] ?? 9;
            if (ra !== rb) return ra - rb;
            const da = new Date(a.start_date || a.created_at).getTime();
            const db = new Date(b.start_date || b.created_at).getTime();
            if (a.status === 'closed') return db - da;
            return da - db;
        });
    }, [sprints]);

    const loadBucket = useCallback(async (key: string, sprintId: number | 'none') => {
        setBuckets((prev) => ({ ...prev, [key]: { ...(prev[key] || emptyBucket()), loading: true, error: null } }));
        try {
            const r = await listTicketsAction({
                ...paramsRef.current,
                sprint_id: sprintId,
                page: 1,
                page_size: BACKLOG_PAGE_SIZE,
            });
            const rows = (r?.data || []) as Ticket[];
            const list = sprintId === 'none' ? rows.filter((t) => !DONE_STATUSES.includes(t.status)) : rows;
            setBuckets((prev) => ({
                ...prev,
                [key]: { tickets: list, total: r?.total ?? list.length, loading: false, error: null },
            }));
        } catch (e) {
            console.error('Error listando tickets del backlog', e);
            setBuckets((prev) => ({
                ...prev,
                [key]: { tickets: [], total: 0, loading: false, error: 'No se pudieron cargar los tickets.' },
            }));
        }
    }, []);

    const firstRun = useRef(true);
    useEffect(() => {
        if (firstRun.current) {
            firstRun.current = false;
            return;
        }
        setBuckets({});
    }, [paramsKey]);

    useEffect(() => {
        setExpanded((prev) => {
            const next = { ...prev };
            let changed = false;
            orderedSprints.forEach((s) => {
                const k = sprintKey(s.id);
                if (next[k] === undefined) {
                    next[k] = s.status === 'active';
                    changed = true;
                }
            });
            return changed ? next : prev;
        });
    }, [orderedSprints]);

    useEffect(() => {
        if (!buckets[BACKLOG_KEY]) void loadBucket(BACKLOG_KEY, 'none');
        orderedSprints.forEach((s) => {
            const k = sprintKey(s.id);
            if (expanded[k] && !buckets[k]) void loadBucket(k, s.id);
        });
    }, [orderedSprints, expanded, buckets, loadBucket]);

    const toggleSprint = (id: number) => {
        const k = sprintKey(id);
        setExpanded((prev) => ({ ...prev, [k]: !prev[k] }));
    };

    const findTicket = (id: number): { key: string; ticket: Ticket } | null => {
        for (const [key, bucket] of Object.entries(buckets)) {
            const ticket = bucket.tickets.find((t) => t.id === id);
            if (ticket) return { key, ticket };
        }
        return null;
    };

    const handleDragStart = (e: React.DragEvent<HTMLDivElement>, ticket: Ticket) => {
        if (!canManage) return;
        setDragId(ticket.id);
        e.dataTransfer.effectAllowed = 'move';
        e.dataTransfer.setData('text/plain', String(ticket.id));
    };

    const handleDragEnd = () => {
        setDragId(null);
        setOverKey(null);
    };

    const handleDragOver = (e: React.DragEvent<HTMLDivElement>, key: string) => {
        if (!canManage) return;
        e.preventDefault();
        e.dataTransfer.dropEffect = 'move';
        if (overKey !== key) setOverKey(key);
    };

    const handleDragLeave = (e: React.DragEvent<HTMLDivElement>, key: string) => {
        if (e.currentTarget.contains(e.relatedTarget as Node)) return;
        if (overKey === key) setOverKey(null);
    };

    const handleDrop = async (e: React.DragEvent<HTMLDivElement>, targetKey: string, targetSprintId: number | null) => {
        e.preventDefault();
        setOverKey(null);
        if (!canManage) return;
        const raw = e.dataTransfer.getData('text/plain');
        const id = Number(raw) || dragId;
        setDragId(null);
        if (!id) return;

        const found = findTicket(id);
        if (!found) return;
        if (found.key === targetKey) return;

        const snapshot = buckets;
        const sprintName = targetSprintId ? (sprints.find((s) => s.id === targetSprintId)?.name || '') : '';
        const moved: Ticket = { ...found.ticket, sprint_id: targetSprintId, sprint_name: sprintName };

        setError(null);
        setMovingId(id);
        setBuckets((prev) => {
            const next: BucketMap = { ...prev };
            const source = next[found.key];
            if (source) {
                next[found.key] = {
                    ...source,
                    tickets: source.tickets.filter((t) => t.id !== id),
                    total: Math.max(0, source.total - 1),
                };
            }
            const target = next[targetKey];
            if (target) {
                next[targetKey] = { ...target, tickets: [moved, ...target.tickets], total: target.total + 1 };
            }
            return next;
        });

        try {
            const updated = await changeTicketSprintAction(id, targetSprintId);
            setBuckets((prev) => {
                const target = prev[targetKey];
                if (!target) return prev;
                return {
                    ...prev,
                    [targetKey]: {
                        ...target,
                        tickets: target.tickets.map((t) => (t.id === id ? (updated as Ticket) : t)),
                    },
                };
            });
            await onSprintsChanged();
        } catch (err) {
            console.error('cambio de sprint fallo', err);
            setBuckets(snapshot);
            setError('No se pudo mover el ticket de sprint. Intenta de nuevo.');
        } finally {
            setMovingId(null);
        }
    };

    const openCreate = () => {
        setEditing(null);
        setForm(emptyForm());
        setFormError(null);
        setShowForm(true);
    };

    const openEdit = (sprint: Sprint) => {
        setEditing(sprint);
        setForm({
            name: sprint.name || '',
            goal: sprint.goal || '',
            start_date: toInputDate(sprint.start_date),
            end_date: toInputDate(sprint.end_date),
            status: sprint.status || 'planned',
        });
        setFormError(null);
        setShowForm(true);
    };

    const submitForm = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!form.name.trim() || !form.start_date || !form.end_date) {
            setFormError('Nombre, fecha de inicio y fecha de fin son obligatorios.');
            return;
        }
        setSaving(true);
        setFormError(null);
        try {
            if (editing) {
                await updateSprintAction(editing.id, {
                    name: form.name.trim(),
                    goal: form.goal,
                    start_date: form.start_date,
                    end_date: form.end_date,
                    status: form.status,
                });
            } else {
                await createSprintAction({
                    name: form.name.trim(),
                    goal: form.goal || undefined,
                    start_date: form.start_date,
                    end_date: form.end_date,
                    status: form.status,
                });
            }
            setShowForm(false);
            await onSprintsChanged();
        } catch (err) {
            console.error('guardado de sprint fallo', err);
            setFormError('No se pudo guardar el sprint. Revisa los datos e intenta de nuevo.');
        } finally {
            setSaving(false);
        }
    };

    const changeStatus = async (sprint: Sprint, status: SprintStatus) => {
        setStatusBusyId(sprint.id);
        setError(null);
        try {
            await changeSprintStatusAction(sprint.id, status);
            await onSprintsChanged();
            setBuckets({});
        } catch (err) {
            console.error('cambio de estado de sprint fallo', err);
            setError('No se pudo cambiar el estado del sprint. Intenta de nuevo.');
        } finally {
            setStatusBusyId(null);
        }
    };

    const renderRow = (t: Ticket) => {
        const avatarUrl = getAvatarUrl ? getAvatarUrl(t) : '';
        const statusMeta = STATUS_META[t.status];
        const isMoving = movingId === t.id;
        return (
            <div
                key={t.id}
                draggable={canManage}
                onDragStart={(e) => handleDragStart(e, t)}
                onDragEnd={handleDragEnd}
                onClick={() => onOpenTicket(t)}
                className={`flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-3 py-2 transition-all hover:border-purple-300 hover:shadow-sm dark:border-gray-700 dark:bg-gray-800 ${
                    canManage ? 'cursor-grab active:cursor-grabbing' : 'cursor-pointer'
                } ${dragId === t.id ? 'opacity-40' : ''} ${isMoving ? 'ring-2 ring-purple-400 animate-pulse' : ''}`}
            >
                <span className="font-mono text-[11px] text-gray-500 dark:text-gray-400 w-20 shrink-0 truncate">{t.code}</span>
                <span className="min-w-0 flex-1 truncate text-sm text-gray-900 dark:text-gray-100">{t.title}</span>
                <span className="hidden shrink-0 sm:block"><TypeBadge type={t.type} /></span>
                <span className="shrink-0"><PriorityBadge priority={t.priority} /></span>
                <span className={`hidden shrink-0 rounded-full px-2 py-0.5 text-[11px] font-medium md:inline-flex ${statusMeta.bg} ${statusMeta.color}`}>
                    {statusMeta.label}
                </span>
                {avatarUrl ? (
                    <img src={avatarUrl} alt="" className="h-6 w-6 shrink-0 rounded-full object-cover ring-1 ring-gray-200 dark:ring-gray-600" />
                ) : (
                    <div className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-gray-200 text-[10px] text-gray-600 dark:bg-gray-600 dark:text-gray-300">
                        {t.assigned_to_name ? t.assigned_to_name[0].toUpperCase() : '-'}
                    </div>
                )}
            </div>
        );
    };

    const renderBucketBody = (key: string, emptyText: string) => {
        const bucket = buckets[key];
        if (!bucket || bucket.loading) {
            return (
                <div className="flex items-center justify-center gap-2 rounded-lg border border-dashed border-gray-300 p-4 text-xs text-gray-500 dark:border-gray-700 dark:text-gray-400">
                    <span className="h-3 w-3 animate-spin rounded-full border-2 border-purple-500 border-t-transparent"></span>
                    Cargando...
                </div>
            );
        }
        if (bucket.error) {
            return (
                <div className="rounded-lg border border-dashed border-red-300 p-4 text-center text-xs text-red-600 dark:border-red-800 dark:text-red-300">
                    {bucket.error}
                </div>
            );
        }
        if (bucket.tickets.length === 0) {
            return (
                <div className="rounded-lg border border-dashed border-gray-300 p-4 text-center text-xs text-gray-400 dark:border-gray-700 dark:text-gray-500">
                    {emptyText}
                </div>
            );
        }
        return (
            <>
                <div className="flex flex-col gap-1.5">{bucket.tickets.map(renderRow)}</div>
                {bucket.total > BACKLOG_PAGE_SIZE && (
                    <p className="mt-2 text-[11px] text-gray-500 dark:text-gray-400">
                        {`Mostrando los primeros ${BACKLOG_PAGE_SIZE} de ${bucket.total}. Usa los filtros para acotar la b\u00fasqueda.`}
                    </p>
                )}
            </>
        );
    };

    const backlogBucket = buckets[BACKLOG_KEY];

    return (
        <div className="flex flex-col gap-4">
            {error && (
                <div className="flex items-start justify-between gap-3 rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/30 dark:text-red-200">
                    <span>{error}</span>
                    <button type="button" onClick={() => setError(null)} className="text-xs font-semibold underline">
                        Cerrar
                    </button>
                </div>
            )}

            <div className="rounded-xl bg-white p-3 shadow-sm dark:bg-gray-800">
                <div className="mb-3 flex items-center justify-between gap-3">
                    <h2 className="text-sm font-bold uppercase tracking-wider text-gray-700 dark:text-gray-200">Sprints</h2>
                    {canManage && (
                        <button
                            type="button"
                            onClick={openCreate}
                            className="rounded-lg px-3 py-1.5 text-xs font-semibold text-white shadow-sm transition-opacity hover:opacity-90"
                            style={{ backgroundColor: 'var(--color-primary)' }}
                        >
                            + Nuevo sprint
                        </button>
                    )}
                </div>

                {sprintsError && (
                    <div className="rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-800 dark:border-amber-800 dark:bg-amber-900/30 dark:text-amber-200">
                        {sprintsError}
                    </div>
                )}

                {!sprintsError && sprintsLoading && (
                    <div className="rounded-lg border border-dashed border-gray-300 p-4 text-center text-xs text-gray-500 dark:border-gray-700 dark:text-gray-400">
                        Cargando sprints...
                    </div>
                )}

                {!sprintsError && !sprintsLoading && orderedSprints.length === 0 && (
                    <div className="rounded-lg border border-dashed border-gray-300 p-4 text-center text-xs text-gray-400 dark:border-gray-700 dark:text-gray-500">
                        {'No hay sprints creados todav\u00eda.'}
                    </div>
                )}

                <div className="flex flex-col gap-3">
                    {orderedSprints.map((s) => {
                        const key = sprintKey(s.id);
                        const isOpen = !!expanded[key];
                        const isOver = overKey === key;
                        const meta = SPRINT_STATUS_META[s.status] || SPRINT_STATUS_META.planned;
                        const busy = statusBusyId === s.id;
                        return (
                            <div
                                key={s.id}
                                onDragOver={(e) => handleDragOver(e, key)}
                                onDragLeave={(e) => handleDragLeave(e, key)}
                                onDrop={(e) => { void handleDrop(e, key, s.id); }}
                                className={`rounded-xl border transition-colors ${
                                    isOver
                                        ? 'border-purple-400 bg-purple-50 dark:border-purple-500 dark:bg-purple-900/20'
                                        : 'border-gray-200 bg-gray-50 dark:border-gray-700 dark:bg-gray-900/40'
                                }`}
                            >
                                <div className="flex flex-wrap items-center gap-2 border-b border-gray-200 px-3 py-2 dark:border-gray-700">
                                    <button
                                        type="button"
                                        onClick={() => toggleSprint(s.id)}
                                        className="flex min-w-0 flex-1 items-center gap-2 text-left"
                                    >
                                        <span className="text-xs text-gray-500 dark:text-gray-400">{isOpen ? 'v' : '>'}</span>
                                        <span className="truncate text-sm font-bold text-gray-900 dark:text-gray-100">{s.name}</span>
                                        <span className={`inline-flex shrink-0 items-center rounded-full px-2 py-0.5 text-[11px] font-semibold ring-1 ring-inset ${meta.bg} ${meta.color} ${meta.ring}`}>
                                            {meta.label}
                                        </span>
                                        <span className="hidden shrink-0 text-[11px] text-gray-500 dark:text-gray-400 sm:inline">
                                            {`${formatDate(s.start_date)} - ${formatDate(s.end_date)}`}
                                        </span>
                                    </button>
                                    <span className="shrink-0 rounded-full bg-gray-200 px-2 py-0.5 text-[11px] font-bold text-gray-700 dark:bg-gray-700 dark:text-gray-200">
                                        {`${s.done_count ?? 0} de ${s.ticket_count ?? 0} terminados`}
                                    </span>
                                    {canManage && (
                                        <div className="flex shrink-0 items-center gap-1">
                                            <button
                                                type="button"
                                                onClick={() => openEdit(s)}
                                                className="rounded-md border border-gray-300 px-2 py-1 text-[11px] font-semibold text-gray-700 hover:bg-gray-100 dark:border-gray-600 dark:text-gray-200 dark:hover:bg-gray-700"
                                            >
                                                Editar
                                            </button>
                                            {s.status !== 'active' && (
                                                <button
                                                    type="button"
                                                    disabled={busy}
                                                    onClick={() => { void changeStatus(s, 'active'); }}
                                                    className="rounded-md border border-emerald-300 px-2 py-1 text-[11px] font-semibold text-emerald-700 hover:bg-emerald-50 disabled:opacity-50 dark:border-emerald-700 dark:text-emerald-300 dark:hover:bg-emerald-900/30"
                                                >
                                                    Activar
                                                </button>
                                            )}
                                            {s.status !== 'closed' && (
                                                <button
                                                    type="button"
                                                    disabled={busy}
                                                    onClick={() => { void changeStatus(s, 'closed'); }}
                                                    className="rounded-md border border-gray-300 px-2 py-1 text-[11px] font-semibold text-gray-700 hover:bg-gray-100 disabled:opacity-50 dark:border-gray-600 dark:text-gray-200 dark:hover:bg-gray-700"
                                                >
                                                    Cerrar
                                                </button>
                                            )}
                                        </div>
                                    )}
                                </div>

                                {s.goal && (
                                    <p className="px-3 pt-2 text-[11px] text-gray-500 dark:text-gray-400">{`Objetivo: ${s.goal}`}</p>
                                )}

                                {isOpen && <div className="p-2">{renderBucketBody(key, 'Sin tickets en este sprint')}</div>}
                            </div>
                        );
                    })}
                </div>
            </div>

            <div
                onDragOver={(e) => handleDragOver(e, BACKLOG_KEY)}
                onDragLeave={(e) => handleDragLeave(e, BACKLOG_KEY)}
                onDrop={(e) => { void handleDrop(e, BACKLOG_KEY, null); }}
                className={`rounded-xl border p-3 shadow-sm transition-colors ${
                    overKey === BACKLOG_KEY
                        ? 'border-purple-400 bg-purple-50 dark:border-purple-500 dark:bg-purple-900/20'
                        : 'border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800'
                }`}
            >
                <div className="mb-3 flex items-center gap-2">
                    <h2 className="text-sm font-bold uppercase tracking-wider text-gray-700 dark:text-gray-200">Backlog</h2>
                    <span className="rounded-full bg-gray-200 px-2 py-0.5 text-[11px] font-bold text-gray-700 dark:bg-gray-700 dark:text-gray-200">
                        {backlogBucket?.tickets.length ?? 0}
                    </span>
                    <span className="text-[11px] text-gray-500 dark:text-gray-400">
                        {'Tickets sin sprint y sin terminar'}
                    </span>
                </div>
                {renderBucketBody(BACKLOG_KEY, 'El backlog est\u00e1 vac\u00edo')}
            </div>

            <Modal isOpen={showForm} onClose={() => setShowForm(false)} title={editing ? 'Editar sprint' : 'Nuevo sprint'} size="lg">
                <form onSubmit={submitForm} className="flex flex-col gap-3">
                    {formError && (
                        <div className="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-xs text-red-700 dark:border-red-800 dark:bg-red-900/30 dark:text-red-200">
                            {formError}
                        </div>
                    )}
                    <label className="flex flex-col gap-1 text-xs font-semibold text-gray-700 dark:text-gray-200">
                        Nombre
                        <input
                            type="text"
                            value={form.name}
                            onChange={(e) => setForm((p) => ({ ...p, name: e.target.value }))}
                            className="rounded-lg border border-gray-300 px-3 py-2 text-sm font-normal dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                            placeholder="Sprint 1"
                        />
                    </label>
                    <label className="flex flex-col gap-1 text-xs font-semibold text-gray-700 dark:text-gray-200">
                        Objetivo
                        <textarea
                            value={form.goal}
                            onChange={(e) => setForm((p) => ({ ...p, goal: e.target.value }))}
                            rows={2}
                            className="rounded-lg border border-gray-300 px-3 py-2 text-sm font-normal dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                        />
                    </label>
                    <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                        <label className="flex flex-col gap-1 text-xs font-semibold text-gray-700 dark:text-gray-200">
                            Fecha de inicio
                            <input
                                type="date"
                                value={form.start_date}
                                onChange={(e) => setForm((p) => ({ ...p, start_date: e.target.value }))}
                                className="rounded-lg border border-gray-300 px-3 py-2 text-sm font-normal dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                            />
                        </label>
                        <label className="flex flex-col gap-1 text-xs font-semibold text-gray-700 dark:text-gray-200">
                            Fecha de fin
                            <input
                                type="date"
                                value={form.end_date}
                                onChange={(e) => setForm((p) => ({ ...p, end_date: e.target.value }))}
                                className="rounded-lg border border-gray-300 px-3 py-2 text-sm font-normal dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                            />
                        </label>
                    </div>
                    <label className="flex flex-col gap-1 text-xs font-semibold text-gray-700 dark:text-gray-200">
                        Estado
                        <select
                            value={form.status}
                            onChange={(e) => setForm((p) => ({ ...p, status: e.target.value as SprintStatus }))}
                            className="rounded-lg border border-gray-300 px-3 py-2 text-sm font-normal dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                        >
                            {SPRINT_STATUSES.map((s) => (
                                <option key={s} value={s}>{SPRINT_STATUS_META[s].label}</option>
                            ))}
                        </select>
                    </label>
                    <div className="mt-2 flex justify-end gap-2">
                        <button
                            type="button"
                            onClick={() => setShowForm(false)}
                            className="rounded-lg border border-gray-300 px-4 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-100 dark:border-gray-600 dark:text-gray-200 dark:hover:bg-gray-700"
                        >
                            Cancelar
                        </button>
                        <button
                            type="submit"
                            disabled={saving}
                            className="rounded-lg px-4 py-2 text-sm font-semibold text-white shadow-sm transition-opacity hover:opacity-90 disabled:opacity-50"
                            style={{ backgroundColor: 'var(--color-primary)' }}
                        >
                            {saving ? 'Guardando...' : 'Guardar'}
                        </button>
                    </div>
                </form>
            </Modal>
        </div>
    );
}

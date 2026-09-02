'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { createPortal } from 'react-dom';
import { Modal, DynamicFilters, FilterOption, ActiveFilter, TablePagination, TICKETS_TABS_SLOT_ID, TICKETS_ACTIONS_SLOT_ID, TICKETS_FILTERS_SLOT_ID } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import {
    Ticket,
    TicketStatus,
    PaginatedTickets,
    ListTicketsParams,
    TICKET_STATUSES,
    TICKET_PRIORITIES,
    TICKET_TYPES,
    TICKET_AREAS,
    STATUS_META,
    PRIORITY_META,
    TYPE_META,
    AREA_META,
} from '../../domain/types';
import {
    listTicketsAction,
    createTicketAction,
    getTicketAction,
    changeTicketStatusAction,
    changeTicketAreaAction,
    assignTicketAction,
    uploadAttachmentAction,
    addCommentAction,
} from '../../infra/actions';
import { getUsersAction } from '@/services/auth/users/infra/actions';
import { Sprint } from '@/services/modules/sprints/domain/types';
import { listSprintsAction } from '@/services/modules/sprints/infra/actions';
import { PriorityBadge, TypeBadge } from './TicketBadges';
import TicketForm, { CreateTicketPayload } from './TicketForm';
import TicketDetail from './TicketDetail';
import TicketsBoard from './TicketsBoard';
import TicketsBacklog from './TicketsBacklog';

type ViewMode = 'table' | 'board' | 'backlog';

const VIEW_MODES: ViewMode[] = ['table', 'board', 'backlog'];
const VIEW_LABELS: Record<ViewMode, string> = { table: 'Tabla', board: 'Board', backlog: 'Backlog' };
const SPRINTS_PAGE_SIZE = 100;

const VIEW_STORAGE_KEY = 'tickets_view_mode';
const BOARD_PAGE_SIZE = 100;
const S3_BASE_URL = process.env.NEXT_PUBLIC_S3_BASE_URL || 'https://probability-media-assets.s3.us-east-1.amazonaws.com';

const resolveAvatarUrl = (avatarUrl?: string) => {
    if (!avatarUrl) return '';
    return avatarUrl.startsWith('http') ? avatarUrl : `${S3_BASE_URL}/${avatarUrl.replace(/^\//, '')}`;
};

const TABLE_COLUMNS = [
    'C\u00f3digo',
    'T\u00edtulo',
    'Tipo',
    '\u00c1rea',
    'Prioridad',
    'Estado',
    'Asignado',
    'Negocio',
    'Creado',
];

const TH_CLASS = 'px-3 sm:px-6 py-3 text-left text-xs font-bold text-white uppercase tracking-widest whitespace-nowrap';
const TH_STYLE = { paddingTop: '10px', paddingBottom: '10px', fontSize: '0.75rem', fontWeight: 800, letterSpacing: '0.06em' };
const TD_CLASS = 'px-3 sm:px-6 py-3';

export default function TicketsManager() {
    const { isSuperAdmin } = usePermissions();
    const [data, setData] = useState<PaginatedTickets | null>(null);
    const [loading, setLoading] = useState(false);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [filters, setFilters] = useState<{
        search?: string;
        status?: string;
        priority?: string;
        type?: string;
        area?: string;
        source?: string;
        only_mine?: boolean;
        escalated?: boolean;
        assigned_to_id?: number;
    }>({});
    const [sortBy, setSortBy] = useState<string>('created_at');
    const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');
    const [showCreate, setShowCreate] = useState(false);
    const [submitting, setSubmitting] = useState(false);
    const [createProgress, setCreateProgress] = useState('');
    const [openTicket, setOpenTicket] = useState<Ticket | null>(null);

    const syncTicketUrl = (t: Ticket | null) => {
        if (typeof window === 'undefined') return;
        const url = new URL(window.location.href);
        if (t) url.searchParams.set('ticket', String(t.id));
        else url.searchParams.delete('ticket');
        window.history.replaceState(null, '', url.pathname + url.search);
    };

    const openTicketDetail = (t: Ticket | null) => {
        setOpenTicket(t);
        syncTicketUrl(t);
    };
    const [users, setUsers] = useState<{ id: number; name: string; email: string; avatar_url?: string }[]>([]);
    const [updatingId, setUpdatingId] = useState<number | null>(null);
    const [view, setView] = useState<ViewMode>('table');
    const [error, setError] = useState<string | null>(null);
    const [sprints, setSprints] = useState<Sprint[]>([]);
    const [sprintsLoading, setSprintsLoading] = useState(true);
    const [sprintsError, setSprintsError] = useState<string | null>(null);
    const [boardSprintId, setBoardSprintId] = useState<string>('');
    const boardSprintDefaulted = useRef(false);
    const [tabsSlot, setTabsSlot] = useState<HTMLElement | null>(null);
    const [actionsSlot, setActionsSlot] = useState<HTMLElement | null>(null);
    const [filtersSlot, setFiltersSlot] = useState<HTMLElement | null>(null);

    useEffect(() => {
        setTabsSlot(document.getElementById(TICKETS_TABS_SLOT_ID));
        setActionsSlot(document.getElementById(TICKETS_ACTIONS_SLOT_ID));
        setFiltersSlot(document.getElementById(TICKETS_FILTERS_SLOT_ID));
    }, []);

    useEffect(() => {
        try {
            const stored = window.localStorage.getItem(VIEW_STORAGE_KEY) as ViewMode | null;
            if (stored && VIEW_MODES.includes(stored)) setView(stored);
        } catch {}
    }, []);

    const loadSprints = useCallback(async () => {
        setSprintsLoading(true);
        try {
            const r = await listSprintsAction({ page: 1, page_size: SPRINTS_PAGE_SIZE });
            setSprints((r?.data || []) as Sprint[]);
            setSprintsError(null);
        } catch (e) {
            console.error('Error listando sprints', e);
            setSprints([]);
            setSprintsError('No se pudieron cargar los sprints. Es posible que el m\u00f3dulo a\u00fan no est\u00e9 disponible.');
        } finally {
            setSprintsLoading(false);
        }
    }, []);

    useEffect(() => { void loadSprints(); }, [loadSprints]);

    useEffect(() => {
        if (boardSprintDefaulted.current || sprints.length === 0) return;
        boardSprintDefaulted.current = true;
        const active = sprints.find((s) => s.status === 'active');
        if (active) setBoardSprintId(String(active.id));
    }, [sprints]);

    const changeView = useCallback((next: ViewMode) => {
        setView(next);
        setPage(1);
        try {
            window.localStorage.setItem(VIEW_STORAGE_KEY, next);
        } catch {}
    }, []);

    useEffect(() => {
        if (!isSuperAdmin) return;
        (async () => {
            try {
                const r: any = await getUsersAction({ page: 1, page_size: 100 } as any);
                const list = (r?.data || []) as any[];
                setUsers(list.filter((u) => !!u.name && (u.scope_code === 'platform' || u.is_super_user)));
            } catch {}
        })();
    }, [isSuperAdmin]);

    const availableFilters: FilterOption[] = useMemo(() => [
        { key: 'search', label: 'Buscar', type: 'text', placeholder: 't\u00edtulo, c\u00f3digo, descripci\u00f3n...' },
        { key: 'status', label: 'Estado', type: 'select', options: TICKET_STATUSES.map(s => ({ value: s, label: STATUS_META[s].label })) },
        { key: 'area', label: '\u00c1rea', type: 'select', options: TICKET_AREAS.map(a => ({ value: a, label: AREA_META[a].label })) },
        { key: 'priority', label: 'Prioridad', type: 'select', options: TICKET_PRIORITIES.map(p => ({ value: p, label: PRIORITY_META[p].label })) },
        { key: 'type', label: 'Tipo', type: 'select', options: TICKET_TYPES.map(t => ({ value: t, label: TYPE_META[t].label })) },
        { key: 'source', label: 'Origen', type: 'select', options: [{ value: 'internal', label: 'Interno' }, { value: 'business', label: 'Negocio' }] },
        { key: 'only_mine', label: 'Solo m\u00edos', type: 'select', options: [{ value: 'true', label: 'Si' }, { value: 'false', label: 'No' }] },
        { key: 'escalated', label: 'Escalado a dev', type: 'select', options: [{ value: 'true', label: 'Si' }, { value: 'false', label: 'No' }] },
        ...(users.length > 0
            ? [{
                key: 'assigned_to_id',
                label: 'Asignado a',
                type: 'select' as const,
                avatarOptions: true,
                options: users.map((u) => ({
                    value: String(u.id),
                    label: u.name,
                    avatarUrl: resolveAvatarUrl(u.avatar_url),
                })),
            }]
            : []),
    ], [users]);

    const activeFilters: ActiveFilter[] = useMemo(() => {
        const out: ActiveFilter[] = [];
        if (filters.search) out.push({ key: 'search', label: 'Buscar', value: filters.search, type: 'text' });
        if (filters.status) out.push({ key: 'status', label: 'Estado', value: STATUS_META[filters.status as keyof typeof STATUS_META]?.label || filters.status, type: 'select' });
        if (filters.area) out.push({ key: 'area', label: '\u00c1rea', value: AREA_META[filters.area as keyof typeof AREA_META]?.label || filters.area, type: 'select' });
        if (filters.priority) out.push({ key: 'priority', label: 'Prioridad', value: PRIORITY_META[filters.priority as keyof typeof PRIORITY_META]?.label || filters.priority, type: 'select' });
        if (filters.type) out.push({ key: 'type', label: 'Tipo', value: TYPE_META[filters.type as keyof typeof TYPE_META]?.label || filters.type, type: 'select' });
        if (filters.source) out.push({ key: 'source', label: 'Origen', value: filters.source === 'internal' ? 'Interno' : 'Negocio', type: 'select' });
        if (filters.only_mine) out.push({ key: 'only_mine', label: 'Solo m\u00edos', value: 'Si', type: 'select' });
        if (filters.escalated) out.push({ key: 'escalated', label: 'Escalado a dev', value: 'Si', type: 'select' });
        if (filters.assigned_to_id) {
            const assigned = users.find((u) => u.id === filters.assigned_to_id);
            out.push({ key: 'assigned_to_id', label: 'Asignado a', value: assigned?.name || `Usuario ${filters.assigned_to_id}`, type: 'select' });
        }
        return out;
    }, [filters, users]);

    const handleAddFilter = useCallback((key: string, value: any) => {
        setPage(1);
        setFilters((prev) => {
            const f: any = { ...prev };
            if (key === 'only_mine' || key === 'escalated') {
                f[key] = value === 'true' || value === true;
            } else if (key === 'assigned_to_id') {
                const parsed = Number(value);
                if (!value || Number.isNaN(parsed)) delete f[key];
                else f[key] = parsed;
            } else {
                f[key] = value;
            }
            return f;
        });
    }, []);

    const handleRemoveFilter = useCallback((key: string) => {
        setPage(1);
        setFilters((prev) => {
            const f: any = { ...prev };
            delete f[key];
            return f;
        });
    }, []);

    const formatDateTime = (iso: string) => {
        try {
            return new Date(iso).toLocaleString(undefined, {
                year: 'numeric', month: '2-digit', day: '2-digit',
                hour: '2-digit', minute: '2-digit',
            });
        } catch {
            return iso;
        }
    };

    const updateLocalTicket = (updated: Ticket) => {
        setData((prev) => prev ? { ...prev, data: prev.data.map(t => t.id === updated.id ? updated : t) } : prev);
    };

    const patchLocalStatus = (id: number, status: TicketStatus) => {
        setData((prev) => prev ? { ...prev, data: prev.data.map(t => t.id === id ? { ...t, status } : t) } : prev);
    };

    const avatarUrlFor = useCallback((t: Ticket) => {
        const assignedUser = users.find(u => u.id === t.assigned_to_id);
        return resolveAvatarUrl(t.assigned_to_avatar_url || assignedUser?.avatar_url);
    }, [users]);

    const handleBoardMove = async (id: number, status: TicketStatus) => {
        const current = data?.data.find(t => t.id === id);
        if (!current || current.status === status) return;
        const previousStatus = current.status;
        setError(null);
        setUpdatingId(id);
        patchLocalStatus(id, status);
        try {
            const updated = await changeTicketStatusAction(id, status);
            updateLocalTicket(updated as Ticket);
        } catch (e) {
            patchLocalStatus(id, previousStatus);
            console.error('cambio de estado fallo', e);
            setError('No se pudo cambiar el estado del ticket. Intenta de nuevo.');
        } finally {
            setUpdatingId(null);
        }
    };

    const handleStatusChange = async (id: number, status: string) => {
        setUpdatingId(id);
        try {
            const updated = await changeTicketStatusAction(id, status);
            updateLocalTicket(updated as Ticket);
        } catch (e) {
            console.error('cambio de estado fallo', e);
        } finally {
            setUpdatingId(null);
        }
    };

    const handleAreaChange = async (id: number, area: string) => {
        setUpdatingId(id);
        try {
            const updated = await changeTicketAreaAction(id, area);
            updateLocalTicket(updated as Ticket);
        } catch (e) {
            console.error('cambio de area fallo', e);
        } finally {
            setUpdatingId(null);
        }
    };

    const handleAssignChange = async (id: number, val: string) => {
        setUpdatingId(id);
        try {
            const userId = val === '' ? null : Number(val);
            const updated = await assignTicketAction(id, userId);
            updateLocalTicket(updated as Ticket);
        } catch (e) {
            console.error('cambio de asignado fallo', e);
        } finally {
            setUpdatingId(null);
        }
    };

    const listParams = useMemo<ListTicketsParams>(() => ({
        search: filters.search || undefined,
        status: filters.status || undefined,
        priority: filters.priority || undefined,
        type: filters.type || undefined,
        area: filters.area || undefined,
        source: filters.source || undefined,
        only_mine: filters.only_mine || undefined,
        escalated: filters.escalated || undefined,
        assigned_to_id: filters.assigned_to_id || undefined,
        sort_by: sortBy,
        sort_order: sortOrder,
    }), [filters, sortBy, sortOrder]);

    const boardSprintParam = useMemo(() => {
        if (!boardSprintId) return undefined;
        if (boardSprintId === 'none') return 'none';
        const parsed = Number(boardSprintId);
        return Number.isNaN(parsed) ? undefined : parsed;
    }, [boardSprintId]);

    const fetchData = useCallback(async () => {
        if (view === 'backlog') return;
        setLoading(true);
        try {
            const r = await listTicketsAction({
                ...listParams,
                page: view === 'board' ? 1 : page,
                page_size: view === 'board' ? BOARD_PAGE_SIZE : pageSize,
                sprint_id: view === 'board' ? boardSprintParam : undefined,
            });
            setData(r);
        } catch (e) {
            console.error('Error listando tickets', e);
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, listParams, view, boardSprintParam]);

    useEffect(() => { fetchData(); }, [fetchData]);

    useEffect(() => {
        const raw = new URLSearchParams(window.location.search).get('ticket');
        const id = raw ? Number(raw) : NaN;
        if (!raw || Number.isNaN(id)) return;
        (async () => {
            try {
                const t = await getTicketAction(id);
                if (t) setOpenTicket(t);
                else setError('No se encontr\u00f3 el ticket ' + id + '.');
            } catch (e: any) {
                const msg = String(e?.message || '');
                if (msg.includes('403') || msg.toLowerCase().includes('super')) {
                    setError('No tienes permiso para ver este ticket.');
                } else if (msg.includes('404')) {
                    setError('No se encontr\u00f3 el ticket ' + id + '.');
                } else {
                    setError('No se pudo abrir el ticket ' + id + '.');
                }
                syncTicketUrl(null);
            }
        })();
    }, []);

    const handleCreate = async ({ dto, files, comment, commentInternal }: CreateTicketPayload) => {
        setSubmitting(true);
        setCreateProgress('Creando...');
        setError(null);
        let created: Ticket;
        try {
            created = await createTicketAction(dto) as Ticket;
        } catch (e) {
            setSubmitting(false);
            setCreateProgress('');
            throw e;
        }

        const label = created?.code || '#' + created?.id;
        const problems: string[] = [];
        let failedFiles = 0;

        for (let i = 0; i < files.length; i++) {
            setCreateProgress('Subiendo adjuntos ' + (i + 1) + ' de ' + files.length + '...');
            try {
                const fd = new FormData();
                fd.append('file', files[i]);
                await uploadAttachmentAction(created.id, fd);
            } catch (e) {
                console.error('fallo al subir adjunto', files[i]?.name, e);
                failedFiles++;
            }
        }
        if (failedFiles > 0) {
            problems.push('no se pudieron subir ' + failedFiles + (failedFiles === 1 ? ' adjunto' : ' adjuntos'));
        }

        if (comment) {
            setCreateProgress('Guardando comentario...');
            try {
                await addCommentAction(created.id, comment, commentInternal);
            } catch (e) {
                console.error('fallo al agregar comentario inicial', e);
                problems.push('no se pudo agregar el comentario inicial');
            }
        }

        setShowCreate(false);
        setSubmitting(false);
        setCreateProgress('');
        if (problems.length > 0) {
            setError('Ticket ' + label + ' creado, pero ' + problems.join(' y ') + '.');
        }
        await fetchData();
    };

    const refreshOpenTicket = async () => {
        if (openTicket) {
            try {
                const fresh = await getTicketAction(openTicket.id);
                setOpenTicket(fresh);
            } catch {}
        }
        fetchData();
    };

    const sprintSelector = (
        <select
            value={boardSprintId}
            onChange={(e) => { setBoardSprintId(e.target.value); setPage(1); }}
            className="rounded-lg border border-gray-300 bg-white px-2 py-1.5 text-xs font-medium text-gray-700 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-200"
        >
            <option value="">Todos los sprints</option>
            <option value="none">Sin sprint</option>
            {sprints.map((s) => (
                <option key={s.id} value={String(s.id)}>{s.name}</option>
            ))}
        </select>
    );

    const viewToggle = (
        <div className="inline-flex items-center rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 p-0.5">
            {VIEW_MODES.map((v) => (
                <button
                    key={v}
                    type="button"
                    onClick={() => changeView(v)}
                    className={`px-2.5 py-1 text-xs font-semibold rounded-md transition-colors ${
                        view === v
                            ? 'text-white shadow-sm'
                            : 'text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white'
                    }`}
                    style={view === v ? { backgroundColor: 'var(--color-primary)' } : undefined}
                >
                    {VIEW_LABELS[v]}
                </button>
            ))}
        </div>
    );

    const headerActions = (
        <>
            {view === 'board' && sprintSelector}
        </>
    );

    const headerFilters = (
        <DynamicFilters
            variant="bar"
            availableFilters={availableFilters}
            activeFilters={activeFilters}
            onAddFilter={handleAddFilter}
            onRemoveFilter={handleRemoveFilter}
            onCreate={() => setShowCreate(true)}
            createButtonIconOnly
            createButtonAriaLabel="Nuevo ticket"
            createButtonPosition="right"
            sortBy={sortBy}
            sortOrder={sortOrder}
            onSortChange={(by, order) => { setSortBy(by); setSortOrder(order); setPage(1); }}
            sortOptions={[
                { value: 'created_at', label: 'Por fecha' },
                { value: 'updated_at', label: 'Por actualizaci\u00f3n' },
                { value: 'priority', label: 'Por prioridad' },
                { value: 'status', label: 'Por estado' },
                { value: 'area', label: 'Por \u00e1rea' },
                { value: 'code', label: 'Por c\u00f3digo' },
                { value: 'due_date', label: 'Por fecha l\u00edmite' },
            ]}
        />
    );

    const headerInSlots = !!(tabsSlot || actionsSlot || filtersSlot);

    return (
        <div>
            {tabsSlot && createPortal(viewToggle, tabsSlot)}
            {actionsSlot && createPortal(headerActions, actionsSlot)}
            {filtersSlot && createPortal(headerFilters, filtersSlot)}

            {!headerInSlots && (
                <div className="mb-4 flex flex-col gap-2">
                    <div className="flex items-center gap-2 flex-wrap">
                        {viewToggle}
                        {headerActions}
                    </div>
                    {headerFilters}
                </div>
            )}

            {error && (
                <div className="mb-3 flex items-start justify-between gap-3 rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/30 dark:text-red-200">
                    <span>{error}</span>
                    <button type="button" onClick={() => setError(null)} className="text-xs font-semibold underline">
                        Cerrar
                    </button>
                </div>
            )}

            {view === 'backlog' ? (
                <TicketsBacklog
                    sprints={sprints}
                    sprintsLoading={sprintsLoading}
                    sprintsError={sprintsError}
                    canManage={isSuperAdmin}
                    baseParams={listParams}
                    onOpenTicket={openTicketDetail}
                    onSprintsChanged={loadSprints}
                    getAvatarUrl={avatarUrlFor}
                />
            ) : view === 'board' ? (
                <div>
                    {data && data.total > BOARD_PAGE_SIZE && (
                        <p className="mb-2 text-xs text-gray-500 dark:text-gray-400">
                            {`Mostrando los primeros ${BOARD_PAGE_SIZE} tickets de ${data.total}. Usa los filtros para acotar la b\u00fasqueda.`}
                        </p>
                    )}
                    <TicketsBoard
                        tickets={data?.data || []}
                        loading={loading}
                        canDrag={isSuperAdmin}
                        updatingId={updatingId}
                        onOpen={openTicketDetail}
                        onMove={handleBoardMove}
                        getAvatarUrl={avatarUrlFor}
                    />
                </div>
            ) : (
            <div className="relative rounded-xl overflow-hidden shadow-sm bg-white dark:bg-gray-800">
                {loading && (
                    <div className="absolute inset-0 bg-white dark:bg-gray-800/80 backdrop-blur-sm z-10 flex items-center justify-center transition-opacity duration-200">
                        <div className="flex flex-col items-center gap-2">
                            <div className="w-8 h-8 border-4 border-purple-500 border-t-transparent rounded-full animate-spin"></div>
                            <p className="text-sm text-gray-600 dark:text-gray-300">Actualizando...</p>
                        </div>
                    </div>
                )}

                <div className="overflow-x-auto">
                    <table className={`min-w-full transition-opacity duration-200 ${loading ? 'opacity-50' : 'opacity-100'}`}>
                        <thead style={{ backgroundColor: 'var(--color-primary)', color: 'var(--color-on-primary, white)' }}>
                            <tr>
                                {TABLE_COLUMNS.map((label) => (
                                    <th key={label} className={TH_CLASS} style={TH_STYLE}>{label}</th>
                                ))}
                            </tr>
                        </thead>
                        <tbody>
                            {!loading && data && data.data.length === 0 && (
                                <tr>
                                    <td colSpan={TABLE_COLUMNS.length} className="px-4 sm:px-6 py-8 text-center text-gray-500 dark:text-gray-400">
                                        No hay tickets disponibles
                                    </td>
                                </tr>
                            )}
                            {data?.data.map((t) => {
                                const isUpdating = updatingId === t.id;
                                const fullAvatarUrl = avatarUrlFor(t);
                                const stop = (e: React.MouseEvent | React.ChangeEvent) => { e.stopPropagation(); };
                                const areaMeta = t.area ? AREA_META[t.area] : null;
                                const statusMeta = STATUS_META[t.status];
                                return (
                                <tr
                                    key={t.id}
                                    className="bg-white dark:bg-gray-800 border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors cursor-pointer"
                                    onClick={() => openTicketDetail(t)}
                                >
                                    <td className={`${TD_CLASS} font-mono text-xs text-gray-600 dark:text-gray-300`}>{t.code}</td>
                                    <td className={`${TD_CLASS} text-gray-900 dark:text-gray-100`}>
                                        <div className="font-medium truncate max-w-md">{t.title}</div>
                                        <div className="text-xs text-gray-500 dark:text-gray-400">{t.comments_count} comentarios | {t.attachments_count} adjuntos</div>
                                    </td>
                                    <td className={TD_CLASS}><TypeBadge type={t.type} /></td>
                                    <td className={TD_CLASS} onClick={stop}>
                                        <select
                                            value={t.area || 'soporte'}
                                            disabled={isUpdating || !isSuperAdmin}
                                            onChange={(e) => handleAreaChange(t.id, e.target.value)}
                                            className={`text-xs font-medium rounded-full px-2 py-1 border-0 cursor-pointer focus:ring-2 focus:ring-offset-1 ${areaMeta?.bg || 'bg-gray-100'} ${areaMeta?.color || 'text-gray-700'} disabled:opacity-60`}
                                        >
                                            {TICKET_AREAS.map((a) => <option key={a} value={a}>{AREA_META[a].label}</option>)}
                                        </select>
                                    </td>
                                    <td className={TD_CLASS}><PriorityBadge priority={t.priority} /></td>
                                    <td className={TD_CLASS} onClick={stop}>
                                        <select
                                            value={t.status}
                                            disabled={isUpdating || !isSuperAdmin}
                                            onChange={(e) => handleStatusChange(t.id, e.target.value)}
                                            className={`text-xs font-medium rounded-full px-2 py-1 border-0 cursor-pointer focus:ring-2 focus:ring-offset-1 ${statusMeta.bg} ${statusMeta.color} disabled:opacity-60`}
                                        >
                                            {TICKET_STATUSES.map((s) => <option key={s} value={s}>{STATUS_META[s].label}</option>)}
                                        </select>
                                    </td>
                                    <td className={TD_CLASS} onClick={stop}>
                                        {isSuperAdmin ? (
                                            <div className="flex items-center gap-2">
                                                {fullAvatarUrl ? (
                                                    <img src={fullAvatarUrl} alt="" className="h-6 w-6 rounded-full object-cover ring-1 ring-gray-200 dark:ring-gray-600" />
                                                ) : (
                                                    <div className="h-6 w-6 rounded-full bg-gray-200 dark:bg-gray-600 flex items-center justify-center text-[10px] text-gray-600 dark:text-gray-300">
                                                        {t.assigned_to_name ? t.assigned_to_name[0].toUpperCase() : '-'}
                                                    </div>
                                                )}
                                                <select
                                                    value={t.assigned_to_id ?? ''}
                                                    disabled={isUpdating}
                                                    onChange={(e) => handleAssignChange(t.id, e.target.value)}
                                                    className="text-xs rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-2 py-1 max-w-[140px] disabled:opacity-60"
                                                >
                                                    <option value="">Sin asignar</option>
                                                    {t.assigned_to_id && !users.some(u => u.id === t.assigned_to_id) && (
                                                        <option value={t.assigned_to_id}>{t.assigned_to_name || `Usuario ${t.assigned_to_id}`}</option>
                                                    )}
                                                    {users.map((u) => <option key={u.id} value={u.id}>{u.name}</option>)}
                                                </select>
                                            </div>
                                        ) : (
                                            <div className="flex items-center gap-2">
                                                {fullAvatarUrl && (
                                                    <img src={fullAvatarUrl} alt="" className="h-6 w-6 rounded-full object-cover ring-1 ring-gray-200 dark:ring-gray-600" />
                                                )}
                                                <span className="text-xs text-gray-700 dark:text-gray-300">{t.assigned_to_name || '-'}</span>
                                            </div>
                                        )}
                                    </td>
                                    <td className={`${TD_CLASS} text-xs text-gray-700 dark:text-gray-300`}>{t.business_name || (t.business_id ? `#${t.business_id}` : 'Interno')}</td>
                                    <td className={`${TD_CLASS} text-xs text-gray-500 dark:text-gray-400 whitespace-nowrap`}>{formatDateTime(t.created_at)}</td>
                                </tr>
                                );
                            })}
                        </tbody>
                    </table>
                </div>

                {data && (
                    <TablePagination
                        currentPage={data.page || page}
                        totalPages={data.total_pages || 1}
                        totalItems={data.total || 0}
                        pageSize={pageSize}
                        onPageChange={setPage}
                        onPageSizeChange={(size) => { setPageSize(size); setPage(1); }}
                    />
                )}
            </div>
            )}

            <Modal isOpen={showCreate} onClose={() => setShowCreate(false)} title="Nuevo ticket" size="wide">
                <TicketForm
                    isSuperAdmin={isSuperAdmin}
                    users={users}
                    sprints={sprints}
                    onSubmit={handleCreate}
                    onCancel={() => setShowCreate(false)}
                    submitting={submitting}
                    submitLabel={createProgress}
                />
            </Modal>

            <Modal isOpen={!!openTicket} onClose={() => openTicketDetail(null)} title="" size="wide" noPadding noBodyScroll>
                {openTicket && (
                    <TicketDetail
                        ticket={openTicket}
                        isSuperAdmin={isSuperAdmin}
                        onClose={() => openTicketDetail(null)}
                        onChanged={refreshOpenTicket}
                    />
                )}
            </Modal>
        </div>
    );
}

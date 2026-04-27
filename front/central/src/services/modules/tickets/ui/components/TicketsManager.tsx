'use client';

import { useCallback, useEffect, useMemo, useState } from 'react';
import { Button, Modal, SuperAdminBusinessSelector, DynamicFilters, FilterOption, ActiveFilter } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import {
    Ticket,
    PaginatedTickets,
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
} from '../../infra/actions';
import { getUsersAction } from '@/services/auth/users/infra/actions';
import { PriorityBadge, TypeBadge } from './TicketBadges';
import TicketForm from './TicketForm';
import TicketDetail from './TicketDetail';

interface Props {
    selectedBusinessId?: number | null;
    onBusinessChange?: (id: number | null) => void;
}

export default function TicketsManager({ selectedBusinessId = null, onBusinessChange }: Props) {
    const { isSuperAdmin } = usePermissions();
    const [data, setData] = useState<PaginatedTickets | null>(null);
    const [loading, setLoading] = useState(false);
    const [page, setPage] = useState(1);
    const [pageSize] = useState(10);
    const [filters, setFilters] = useState<{
        search?: string;
        status?: string;
        priority?: string;
        type?: string;
        area?: string;
        source?: string;
        only_mine?: boolean;
        escalated?: boolean;
    }>({});
    const [sortBy, setSortBy] = useState<string>('created_at');
    const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');
    const [showCreate, setShowCreate] = useState(false);
    const [submitting, setSubmitting] = useState(false);
    const [openTicket, setOpenTicket] = useState<Ticket | null>(null);
    const [users, setUsers] = useState<{ id: number; name: string; email: string; avatar_url?: string }[]>([]);
    const [updatingId, setUpdatingId] = useState<number | null>(null);

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
        { key: 'search', label: 'Buscar', type: 'text', placeholder: 'titulo, codigo, descripcion...' },
        { key: 'status', label: 'Estado', type: 'select', options: TICKET_STATUSES.map(s => ({ value: s, label: STATUS_META[s].label })) },
        { key: 'area', label: 'Area', type: 'select', options: TICKET_AREAS.map(a => ({ value: a, label: AREA_META[a].label })) },
        { key: 'priority', label: 'Prioridad', type: 'select', options: TICKET_PRIORITIES.map(p => ({ value: p, label: PRIORITY_META[p].label })) },
        { key: 'type', label: 'Tipo', type: 'select', options: TICKET_TYPES.map(t => ({ value: t, label: TYPE_META[t].label })) },
        { key: 'source', label: 'Origen', type: 'select', options: [{ value: 'internal', label: 'Interno' }, { value: 'business', label: 'Negocio' }] },
        { key: 'only_mine', label: 'Solo mios', type: 'select', options: [{ value: 'true', label: 'Si' }, { value: 'false', label: 'No' }] },
        { key: 'escalated', label: 'Escalado a dev', type: 'select', options: [{ value: 'true', label: 'Si' }, { value: 'false', label: 'No' }] },
    ], []);

    const activeFilters: ActiveFilter[] = useMemo(() => {
        const out: ActiveFilter[] = [];
        if (filters.search) out.push({ key: 'search', label: 'Buscar', value: filters.search, type: 'text' });
        if (filters.status) out.push({ key: 'status', label: 'Estado', value: STATUS_META[filters.status as keyof typeof STATUS_META]?.label || filters.status, type: 'select' });
        if (filters.area) out.push({ key: 'area', label: 'Area', value: AREA_META[filters.area as keyof typeof AREA_META]?.label || filters.area, type: 'select' });
        if (filters.priority) out.push({ key: 'priority', label: 'Prioridad', value: PRIORITY_META[filters.priority as keyof typeof PRIORITY_META]?.label || filters.priority, type: 'select' });
        if (filters.type) out.push({ key: 'type', label: 'Tipo', value: TYPE_META[filters.type as keyof typeof TYPE_META]?.label || filters.type, type: 'select' });
        if (filters.source) out.push({ key: 'source', label: 'Origen', value: filters.source === 'internal' ? 'Interno' : 'Negocio', type: 'select' });
        if (filters.only_mine) out.push({ key: 'only_mine', label: 'Solo mios', value: 'Si', type: 'select' });
        if (filters.escalated) out.push({ key: 'escalated', label: 'Escalado a dev', value: 'Si', type: 'select' });
        return out;
    }, [filters]);

    const handleAddFilter = useCallback((key: string, value: any) => {
        setPage(1);
        setFilters((prev) => {
            const f: any = { ...prev };
            if (key === 'only_mine' || key === 'escalated') {
                f[key] = value === 'true' || value === true;
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

    const fetchData = useCallback(async () => {
        setLoading(true);
        try {
            const r = await listTicketsAction({
                page,
                page_size: pageSize,
                business_id: isSuperAdmin && selectedBusinessId ? selectedBusinessId : undefined,
                search: filters.search || undefined,
                status: filters.status || undefined,
                priority: filters.priority || undefined,
                type: filters.type || undefined,
                area: filters.area || undefined,
                source: filters.source || undefined,
                only_mine: filters.only_mine || undefined,
                escalated: filters.escalated || undefined,
                sort_by: sortBy,
                sort_order: sortOrder,
            });
            setData(r);
        } catch (e) {
            console.error('Error listando tickets', e);
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, isSuperAdmin, selectedBusinessId, filters, sortBy, sortOrder]);

    useEffect(() => { fetchData(); }, [fetchData]);

    const handleCreate = async (dto: any) => {
        setSubmitting(true);
        try {
            await createTicketAction(dto);
            setShowCreate(false);
            await fetchData();
        } finally {
            setSubmitting(false);
        }
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

    return (
        <div className="space-y-4">
            <div className="flex items-center justify-between flex-wrap gap-3">
                <div>
                    <h1 className="text-2xl font-semibold text-gray-900 dark:text-gray-100">Tickets</h1>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Soporte interno y de negocios. Bugs, mejoras, integraciones, datos.</p>
                </div>
                <div className="flex items-center gap-2">
                    {isSuperAdmin && onBusinessChange && (
                        <SuperAdminBusinessSelector value={selectedBusinessId} onChange={onBusinessChange} placeholder="Todos los negocios" />
                    )}
                    <Button variant="primary" onClick={() => setShowCreate(true)}>+ Nuevo ticket</Button>
                </div>
            </div>

            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-4">
                <DynamicFilters
                    availableFilters={availableFilters}
                    activeFilters={activeFilters}
                    onAddFilter={handleAddFilter}
                    onRemoveFilter={handleRemoveFilter}
                    sortBy={sortBy}
                    sortOrder={sortOrder}
                    onSortChange={(by, order) => { setSortBy(by); setSortOrder(order); setPage(1); }}
                    sortOptions={[
                        { value: 'created_at', label: 'Ordenar por fecha de creacion' },
                        { value: 'updated_at', label: 'Ordenar por ultima actualizacion' },
                        { value: 'priority', label: 'Ordenar por prioridad' },
                        { value: 'status', label: 'Ordenar por estado' },
                        { value: 'area', label: 'Ordenar por area' },
                        { value: 'code', label: 'Ordenar por codigo' },
                        { value: 'due_date', label: 'Ordenar por fecha limite' },
                    ]}
                    className="!p-0 !border-0 !shadow-none"
                />
            </div>

            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700 text-sm">
                        <thead className="bg-gray-50 dark:bg-gray-900">
                            <tr className="text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                                <th className="px-4 py-3">Codigo</th>
                                <th className="px-4 py-3">Titulo</th>
                                <th className="px-4 py-3">Tipo</th>
                                <th className="px-4 py-3">Area</th>
                                <th className="px-4 py-3">Prioridad</th>
                                <th className="px-4 py-3">Estado</th>
                                <th className="px-4 py-3">Asignado</th>
                                <th className="px-4 py-3">Negocio</th>
                                <th className="px-4 py-3">Creado</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                            {loading && (
                                <tr><td colSpan={9} className="px-4 py-6 text-center text-gray-500">Cargando...</td></tr>
                            )}
                            {!loading && data && data.data.length === 0 && (
                                <tr><td colSpan={9} className="px-4 py-6 text-center text-gray-500">Sin tickets</td></tr>
                            )}
                            {!loading && data?.data.map((t) => {
                                const isUpdating = updatingId === t.id;
                                const assignedUser = users.find(u => u.id === t.assigned_to_id);
                                const avatarUrl = t.assigned_to_avatar_url || assignedUser?.avatar_url || '';
                                const fullAvatarUrl = avatarUrl && !avatarUrl.startsWith('http')
                                    ? `${process.env.NEXT_PUBLIC_S3_BASE_URL || 'https://probability-media-assets.s3.us-east-1.amazonaws.com'}/${avatarUrl.replace(/^\//, '')}`
                                    : avatarUrl;
                                const stop = (e: React.MouseEvent | React.ChangeEvent) => { e.stopPropagation(); };
                                const areaMeta = t.area ? AREA_META[t.area] : null;
                                const statusMeta = STATUS_META[t.status];
                                return (
                                <tr key={t.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/50 cursor-pointer" onClick={() => setOpenTicket(t)}>
                                    <td className="px-4 py-3 font-mono text-xs text-gray-600 dark:text-gray-300">{t.code}</td>
                                    <td className="px-4 py-3 text-gray-900 dark:text-gray-100">
                                        <div className="font-medium truncate max-w-md">{t.title}</div>
                                        <div className="text-xs text-gray-500">{t.comments_count} comentarios | {t.attachments_count} adjuntos</div>
                                    </td>
                                    <td className="px-4 py-3"><TypeBadge type={t.type} /></td>
                                    <td className="px-4 py-3" onClick={stop}>
                                        <select
                                            value={t.area || 'soporte'}
                                            disabled={isUpdating || !isSuperAdmin}
                                            onChange={(e) => handleAreaChange(t.id, e.target.value)}
                                            className={`text-xs font-medium rounded-full px-2 py-1 border-0 cursor-pointer focus:ring-2 focus:ring-offset-1 ${areaMeta?.bg || 'bg-gray-100'} ${areaMeta?.color || 'text-gray-700'} disabled:opacity-60`}
                                        >
                                            {TICKET_AREAS.map((a) => <option key={a} value={a}>{AREA_META[a].label}</option>)}
                                        </select>
                                    </td>
                                    <td className="px-4 py-3"><PriorityBadge priority={t.priority} /></td>
                                    <td className="px-4 py-3" onClick={stop}>
                                        <select
                                            value={t.status}
                                            disabled={isUpdating || !isSuperAdmin}
                                            onChange={(e) => handleStatusChange(t.id, e.target.value)}
                                            className={`text-xs font-medium rounded-full px-2 py-1 border-0 cursor-pointer focus:ring-2 focus:ring-offset-1 ${statusMeta.bg} ${statusMeta.color} disabled:opacity-60`}
                                        >
                                            {TICKET_STATUSES.map((s) => <option key={s} value={s}>{STATUS_META[s].label}</option>)}
                                        </select>
                                    </td>
                                    <td className="px-4 py-3" onClick={stop}>
                                        {isSuperAdmin ? (
                                            <div className="flex items-center gap-2">
                                                {fullAvatarUrl ? (
                                                    /* eslint-disable-next-line @next/next/no-img-element */
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
                                                    /* eslint-disable-next-line @next/next/no-img-element */
                                                    <img src={fullAvatarUrl} alt="" className="h-6 w-6 rounded-full object-cover ring-1 ring-gray-200 dark:ring-gray-600" />
                                                )}
                                                <span className="text-xs text-gray-700 dark:text-gray-300">{t.assigned_to_name || '-'}</span>
                                            </div>
                                        )}
                                    </td>
                                    <td className="px-4 py-3 text-xs">{t.business_name || (t.business_id ? `#${t.business_id}` : 'Interno')}</td>
                                    <td className="px-4 py-3 text-xs text-gray-500 whitespace-nowrap">{formatDateTime(t.created_at)}</td>
                                </tr>
                                );
                            })}
                        </tbody>
                    </table>
                </div>
                {data && data.total_pages > 1 && (
                    <div className="flex items-center justify-between px-4 py-3 border-t border-gray-200 dark:border-gray-700 text-sm">
                        <span className="text-gray-500">Pagina {data.page} de {data.total_pages} ({data.total} totales)</span>
                        <div className="flex gap-2">
                            <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage(page - 1)}>Anterior</Button>
                            <Button variant="outline" size="sm" disabled={page >= data.total_pages} onClick={() => setPage(page + 1)}>Siguiente</Button>
                        </div>
                    </div>
                )}
            </div>

            <Modal isOpen={showCreate} onClose={() => setShowCreate(false)} title="Nuevo ticket" size="lg">
                <TicketForm
                    isSuperAdmin={isSuperAdmin}
                    selectedBusinessId={selectedBusinessId}
                    onSubmit={handleCreate}
                    onCancel={() => setShowCreate(false)}
                    submitting={submitting}
                />
            </Modal>

            <Modal isOpen={!!openTicket} onClose={() => setOpenTicket(null)} title="" size="xl">
                {openTicket && (
                    <TicketDetail
                        ticket={openTicket}
                        isSuperAdmin={isSuperAdmin}
                        onClose={() => setOpenTicket(null)}
                        onChanged={refreshOpenTicket}
                    />
                )}
            </Modal>
        </div>
    );
}

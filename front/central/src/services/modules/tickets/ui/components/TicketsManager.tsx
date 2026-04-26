'use client';

import { useCallback, useEffect, useState } from 'react';
import { Button, Modal, Input, SuperAdminBusinessSelector } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import {
    Ticket,
    PaginatedTickets,
    TICKET_STATUSES,
    TICKET_PRIORITIES,
    TICKET_TYPES,
    STATUS_META,
    PRIORITY_META,
    TYPE_META,
} from '../../domain/types';
import { listTicketsAction, createTicketAction, getTicketAction } from '../../infra/actions';
import { StatusBadge, PriorityBadge, TypeBadge } from './TicketBadges';
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
    const [search, setSearch] = useState('');
    const [statusFilter, setStatusFilter] = useState<string>('');
    const [priorityFilter, setPriorityFilter] = useState<string>('');
    const [typeFilter, setTypeFilter] = useState<string>('');
    const [showCreate, setShowCreate] = useState(false);
    const [submitting, setSubmitting] = useState(false);
    const [openTicket, setOpenTicket] = useState<Ticket | null>(null);

    const fetchData = useCallback(async () => {
        setLoading(true);
        try {
            const r = await listTicketsAction({
                page,
                page_size: pageSize,
                business_id: isSuperAdmin && selectedBusinessId ? selectedBusinessId : undefined,
                search: search || undefined,
                status: statusFilter || undefined,
                priority: priorityFilter || undefined,
                type: typeFilter || undefined,
            });
            setData(r);
        } catch (e) {
            console.error('Error listando tickets', e);
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, isSuperAdmin, selectedBusinessId, search, statusFilter, priorityFilter, typeFilter]);

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

            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-4 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
                <Input placeholder="Buscar por titulo, codigo, descripcion..." value={search} onChange={(e) => { setSearch(e.target.value); setPage(1); }} />
                <select value={statusFilter} onChange={(e) => { setStatusFilter(e.target.value); setPage(1); }} className="rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm px-3 py-2">
                    <option value="">Todos los estados</option>
                    {TICKET_STATUSES.map((s) => <option key={s} value={s}>{STATUS_META[s].label}</option>)}
                </select>
                <select value={priorityFilter} onChange={(e) => { setPriorityFilter(e.target.value); setPage(1); }} className="rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm px-3 py-2">
                    <option value="">Todas las prioridades</option>
                    {TICKET_PRIORITIES.map((p) => <option key={p} value={p}>{PRIORITY_META[p].label}</option>)}
                </select>
                <select value={typeFilter} onChange={(e) => { setTypeFilter(e.target.value); setPage(1); }} className="rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm px-3 py-2">
                    <option value="">Todos los tipos</option>
                    {TICKET_TYPES.map((t) => <option key={t} value={t}>{TYPE_META[t].label}</option>)}
                </select>
            </div>

            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700 text-sm">
                        <thead className="bg-gray-50 dark:bg-gray-900">
                            <tr className="text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                                <th className="px-4 py-3">Codigo</th>
                                <th className="px-4 py-3">Titulo</th>
                                <th className="px-4 py-3">Tipo</th>
                                <th className="px-4 py-3">Prioridad</th>
                                <th className="px-4 py-3">Estado</th>
                                <th className="px-4 py-3">Asignado</th>
                                <th className="px-4 py-3">Negocio</th>
                                <th className="px-4 py-3">Creado</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                            {loading && (
                                <tr><td colSpan={8} className="px-4 py-6 text-center text-gray-500">Cargando...</td></tr>
                            )}
                            {!loading && data && data.data.length === 0 && (
                                <tr><td colSpan={8} className="px-4 py-6 text-center text-gray-500">Sin tickets</td></tr>
                            )}
                            {!loading && data?.data.map((t) => (
                                <tr key={t.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/50 cursor-pointer" onClick={() => setOpenTicket(t)}>
                                    <td className="px-4 py-3 font-mono text-xs text-gray-600 dark:text-gray-300">{t.code}</td>
                                    <td className="px-4 py-3 text-gray-900 dark:text-gray-100">
                                        <div className="font-medium truncate max-w-md">{t.title}</div>
                                        <div className="text-xs text-gray-500">{t.comments_count} comentarios | {t.attachments_count} adjuntos</div>
                                    </td>
                                    <td className="px-4 py-3"><TypeBadge type={t.type} /></td>
                                    <td className="px-4 py-3"><PriorityBadge priority={t.priority} /></td>
                                    <td className="px-4 py-3"><StatusBadge status={t.status} /></td>
                                    <td className="px-4 py-3 text-xs">{t.assigned_to_name || '-'}</td>
                                    <td className="px-4 py-3 text-xs">{t.business_name || (t.business_id ? `#${t.business_id}` : 'Interno')}</td>
                                    <td className="px-4 py-3 text-xs text-gray-500">{new Date(t.created_at).toLocaleDateString()}</td>
                                </tr>
                            ))}
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

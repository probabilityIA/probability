'use client';

import { useEffect, useRef, useState } from 'react';
import { Button } from '@/shared/ui';
import {
    Ticket,
    TicketComment,
    TicketAttachment,
    TicketHistoryEntry,
    TICKET_STATUSES,
    STATUS_META,
    PRIORITY_META,
    TYPE_META,
    TicketStatus,
} from '../../domain/types';
import {
    listCommentsAction,
    addCommentAction,
    listAttachmentsAction,
    uploadAttachmentAction,
    deleteAttachmentAction,
    listTicketHistoryAction,
    changeTicketStatusAction,
    escalateTicketAction,
    deleteTicketAction,
    assignTicketAction,
} from '../../infra/actions';
import { getUsersAction } from '@/services/auth/users/infra/actions';
import { StatusBadge, PriorityBadge, TypeBadge } from './TicketBadges';

interface Props {
    ticket: Ticket;
    isSuperAdmin: boolean;
    onClose: () => void;
    onChanged: () => void;
}

export default function TicketDetail({ ticket, isSuperAdmin, onClose, onChanged }: Props) {
    const [comments, setComments] = useState<TicketComment[]>([]);
    const [attachments, setAttachments] = useState<TicketAttachment[]>([]);
    const [history, setHistory] = useState<TicketHistoryEntry[]>([]);
    const [loading, setLoading] = useState(true);
    const [newComment, setNewComment] = useState('');
    const [internalComment, setInternalComment] = useState(false);
    const [posting, setPosting] = useState(false);
    const [uploading, setUploading] = useState(false);
    const [statusNote, setStatusNote] = useState('');
    const [users, setUsers] = useState<{ id: number; name: string; email: string }[]>([]);
    const [assigning, setAssigning] = useState(false);
    const fileRef = useRef<HTMLInputElement>(null);

    const refreshAll = async () => {
        setLoading(true);
        try {
            const [c, a, h] = await Promise.all([
                listCommentsAction(ticket.id),
                listAttachmentsAction(ticket.id),
                listTicketHistoryAction(ticket.id),
            ]);
            setComments(c);
            setAttachments(a);
            setHistory(h);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => { refreshAll(); }, [ticket.id]);

    useEffect(() => {
        if (!isSuperAdmin) return;
        (async () => {
            try {
                const r: any = await getUsersAction({ page: 1, page_size: 100 } as any);
                const list = (r?.data || []) as Array<{ id: number; name: string; email: string; scope_code?: string; is_super_user?: boolean }>;
                setUsers(list.filter((u) => !!u.name && (u.scope_code === 'platform' || u.is_super_user)));
            } catch {}
        })();
    }, [isSuperAdmin]);

    const handleAssign = async (val: string) => {
        setAssigning(true);
        try {
            const id = val === '' ? null : Number(val);
            await assignTicketAction(ticket.id, id);
            await refreshAll();
            onChanged();
        } finally {
            setAssigning(false);
        }
    };

    const submitComment = async () => {
        if (!newComment.trim()) return;
        setPosting(true);
        try {
            await addCommentAction(ticket.id, newComment.trim(), isSuperAdmin && internalComment);
            setNewComment('');
            setInternalComment(false);
            await refreshAll();
            onChanged();
        } finally {
            setPosting(false);
        }
    };

    const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;
        setUploading(true);
        try {
            const fd = new FormData();
            fd.append('file', file);
            await uploadAttachmentAction(ticket.id, fd);
            await refreshAll();
            onChanged();
        } finally {
            setUploading(false);
            if (fileRef.current) fileRef.current.value = '';
        }
    };

    const handleStatusChange = async (newStatus: TicketStatus) => {
        if (newStatus === ticket.status) return;
        await changeTicketStatusAction(ticket.id, newStatus, statusNote);
        setStatusNote('');
        await refreshAll();
        onChanged();
    };

    const handleEscalate = async () => {
        await escalateTicketAction(ticket.id, 'Escalado a desarrollo');
        await refreshAll();
        onChanged();
    };

    const handleDelete = async () => {
        if (!confirm('Eliminar este ticket de forma definitiva?')) return;
        await deleteTicketAction(ticket.id);
        onClose();
        onChanged();
    };

    const isImage = (mime: string) => mime?.startsWith('image/');

    return (
        <div className="space-y-6">
            <div className="flex items-start justify-between gap-4">
                <div className="space-y-2">
                    <div className="flex items-center gap-2 flex-wrap">
                        <span className="font-mono text-sm text-gray-500 dark:text-gray-400">{ticket.code}</span>
                        <StatusBadge status={ticket.status} />
                        <PriorityBadge priority={ticket.priority} />
                        <TypeBadge type={ticket.type} />
                        {ticket.escalated_to_dev && (
                            <span className="text-xs px-2 py-0.5 rounded bg-fuchsia-100 dark:bg-fuchsia-900/40 text-fuchsia-700 dark:text-fuchsia-200">Escalado a dev</span>
                        )}
                    </div>
                    <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">{ticket.title}</h2>
                    <div className="text-xs text-gray-500 dark:text-gray-400 flex flex-wrap gap-x-4 gap-y-1">
                        <span>Creado por: {ticket.created_by_name || `#${ticket.created_by_id}`}</span>
                        {ticket.business_name && <span>Negocio: {ticket.business_name}</span>}
                        {ticket.assigned_to_name && <span>Asignado a: {ticket.assigned_to_name}</span>}
                        {ticket.category && <span>Categoria: {ticket.category}</span>}
                    </div>
                </div>
                {isSuperAdmin && (
                    <div className="flex flex-col gap-2 items-end">
                        {!ticket.escalated_to_dev && (
                            <Button variant="purple" size="sm" onClick={handleEscalate}>Escalar a dev</Button>
                        )}
                        <Button variant="danger" size="sm" onClick={handleDelete}>Eliminar</Button>
                    </div>
                )}
            </div>

            <div className="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 text-sm text-gray-700 dark:text-gray-200 whitespace-pre-wrap">
                {ticket.description}
            </div>

            {isSuperAdmin && (
                <div className="border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-2">
                    <div className="text-sm font-semibold text-gray-700 dark:text-gray-200">Asignar a</div>
                    <select
                        value={ticket.assigned_to_id ?? ''}
                        onChange={(e) => handleAssign(e.target.value)}
                        disabled={assigning}
                        className="block w-full sm:w-80 rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm px-3 py-2"
                    >
                        <option value="">Sin asignar</option>
                        {users.map((u) => (
                            <option key={u.id} value={u.id}>{u.name} ({u.email})</option>
                        ))}
                    </select>
                </div>
            )}

            <div className="border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-3">
                <div className="text-sm font-semibold text-gray-700 dark:text-gray-200">Cambiar estado</div>
                <div className="flex flex-wrap gap-2">
                    {TICKET_STATUSES.map((s) => {
                        const m = STATUS_META[s];
                        const active = s === ticket.status;
                        return (
                            <button
                                key={s}
                                type="button"
                                onClick={() => handleStatusChange(s)}
                                className={`px-2.5 py-1 rounded-full text-xs font-semibold transition ${m.bg} ${m.color} ${active ? 'ring-2 ring-offset-1 ring-blue-500' : 'opacity-70 hover:opacity-100'}`}
                            >
                                {m.label}
                            </button>
                        );
                    })}
                </div>
                <input
                    type="text"
                    value={statusNote}
                    onChange={(e) => setStatusNote(e.target.value)}
                    placeholder="Nota opcional al cambiar de estado"
                    className="block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm px-3 py-1.5"
                />
            </div>

            <div className="border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-3">
                <div className="flex items-center justify-between">
                    <div className="text-sm font-semibold text-gray-700 dark:text-gray-200">Adjuntos ({attachments.length})</div>
                    <div>
                        <input ref={fileRef} type="file" onChange={handleUpload} className="hidden" id="ticket-file-input" />
                        <label htmlFor="ticket-file-input" className="inline-flex items-center px-3 py-1.5 rounded-md text-xs font-medium bg-blue-600 hover:bg-blue-700 text-white cursor-pointer">
                            {uploading ? 'Subiendo...' : 'Subir archivo'}
                        </label>
                    </div>
                </div>
                {attachments.length === 0 ? (
                    <div className="text-xs text-gray-500 dark:text-gray-400">Sin adjuntos</div>
                ) : (
                    <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
                        {attachments.map((a) => (
                            <div key={a.id} className="border border-gray-200 dark:border-gray-700 rounded-lg p-2 text-xs">
                                {isImage(a.mime_type) ? (
                                    <a href={a.file_url} target="_blank" rel="noreferrer">
                                        <img src={a.file_url} alt={a.file_name} className="w-full h-24 object-cover rounded" />
                                    </a>
                                ) : (
                                    <a href={a.file_url} target="_blank" rel="noreferrer" className="block w-full h-24 rounded bg-gray-100 dark:bg-gray-800 flex items-center justify-center text-gray-500">
                                        Archivo
                                    </a>
                                )}
                                <div className="mt-1 truncate" title={a.file_name}>{a.file_name}</div>
                                <div className="flex justify-between items-center mt-1">
                                    <span className="text-[10px] text-gray-500">{(a.size / 1024).toFixed(1)} KB</span>
                                    {(isSuperAdmin) && (
                                        <button
                                            onClick={async () => {
                                                if (confirm('Eliminar este adjunto?')) {
                                                    await deleteAttachmentAction(a.id);
                                                    refreshAll();
                                                }
                                            }}
                                            className="text-red-600 hover:text-red-700 text-[10px]"
                                        >
                                            Eliminar
                                        </button>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>

            <div className="border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-4">
                <div className="text-sm font-semibold text-gray-700 dark:text-gray-200">Comentarios ({comments.length})</div>
                {loading && <div className="text-xs text-gray-500">Cargando...</div>}
                <div className="space-y-3 max-h-80 overflow-y-auto">
                    {comments.map((c) => (
                        <div key={c.id} className={`p-3 rounded-lg text-sm ${c.is_internal ? 'bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800' : 'bg-gray-50 dark:bg-gray-800/40'}`}>
                            <div className="flex items-center justify-between text-xs text-gray-500 dark:text-gray-400 mb-1">
                                <span className="font-medium text-gray-700 dark:text-gray-200">{c.user_name || `Usuario ${c.user_id}`}</span>
                                <span>{new Date(c.created_at).toLocaleString()}</span>
                            </div>
                            <div className="whitespace-pre-wrap text-gray-800 dark:text-gray-200">{c.body}</div>
                            {c.is_internal && <div className="mt-1 text-[10px] text-amber-700 dark:text-amber-300">Nota interna</div>}
                        </div>
                    ))}
                    {!loading && comments.length === 0 && (
                        <div className="text-xs text-gray-500 dark:text-gray-400">Sin comentarios todavia</div>
                    )}
                </div>
                <div className="space-y-2">
                    <textarea
                        value={newComment}
                        onChange={(e) => setNewComment(e.target.value)}
                        rows={3}
                        placeholder="Escribe un comentario..."
                        className="block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm px-3 py-2"
                    />
                    <div className="flex items-center justify-between">
                        {isSuperAdmin ? (
                            <label className="flex items-center gap-2 text-xs text-gray-600 dark:text-gray-300">
                                <input type="checkbox" checked={internalComment} onChange={(e) => setInternalComment(e.target.checked)} />
                                Nota interna (solo super admins)
                            </label>
                        ) : <span />}
                        <Button variant="primary" size="sm" onClick={submitComment} disabled={posting || !newComment.trim()}>
                            {posting ? 'Enviando...' : 'Comentar'}
                        </Button>
                    </div>
                </div>
            </div>

            <div className="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
                <div className="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">Historial</div>
                {history.length === 0 ? (
                    <div className="text-xs text-gray-500">Sin movimientos</div>
                ) : (
                    <ul className="space-y-2 text-xs text-gray-600 dark:text-gray-300">
                        {history.map((h) => (
                            <li key={h.id} className="flex items-start gap-2">
                                <span className="text-gray-400">{new Date(h.created_at).toLocaleString()}</span>
                                <span>
                                    {h.changed_by_name || `Usuario ${h.changed_by_id}`}
                                    {h.from_status && h.from_status !== h.to_status ? ` cambio ${h.from_status} -> ${h.to_status}` : ` actualizo`}
                                    {h.note ? ` (${h.note})` : ''}
                                </span>
                            </li>
                        ))}
                    </ul>
                )}
            </div>

            <div className="flex justify-end pt-2">
                <Button variant="outline" onClick={onClose}>Cerrar</Button>
            </div>
        </div>
    );
}

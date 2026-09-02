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
import { ACCEPTED_TYPES, attachmentError, formatSize, isImageMime, isPdfMime } from './attachment-rules';

interface Props {
    ticket: Ticket;
    isSuperAdmin: boolean;
    onClose: () => void;
    onChanged: () => void;
}

const formatDate = (value?: string) => (value ? new Date(value).toLocaleString() : '');

const initials = (name?: string) => {
    if (!name) return '?';
    const parts = name.trim().split(/\s+/);
    const first = parts[0]?.[0] || '';
    const second = parts.length > 1 ? parts[parts.length - 1][0] : '';
    return (first + second).toUpperCase();
};

const sectionTitle = 'text-[11px] font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400';
const panel = 'rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/60 p-4';

export default function TicketDetail({ ticket, isSuperAdmin, onClose, onChanged }: Props) {
    const [comments, setComments] = useState<TicketComment[]>([]);
    const [attachments, setAttachments] = useState<TicketAttachment[]>([]);
    const [history, setHistory] = useState<TicketHistoryEntry[]>([]);
    const [loading, setLoading] = useState(true);
    const [newComment, setNewComment] = useState('');
    const [internalComment, setInternalComment] = useState(false);
    const [posting, setPosting] = useState(false);
    const [uploading, setUploading] = useState(false);
    const [uploadingName, setUploadingName] = useState('');
    const [uploadError, setUploadError] = useState('');
    const [dragOver, setDragOver] = useState(false);
    const [statusNote, setStatusNote] = useState('');
    const [users, setUsers] = useState<{ id: number; name: string; email: string }[]>([]);
    const [assigning, setAssigning] = useState(false);
    const [copied, setCopied] = useState(false);
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

    const uploadFile = async (file: File) => {
        setUploadError('');
        const invalid = attachmentError(file);
        if (invalid) {
            setUploadError(invalid);
            return;
        }
        setUploading(true);
        setUploadingName(file.name);
        try {
            const fd = new FormData();
            fd.append('file', file);
            await uploadAttachmentAction(ticket.id, fd);
            await refreshAll();
            onChanged();
        } catch {
            setUploadError('No se pudo subir ' + file.name + '. Intenta de nuevo.');
        } finally {
            setUploading(false);
            setUploadingName('');
            if (fileRef.current) fileRef.current.value = '';
        }
    };

    const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) await uploadFile(file);
    };

    const handleDrop = async (e: React.DragEvent) => {
        e.preventDefault();
        setDragOver(false);
        const file = e.dataTransfer.files?.[0];
        if (file) await uploadFile(file);
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
        if (!confirm('\u00bfEliminar este ticket de forma definitiva?')) return;
        await deleteTicketAction(ticket.id);
        onClose();
        onChanged();
    };

    const copyLink = async () => {
        try {
            const url = window.location.origin + window.location.pathname + '?ticket=' + ticket.id;
            await navigator.clipboard.writeText(url);
            setCopied(true);
            setTimeout(() => setCopied(false), 2000);
        } catch {}
    };

    const removeAttachment = async (id: number) => {
        if (!confirm('\u00bfEliminar este adjunto?')) return;
        await deleteAttachmentAction(id);
        await refreshAll();
    };

    return (
        <div className="flex flex-col gap-5">

            <div className="flex items-start justify-between gap-4 pb-4 border-b border-gray-200 dark:border-gray-700">
                <div className="min-w-0 flex flex-col gap-2">
                    <div className="flex items-center gap-2 flex-wrap">
                        <span className="font-mono text-sm text-gray-500 dark:text-gray-400">{ticket.code}</span>
                        <StatusBadge status={ticket.status} />
                        <PriorityBadge priority={ticket.priority} />
                        <TypeBadge type={ticket.type} />
                        {ticket.escalated_to_dev && (
                            <span className="text-xs px-2 py-0.5 rounded-full bg-fuchsia-100 dark:bg-fuchsia-900/40 text-fuchsia-700 dark:text-fuchsia-200">Escalado a dev</span>
                        )}
                    </div>
                    <h2 className="text-xl sm:text-2xl font-semibold text-gray-900 dark:text-gray-100 break-words">{ticket.title}</h2>
                    <div className="text-xs text-gray-500 dark:text-gray-400 flex flex-wrap items-center gap-x-3 gap-y-1">
                        <span>Creado por {ticket.created_by_name || '#' + ticket.created_by_id}</span>
                        <span className="text-gray-300 dark:text-gray-600">|</span>
                        <span>{ticket.business_name || 'Interno'}</span>
                        {ticket.category && (
                            <>
                                <span className="text-gray-300 dark:text-gray-600">|</span>
                                <span>{'Categor\u00eda: '}{ticket.category}</span>
                            </>
                        )}
                        <span className="text-gray-300 dark:text-gray-600">|</span>
                        <span>{formatDate(ticket.created_at)}</span>
                    </div>
                </div>
                <div className="flex items-center gap-2 shrink-0">
                    <button
                        type="button"
                        onClick={copyLink}
                        title="Copiar enlace a este ticket"
                        className="inline-flex items-center gap-1.5 h-8 px-3 rounded-lg border border-gray-300 dark:border-gray-600 text-xs font-medium text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 transition"
                    >
                        <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M13.19 8.688a4.5 4.5 0 011.242 7.244l-4.5 4.5a4.5 4.5 0 01-6.364-6.364l1.757-1.757m13.35-.622l1.757-1.757a4.5 4.5 0 00-6.364-6.364l-4.5 4.5a4.5 4.5 0 001.242 7.244" />
                        </svg>
                        {copied ? 'Copiado' : 'Enlace'}
                    </button>
                    {isSuperAdmin && !ticket.escalated_to_dev && (
                        <Button variant="purple" size="sm" onClick={handleEscalate}>Escalar a dev</Button>
                    )}
                    {isSuperAdmin && (
                        <button
                            type="button"
                            onClick={handleDelete}
                            title="Eliminar ticket"
                            className="inline-flex items-center justify-center w-8 h-8 rounded-lg border border-gray-300 dark:border-gray-600 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/30 transition"
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                            </svg>
                        </button>
                    )}
                </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-[minmax(0,1fr)_340px] gap-6">

                <div className="min-w-0 flex flex-col gap-6">

                    <div className="flex flex-col gap-2">
                        <div className={sectionTitle}>{'Descripci\u00f3n'}</div>
                        <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/60 p-4 text-sm leading-relaxed text-gray-700 dark:text-gray-200 whitespace-pre-wrap break-words">
                            {ticket.description}
                        </div>
                    </div>

                    <div className="flex flex-col gap-2">
                        <div className="flex items-center justify-between">
                            <div className={sectionTitle}>Adjuntos <span className="text-gray-400 dark:text-gray-500">({attachments.length})</span></div>
                            <span className="text-[11px] text-gray-400 dark:text-gray-500">{'Im\u00e1genes y PDF | hasta 10 MB'}</span>
                        </div>

                        {uploadError && (
                            <div className="flex items-start gap-2 rounded-xl border border-red-200 dark:border-red-800 bg-red-50 dark:bg-red-900/20 px-3 py-2">
                                <svg className="w-4 h-4 mt-0.5 shrink-0 text-red-500" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z" />
                                </svg>
                                <span className="flex-1 text-xs text-red-700 dark:text-red-300">{uploadError}</span>
                                <button type="button" onClick={() => setUploadError('')} className="text-red-500 hover:text-red-700 text-xs">x</button>
                            </div>
                        )}

                        {uploading && (
                            <div className="flex items-center gap-3 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/60 px-4 py-3">
                                <svg className="w-4 h-4 animate-spin text-blue-500" fill="none" viewBox="0 0 24 24">
                                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                                </svg>
                                <span className="text-sm text-gray-600 dark:text-gray-300 truncate">{'Subiendo ' + uploadingName}</span>
                            </div>
                        )}

                        <div className="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-3">
                            {attachments.map((a) => (
                                <div key={a.id} className="rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden bg-white dark:bg-gray-900/60">
                                    <a href={a.file_url} target="_blank" rel="noreferrer" className="block h-24 bg-gray-100 dark:bg-gray-800">
                                        {isImageMime(a.mime_type) ? (
                                            <img src={a.file_url} alt={a.file_name} className="w-full h-24 object-cover" />
                                        ) : (
                                            <span className={`w-full h-24 flex flex-col items-center justify-center gap-1 ${isPdfMime(a.mime_type) ? 'bg-red-50 dark:bg-red-900/20' : ''}`}>
                                                <svg className={`w-6 h-6 ${isPdfMime(a.mime_type) ? 'text-red-400' : 'text-gray-400'}`} fill="none" stroke="currentColor" strokeWidth={1.5} viewBox="0 0 24 24">
                                                    <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25M9 16.5v.75m3-3v3M15 12v5.25m-4.5-15H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
                                                </svg>
                                                {isPdfMime(a.mime_type) && <span className="font-mono text-[10px] font-medium text-red-400 tracking-wide">PDF</span>}
                                            </span>
                                        )}
                                    </a>
                                    <div className="p-2 flex flex-col gap-1">
                                        <div className="text-xs text-gray-700 dark:text-gray-200 truncate" title={a.file_name}>{a.file_name}</div>
                                        <div className="flex items-center justify-between">
                                            <span className="text-[10px] text-gray-400 dark:text-gray-500">{formatSize(a.size)}</span>
                                            {isSuperAdmin && (
                                                <button type="button" onClick={() => removeAttachment(a.id)} className="text-[10px] text-red-500 hover:text-red-600">Eliminar</button>
                                            )}
                                        </div>
                                    </div>
                                </div>
                            ))}

                            <label
                                htmlFor="ticket-file-input"
                                onDragOver={(e) => { e.preventDefault(); setDragOver(true); }}
                                onDragLeave={() => setDragOver(false)}
                                onDrop={handleDrop}
                                className={`col-span-2 min-h-[138px] flex flex-col items-center justify-center gap-1.5 rounded-xl border-[1.5px] border-dashed p-4 cursor-pointer transition ${dragOver ? 'border-sky-400 bg-sky-50 dark:bg-sky-900/20' : 'border-gray-300 dark:border-gray-600 bg-gray-50/60 dark:bg-gray-900/40 hover:border-gray-400 dark:hover:border-gray-500'}`}
                            >
                                <svg className={`w-6 h-6 ${dragOver ? 'text-sky-400' : 'text-gray-400 dark:text-gray-500'}`} fill="none" stroke="currentColor" strokeWidth={1.5} viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M12 16.5V9.75m0 0l3 3m-3-3l-3 3M6.75 19.5a4.5 4.5 0 01-1.41-8.775 5.25 5.25 0 0110.233-2.33 3 3 0 013.758 3.848A3.752 3.752 0 0118 19.5H6.75z" />
                                </svg>
                                <span className="text-sm text-gray-600 dark:text-gray-300">
                                    {dragOver ? 'Suelta para adjuntar' : 'Arrastra un archivo aqu\u00ed'}
                                </span>
                                {!dragOver && (
                                    <span className="text-[11px] text-gray-400 dark:text-gray-500">
                                        {'o '}<span className="text-sky-600 dark:text-sky-400 font-medium">{'selecci\u00f3nalo'}</span>{' desde tu equipo'}
                                    </span>
                                )}
                                <input ref={fileRef} type="file" accept={ACCEPTED_TYPES} onChange={handleUpload} className="hidden" id="ticket-file-input" />
                            </label>
                        </div>
                    </div>

                    <div className="flex flex-col gap-3">
                        <div className={sectionTitle}>Comentarios <span className="text-gray-400 dark:text-gray-500">({comments.length})</span></div>

                        <div className="flex flex-col gap-2.5 max-h-96 overflow-y-auto">
                            {comments.map((c) => (
                                <div key={c.id} className="flex gap-2.5">
                                    <div className={`w-8 h-8 shrink-0 rounded-full flex items-center justify-center text-[11px] font-semibold ${c.is_internal ? 'bg-amber-100 dark:bg-amber-900/50 text-amber-700 dark:text-amber-200' : 'bg-sky-100 dark:bg-sky-900/50 text-sky-700 dark:text-sky-200'}`}>
                                        {initials(c.user_name)}
                                    </div>
                                    <div className={`flex-1 min-w-0 rounded-xl border px-3 py-2.5 ${c.is_internal ? 'bg-amber-50 dark:bg-amber-900/20 border-amber-200 dark:border-amber-800' : 'bg-gray-50 dark:bg-gray-900/60 border-gray-200 dark:border-gray-700'}`}>
                                        <div className="flex items-center justify-between gap-2 mb-1">
                                            <div className="flex items-center gap-2 min-w-0">
                                                <span className="text-xs font-semibold text-gray-700 dark:text-gray-200 truncate">{c.user_name || 'Usuario ' + c.user_id}</span>
                                                {c.is_internal && (
                                                    <span className="text-[10px] font-semibold px-1.5 py-0.5 rounded-full bg-amber-200 dark:bg-amber-800/60 text-amber-800 dark:text-amber-200 shrink-0">Nota interna</span>
                                                )}
                                            </div>
                                            <span className="text-[11px] text-gray-400 dark:text-gray-500 shrink-0">{formatDate(c.created_at)}</span>
                                        </div>
                                        <div className="text-sm leading-relaxed text-gray-700 dark:text-gray-200 whitespace-pre-wrap break-words">{c.body}</div>
                                    </div>
                                </div>
                            ))}
                            {loading && <div className="text-xs text-gray-500">Cargando...</div>}
                            {!loading && comments.length === 0 && (
                                <div className="text-xs text-gray-500 dark:text-gray-400">{'Sin comentarios todav\u00eda'}</div>
                            )}
                        </div>

                        <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/60 p-3 flex flex-col gap-2">
                            <textarea
                                value={newComment}
                                onChange={(e) => setNewComment(e.target.value)}
                                rows={3}
                                placeholder="Escribe un comentario..."
                                className="block w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm px-3 py-2 text-gray-800 dark:text-gray-100"
                            />
                            <div className="flex items-center justify-between gap-3">
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
                </div>

                <div className="flex flex-col gap-4">

                    {isSuperAdmin && (
                        <div className={panel}>
                            <div className={sectionTitle + ' mb-2'}>Asignado a</div>
                            <select
                                value={ticket.assigned_to_id ?? ''}
                                onChange={(e) => handleAssign(e.target.value)}
                                disabled={assigning}
                                className="block w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm px-3 py-2 text-gray-800 dark:text-gray-100"
                            >
                                <option value="">Sin asignar</option>
                                {users.map((u) => (
                                    <option key={u.id} value={u.id}>{u.name} ({u.email})</option>
                                ))}
                            </select>
                        </div>
                    )}

                    <div className={panel}>
                        <div className={sectionTitle + ' mb-3'}>Estado</div>
                        <div className="flex flex-wrap gap-1.5 mb-3">
                            {TICKET_STATUSES.map((s) => {
                                const m = STATUS_META[s];
                                const active = s === ticket.status;
                                return (
                                    <button
                                        key={s}
                                        type="button"
                                        onClick={() => handleStatusChange(s)}
                                        className={`px-2.5 py-1 rounded-full text-xs font-semibold transition ${m.bg} ${m.color} ${active ? 'ring-2 ring-offset-1 ring-offset-gray-50 dark:ring-offset-gray-900 ring-blue-500' : 'opacity-70 hover:opacity-100'}`}
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
                            className="block w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-xs px-3 py-2 text-gray-800 dark:text-gray-100"
                        />
                    </div>

                    <div className={panel}>
                        <div className={sectionTitle + ' mb-3'}>Historial</div>
                        {history.length === 0 ? (
                            <div className="text-xs text-gray-500">Sin movimientos</div>
                        ) : (
                            <ul className="flex flex-col gap-3">
                                {history.map((h, i) => (
                                    <li key={h.id} className="flex gap-2.5">
                                        <div className="flex flex-col items-center pt-1.5">
                                            <span className={`w-1.5 h-1.5 rounded-full ${i === 0 ? 'bg-blue-500' : 'bg-gray-400 dark:bg-gray-600'}`} />
                                            {i < history.length - 1 && <span className="w-px flex-1 bg-gray-200 dark:bg-gray-700 mt-1" />}
                                        </div>
                                        <div className="flex flex-col gap-0.5 pb-1">
                                            <div className="text-xs text-gray-700 dark:text-gray-200">
                                                <span className="font-medium">{h.changed_by_name || 'Usuario ' + h.changed_by_id}</span>
                                                <span className="text-gray-500 dark:text-gray-400">
                                                    {h.from_status && h.from_status !== h.to_status
                                                        ? ' ' + (STATUS_META[h.from_status as TicketStatus]?.label || h.from_status) + ' -> ' + (STATUS_META[h.to_status as TicketStatus]?.label || h.to_status)
                                                        : ' actualiz\u00f3 el ticket'}
                                                    {h.note ? ' (' + h.note + ')' : ''}
                                                </span>
                                            </div>
                                            <div className="text-[11px] text-gray-400 dark:text-gray-500">{formatDate(h.created_at)}</div>
                                        </div>
                                    </li>
                                ))}
                            </ul>
                        )}
                    </div>
                </div>
            </div>

            <div className="flex items-center justify-between gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
                <span className="text-[11px] text-gray-400 dark:text-gray-500">
                    {ticket.updated_at ? '\u00daltima actualizaci\u00f3n ' + formatDate(ticket.updated_at) : ''}
                </span>
                <Button variant="outline" onClick={onClose}>Cerrar</Button>
            </div>
        </div>
    );
}

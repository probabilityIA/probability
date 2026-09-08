'use client';

import { useEffect, useRef, useState } from 'react';
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
import { UserSelect } from '@/shared/ui';
import { getUsersAction } from '@/services/auth/users/infra/actions';
import { TypeChip, PriorityDot } from './TicketBadges';
import { META_TONE, initials, slaMeta } from './ticket-meta';
import { ACCEPTED_TYPES, attachmentError, formatSize, isImageMime, isPdfMime } from './attachment-rules';

interface Props {
    ticket: Ticket;
    isSuperAdmin: boolean;
    onClose: () => void;
    onChanged: () => void;
}

const formatDate = (value?: string) => (value ? new Date(value).toLocaleString() : '');

const SECTION = 'text-[10.5px] font-semibold uppercase tracking-[0.08em] text-gray-500 dark:text-gray-400';
const SUNKEN = 'rounded-lg border border-gray-200 bg-gray-50 dark:border-white/8 dark:bg-gray-950/60';
const CONTROL = 'block w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-gray-800 dark:border-white/10 dark:bg-gray-950/60 dark:text-gray-100';

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
    const [users, setUsers] = useState<{ id: number; name: string; email: string; avatar_url?: string }[]>([]);
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
                const list = (r?.data || []) as Array<{ id: number; name: string; email: string; avatar_url?: string; scope_code?: string; is_super_user?: boolean }>;
                setUsers(list.filter((u) => !!u.name && (u.scope_code === 'platform' || u.is_super_user)));
            } catch {}
        })();
    }, [isSuperAdmin]);

    const handleAssign = async (id: number | null) => {
        setAssigning(true);
        try {
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

    const statusMeta = STATUS_META[ticket.status];
    const sla = slaMeta(ticket);

    return (
        <div className="flex h-full min-h-0 flex-col bg-white dark:bg-gray-900">

            <div className="flex shrink-0 flex-wrap items-center gap-2 rounded-t-2xl border-b border-gray-200 bg-gray-50 px-5 py-3 dark:border-white/8 dark:bg-gray-900/80">
                <span className="font-mono text-[12.5px] font-semibold text-indigo-600 dark:text-indigo-300">{ticket.code}</span>
                <TypeChip type={ticket.type} />
                <PriorityDot priority={ticket.priority} />
                <span className={`inline-flex items-center rounded-full px-2.5 py-[3px] text-[10.5px] font-semibold ${statusMeta.bg} ${statusMeta.color}`}>
                    {statusMeta.label}
                </span>
                {ticket.escalated_to_dev && (
                    <span className="inline-flex items-center rounded-full bg-fuchsia-100 px-2.5 py-[3px] text-[10.5px] font-semibold text-fuchsia-700 dark:bg-fuchsia-500/15 dark:text-fuchsia-300">
                        Escalado a dev
                    </span>
                )}

                <span className="flex-1"></span>

                <button
                    type="button"
                    onClick={copyLink}
                    title="Copiar enlace a este ticket"
                    className="inline-flex items-center gap-1.5 rounded-lg border border-gray-300 px-2.5 py-1.5 text-[11.5px] font-medium text-gray-600 transition hover:border-indigo-400 hover:text-gray-900 dark:border-white/12 dark:text-gray-300 dark:hover:border-indigo-500 dark:hover:text-white"
                >
                    <svg className="h-3.5 w-3.5" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" d="M13.19 8.688a4.5 4.5 0 011.242 7.244l-4.5 4.5a4.5 4.5 0 01-6.364-6.364l1.757-1.757m13.35-.622l1.757-1.757a4.5 4.5 0 00-6.364-6.364l-4.5 4.5a4.5 4.5 0 001.242 7.244" />
                    </svg>
                    {copied ? 'Copiado' : 'Enlace'}
                </button>

                {isSuperAdmin && !ticket.escalated_to_dev && (
                    <button
                        type="button"
                        onClick={handleEscalate}
                        className="rounded-lg bg-indigo-500 px-3 py-1.5 text-[11.5px] font-semibold text-white transition hover:bg-indigo-400"
                    >
                        Escalar a dev
                    </button>
                )}

                <button
                    type="button"
                    onClick={onClose}
                    title="Cerrar"
                    aria-label="Cerrar"
                    className="inline-flex h-7 w-7 items-center justify-center rounded-lg text-gray-500 transition hover:bg-gray-200 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-white/10 dark:hover:text-white"
                >
                    <svg className="h-4 w-4" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>
            </div>

            <div className="shrink-0 px-5 pb-1 pt-4">
                <h2 className="text-xl font-semibold leading-[1.25] text-gray-900 dark:text-white break-words">{ticket.title}</h2>
                <div className="mt-2 flex flex-wrap items-center gap-x-2.5 gap-y-1 text-[11.5px] text-gray-500 dark:text-gray-400">
                    <span>{'Creado por '}<span className="text-gray-700 dark:text-gray-200">{ticket.created_by_name || '#' + ticket.created_by_id}</span></span>
                    <span>{'\u00b7'}</span>
                    <span>{ticket.business_name || 'Interno'}</span>
                    {ticket.category && (
                        <>
                            <span>{'\u00b7'}</span>
                            <span>{ticket.category}</span>
                        </>
                    )}
                    <span>{'\u00b7'}</span>
                    <span className="font-mono text-[11px]">{formatDate(ticket.created_at)}</span>
                    {sla && (
                        <>
                            <span>{'\u00b7'}</span>
                            <span className={META_TONE[sla.tone]}>{sla.label}</span>
                        </>
                    )}
                </div>
            </div>

            <div className="min-h-0 flex-1 overflow-y-auto overflow-x-hidden px-5 pb-5 pt-3.5">
            <div className="grid grid-cols-1 gap-0 lg:grid-cols-[minmax(0,1fr)_320px]">

                <div className="flex min-w-0 flex-col gap-4 lg:pr-5">

                    <div className="flex flex-col gap-[7px]">
                        <div className={SECTION}>{'Descripci\u00f3n'}</div>
                        <div className={`${SUNKEN} px-3.5 py-3 text-[13px] leading-[1.55] text-gray-700 dark:text-gray-300 whitespace-pre-wrap break-words`}>
                            {ticket.description}
                        </div>
                    </div>

                    <div className="flex flex-col gap-[7px]">
                        <div className="flex items-baseline gap-2">
                            <span className={SECTION}>Adjuntos</span>
                            <span className="text-[10.5px] text-gray-400 dark:text-gray-500">{'im\u00e1genes y PDF \u00b7 hasta 10 MB'}</span>
                        </div>

                        {uploadError && (
                            <div className="flex items-start gap-2 rounded-lg border border-red-200 bg-red-50 px-3 py-2 dark:border-red-500/30 dark:bg-red-500/10">
                                <svg className="mt-0.5 h-4 w-4 shrink-0 text-red-500" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z" />
                                </svg>
                                <span className="flex-1 text-xs text-red-700 dark:text-red-300">{uploadError}</span>
                                <button type="button" onClick={() => setUploadError('')} className="text-xs text-red-500 hover:text-red-700">x</button>
                            </div>
                        )}

                        {uploading && (
                            <div className={`${SUNKEN} flex items-center gap-3 px-3.5 py-2.5`}>
                                <svg className="h-4 w-4 animate-spin text-indigo-500" fill="none" viewBox="0 0 24 24">
                                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                                </svg>
                                <span className="truncate text-[12.5px] text-gray-600 dark:text-gray-300">{'Subiendo ' + uploadingName}</span>
                            </div>
                        )}

                        {attachments.length > 0 && (
                            <div className="grid grid-cols-3 gap-2.5 sm:grid-cols-4">
                                {attachments.map((a) => (
                                    <div key={a.id} className={`${SUNKEN} overflow-hidden`}>
                                        <a href={a.file_url} target="_blank" rel="noreferrer" className="block h-20 bg-gray-100 dark:bg-gray-900">
                                            {isImageMime(a.mime_type) ? (
                                                <img src={a.file_url} alt={a.file_name} className="h-20 w-full object-cover" />
                                            ) : (
                                                <span className={`flex h-20 w-full flex-col items-center justify-center gap-1 ${isPdfMime(a.mime_type) ? 'bg-red-50 dark:bg-red-500/10' : ''}`}>
                                                    <svg className={`h-6 w-6 ${isPdfMime(a.mime_type) ? 'text-red-400' : 'text-gray-400'}`} fill="none" stroke="currentColor" strokeWidth={1.5} viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25M9 16.5v.75m3-3v3M15 12v5.25m-4.5-15H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
                                                    </svg>
                                                    {isPdfMime(a.mime_type) && <span className="font-mono text-[10px] font-medium tracking-wide text-red-400">PDF</span>}
                                                </span>
                                            )}
                                        </a>
                                        <div className="flex flex-col gap-0.5 p-2">
                                            <div className="truncate text-[11px] text-gray-700 dark:text-gray-200" title={a.file_name}>{a.file_name}</div>
                                            <div className="flex items-center justify-between">
                                                <span className="text-[10px] text-gray-400 dark:text-gray-500">{formatSize(a.size)}</span>
                                                {isSuperAdmin && (
                                                    <button type="button" onClick={() => removeAttachment(a.id)} className="text-[10px] text-red-500 hover:text-red-600">Eliminar</button>
                                                )}
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}

                        <label
                            htmlFor="ticket-file-input"
                            onDragOver={(e) => { e.preventDefault(); setDragOver(true); }}
                            onDragLeave={() => setDragOver(false)}
                            onDrop={handleDrop}
                            className={`cursor-pointer rounded-lg border border-dashed p-4 text-center text-[12px] transition ${
                                dragOver
                                    ? 'border-indigo-400 bg-indigo-50 text-indigo-600 dark:bg-indigo-500/10 dark:text-indigo-300'
                                    : 'border-gray-300 text-gray-500 hover:border-gray-400 dark:border-white/18 dark:text-gray-400 dark:hover:border-white/30'
                            }`}
                        >
                            {dragOver ? (
                                'Suelta para adjuntar'
                            ) : (
                                <>
                                    {'Arrastra un archivo aqu\u00ed o '}
                                    <span className="font-medium text-indigo-600 dark:text-indigo-300">{'selecci\u00f3nalo'}</span>
                                </>
                            )}
                            <input ref={fileRef} type="file" accept={ACCEPTED_TYPES} onChange={handleUpload} className="hidden" id="ticket-file-input" />
                        </label>
                    </div>

                    <div className="flex flex-col gap-2.5">
                        <div className={SECTION}>{'Comentarios \u00b7 ' + comments.length}</div>

                        <div className="flex max-h-96 flex-col gap-2 overflow-y-auto">
                            {comments.map((c) => (
                                <div key={c.id} className="flex gap-2.5">
                                    <span className={`flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-[9.5px] font-bold text-white ${c.is_internal ? 'bg-amber-500' : 'bg-sky-500'}`}>
                                        {initials(c.user_name)}
                                    </span>
                                    <div className={`min-w-0 flex-1 rounded-lg border px-3 py-2.5 ${c.is_internal ? 'border-amber-200 bg-amber-50 dark:border-amber-500/25 dark:bg-amber-500/10' : 'border-gray-200 bg-gray-50 dark:border-white/8 dark:bg-gray-950/60'}`}>
                                        <div className="flex items-center gap-2">
                                            <span className="truncate text-[11.5px] font-semibold text-gray-700 dark:text-gray-200">{c.user_name || 'Usuario ' + c.user_id}</span>
                                            {c.is_internal && (
                                                <span className="shrink-0 rounded-full bg-amber-200 px-1.5 text-[10px] font-semibold text-amber-800 dark:bg-amber-500/25 dark:text-amber-200">Nota interna</span>
                                            )}
                                            <span className="ml-auto shrink-0 text-[10.5px] text-gray-400 dark:text-gray-500">{formatDate(c.created_at)}</span>
                                        </div>
                                        <div className="mt-1 whitespace-pre-wrap break-words text-[12.5px] leading-[1.5] text-gray-700 dark:text-gray-300">{c.body}</div>
                                    </div>
                                </div>
                            ))}
                            {loading && <div className="text-xs text-gray-500">Cargando...</div>}
                            {!loading && comments.length === 0 && (
                                <div className="text-xs text-gray-400 dark:text-gray-500">{'Sin comentarios todav\u00eda'}</div>
                            )}
                        </div>

                        <div className={`${SUNKEN} flex flex-col gap-2 px-3 py-2.5`}>
                            <textarea
                                value={newComment}
                                onChange={(e) => setNewComment(e.target.value)}
                                rows={2}
                                placeholder="Escribe un comentario..."
                                className="block w-full resize-y bg-transparent text-[12.5px] text-gray-800 outline-none placeholder:text-gray-400 dark:text-gray-100 dark:placeholder:text-gray-500"
                            />
                            <div className="flex items-center gap-2">
                                {isSuperAdmin ? (
                                    <label className="flex items-center gap-2 text-[11.5px] text-gray-600 dark:text-gray-400">
                                        <input type="checkbox" checked={internalComment} onChange={(e) => setInternalComment(e.target.checked)} />
                                        Nota interna (solo super admins)
                                    </label>
                                ) : <span />}
                                <span className="flex-1"></span>
                                <button
                                    type="button"
                                    onClick={submitComment}
                                    disabled={posting || !newComment.trim()}
                                    className="rounded-md bg-gray-900 px-3.5 py-1.5 text-[12px] font-semibold text-white transition hover:bg-gray-700 disabled:opacity-40 dark:bg-gray-100 dark:text-gray-900 dark:hover:bg-white"
                                >
                                    {posting ? 'Enviando...' : 'Comentar'}
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="mt-5 flex flex-col gap-3.5 border-gray-200 lg:mt-0 lg:border-l lg:pl-5 dark:border-white/8">

                    {isSuperAdmin && (
                        <div className="flex flex-col gap-[7px]">
                            <div className={SECTION}>Asignado a</div>
                            <UserSelect
                                value={ticket.assigned_to_id ?? null}
                                options={users}
                                onChange={handleAssign}
                                disabled={assigning}
                                fallbackName={ticket.assigned_to_name}
                                fallbackAvatarUrl={ticket.assigned_to_avatar_url}
                                showEmail
                                buttonClassName="py-2 text-[12.5px] font-medium"
                            />
                        </div>
                    )}

                    <div className="flex flex-col gap-[7px]">
                        <div className={SECTION}>Estado</div>
                        <div className="flex flex-wrap gap-1.5">
                            {TICKET_STATUSES.map((s) => {
                                const m = STATUS_META[s];
                                const active = s === ticket.status;
                                return (
                                    <button
                                        key={s}
                                        type="button"
                                        onClick={() => handleStatusChange(s)}
                                        className={`rounded-full px-2.5 py-1 text-[11px] font-semibold transition ${m.bg} ${m.color} ${
                                            active ? 'ring-2 ring-inset ring-indigo-500' : 'opacity-70 hover:opacity-100'
                                        }`}
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
                            className={`${CONTROL} text-[11.5px]`}
                        />
                    </div>

                    <div className="flex min-h-0 flex-col gap-2">
                        <div className={SECTION}>Historial</div>
                        {history.length === 0 ? (
                            <div className="text-xs text-gray-400 dark:text-gray-500">Sin movimientos</div>
                        ) : (
                            <ul className="flex max-h-[230px] flex-col overflow-y-auto pr-1">
                                {history.map((h, i) => (
                                    <li key={h.id} className="flex gap-2.5 py-[5px]">
                                        <div className="flex w-2 flex-none flex-col items-center">
                                            <span className={`mt-1 h-1.5 w-1.5 rounded-full ${i === 0 ? 'bg-indigo-500' : 'bg-gray-400 dark:bg-gray-600'}`} />
                                            {i < history.length - 1 && <span className="w-px flex-1 bg-gray-200 dark:bg-white/10" />}
                                        </div>
                                        <div className="pb-1 text-[11.5px] leading-[1.45] text-gray-500 dark:text-gray-400">
                                            <span className="font-semibold text-gray-700 dark:text-gray-300">{h.changed_by_name || 'Usuario ' + h.changed_by_id}</span>
                                            <span>
                                                {h.from_status && h.from_status !== h.to_status
                                                    ? ' ' + (STATUS_META[h.from_status as TicketStatus]?.label || h.from_status) + ' -> ' + (STATUS_META[h.to_status as TicketStatus]?.label || h.to_status)
                                                    : ' actualiz\u00f3 el ticket'}
                                                {h.note ? ' (' + h.note + ')' : ''}
                                            </span>
                                            <div className="mt-px font-mono text-[10px] text-gray-400 dark:text-gray-500">{formatDate(h.created_at)}</div>
                                        </div>
                                    </li>
                                ))}
                            </ul>
                        )}
                    </div>
                </div>
            </div>
            </div>

            <div className="flex shrink-0 items-center gap-3 rounded-b-2xl border-t border-gray-200 bg-gray-50 px-5 py-2.5 dark:border-white/8 dark:bg-gray-900/80">
                <span className="font-mono text-[11px] text-gray-400 dark:text-gray-500">
                    {ticket.updated_at ? '\u00daltima actualizaci\u00f3n ' + formatDate(ticket.updated_at) : ''}
                </span>
                <span className="flex-1"></span>
                {isSuperAdmin && (
                    <button
                        type="button"
                        onClick={handleDelete}
                        title="Eliminar ticket"
                        className="inline-flex items-center gap-1.5 rounded-lg border border-red-300 px-2.5 py-1.5 text-[11.5px] font-semibold text-red-600 transition hover:bg-red-50 dark:border-red-500/35 dark:text-red-400 dark:hover:bg-red-500/10"
                    >
                        <svg className="h-3.5 w-3.5" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                        </svg>
                        Eliminar ticket
                    </button>
                )}
            </div>
        </div>
    );
}

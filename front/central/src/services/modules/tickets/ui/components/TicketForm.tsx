'use client';

import { useRef, useState } from 'react';
import { Button, Input } from '@/shared/ui';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import { Sprint } from '@/services/modules/sprints/domain/types';
import {
    CreateTicketDTO,
    TICKET_TYPES,
    TICKET_PRIORITIES,
    TICKET_AREAS,
    TICKET_SEVERITIES,
    TICKET_SOURCES,
    TYPE_META,
    PRIORITY_META,
    AREA_META,
    SEVERITY_META,
    SOURCE_META,
    TicketType,
    TicketPriority,
    TicketSeverity,
    TicketSource,
    TicketArea,
} from '../../domain/types';
import { ACCEPTED_TYPES, attachmentError, formatSize } from './attachment-rules';

export interface CreateTicketPayload {
    dto: CreateTicketDTO;
    files: File[];
    comment: string;
    commentInternal: boolean;
}

interface UserOption {
    id: number;
    name: string;
}

interface Props {
    isSuperAdmin: boolean;
    users?: UserOption[];
    sprints?: Sprint[];
    onSubmit: (payload: CreateTicketPayload) => Promise<void>;
    onCancel: () => void;
    submitting?: boolean;
    submitLabel?: string;
}

const labelClass = 'block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1';
const controlClass = 'block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500';

export default function TicketForm({ isSuperAdmin, users = [], sprints = [], onSubmit, onCancel, submitting, submitLabel }: Props) {
    const { businesses } = useBusinessesSimple();
    const fileRef = useRef<HTMLInputElement>(null);
    const [businessId, setBusinessId] = useState<string>('');
    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [type, setType] = useState<TicketType>('support');
    const [priority, setPriority] = useState<TicketPriority>('medium');
    const [severity, setSeverity] = useState<TicketSeverity>('');
    const [area, setArea] = useState<TicketArea>('soporte');
    const [source, setSource] = useState<TicketSource>(isSuperAdmin ? 'internal' : 'business');
    const [category, setCategory] = useState('');
    const [dueDate, setDueDate] = useState('');
    const [assignedToId, setAssignedToId] = useState<string>('');
    const [sprintId, setSprintId] = useState<string>('');
    const [files, setFiles] = useState<File[]>([]);
    const [fileError, setFileError] = useState<string | null>(null);
    const [dragOver, setDragOver] = useState(false);
    const [comment, setComment] = useState('');
    const [commentInternal, setCommentInternal] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const addFiles = (incoming: FileList | null) => {
        if (!incoming || incoming.length === 0) return;
        setFileError(null);
        const accepted: File[] = [];
        for (const file of Array.from(incoming)) {
            const invalid = attachmentError(file);
            if (invalid) {
                setFileError(invalid);
                continue;
            }
            accepted.push(file);
        }
        if (accepted.length > 0) {
            setFiles((prev) => [
                ...prev,
                ...accepted.filter((f) => !prev.some((p) => p.name === f.name && p.size === f.size)),
            ]);
        }
    };

    const removeFile = (index: number) => {
        setFiles((prev) => prev.filter((_, i) => i !== index));
    };

    const handlePick = (e: React.ChangeEvent<HTMLInputElement>) => {
        addFiles(e.target.files);
        if (fileRef.current) fileRef.current.value = '';
    };

    const handleDrop = (e: React.DragEvent) => {
        e.preventDefault();
        setDragOver(false);
        addFiles(e.dataTransfer.files);
    };

    const submit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        if (!title.trim() || !description.trim()) {
            setError('T\u00edtulo y descripci\u00f3n son obligatorios');
            return;
        }
        try {
            await onSubmit({
                dto: {
                    title: title.trim(),
                    description: description.trim(),
                    type,
                    priority,
                    severity: severity || undefined,
                    area,
                    source,
                    category: category.trim() || undefined,
                    due_date: dueDate || undefined,
                    assigned_to_id: assignedToId ? Number(assignedToId) : null,
                    sprint_id: sprintId ? Number(sprintId) : null,
                    business_id: isSuperAdmin ? (businessId ? Number(businessId) : null) : undefined,
                },
                files,
                comment: comment.trim(),
                commentInternal: isSuperAdmin ? commentInternal : false,
            });
        } catch (err: any) {
            setError(err?.message || 'Error al crear ticket');
        }
    };

    return (
        <form onSubmit={submit} className="space-y-4">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-5">
                <div className="space-y-4">
                    <div>
                        <label className={labelClass}>{'T\u00edtulo *'}</label>
                        <Input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="Resumen breve" />
                    </div>

                    <div>
                        <label className={labelClass}>{'Descripci\u00f3n *'}</label>
                        <textarea
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            rows={8}
                            className={controlClass}
                            placeholder="Detalla el problema, mejora o solicitud"
                        />
                    </div>

                    <div>
                        <label className={labelClass}>Comentario inicial</label>
                        <textarea
                            value={comment}
                            onChange={(e) => setComment(e.target.value)}
                            rows={3}
                            className={controlClass}
                            placeholder="Opcional: contexto adicional para el equipo"
                        />
                        {isSuperAdmin && (
                            <label className="mt-2 flex items-center gap-2 text-xs text-gray-600 dark:text-gray-300">
                                <input type="checkbox" checked={commentInternal} onChange={(e) => setCommentInternal(e.target.checked)} />
                                Nota interna (solo super admins)
                            </label>
                        )}
                    </div>
                </div>

                <div className="space-y-4">
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        <div>
                            <label className={labelClass}>Tipo</label>
                            <select value={type} onChange={(e) => setType(e.target.value as TicketType)} className={controlClass}>
                                {TICKET_TYPES.map((t) => (
                                    <option key={t} value={t}>{TYPE_META[t].label}</option>
                                ))}
                            </select>
                        </div>

                        <div>
                            <label className={labelClass}>Prioridad</label>
                            <select value={priority} onChange={(e) => setPriority(e.target.value as TicketPriority)} className={controlClass}>
                                {TICKET_PRIORITIES.map((p) => (
                                    <option key={p} value={p}>{PRIORITY_META[p].label}</option>
                                ))}
                            </select>
                        </div>

                        <div>
                            <label className={labelClass}>Severidad</label>
                            <select value={severity} onChange={(e) => setSeverity(e.target.value as TicketSeverity)} className={controlClass}>
                                {TICKET_SEVERITIES.map((s) => (
                                    <option key={s || 'none'} value={s}>{SEVERITY_META[s].label}</option>
                                ))}
                            </select>
                        </div>

                        <div>
                            <label className={labelClass}>{'\u00c1rea'}</label>
                            <select value={area} onChange={(e) => setArea(e.target.value as TicketArea)} className={controlClass}>
                                {TICKET_AREAS.map((a) => (
                                    <option key={a} value={a}>{AREA_META[a].label}</option>
                                ))}
                            </select>
                        </div>

                        <div>
                            <label className={labelClass}>Origen</label>
                            <select value={source} onChange={(e) => setSource(e.target.value as TicketSource)} className={controlClass}>
                                {TICKET_SOURCES.map((s) => (
                                    <option key={s} value={s}>{SOURCE_META[s].label}</option>
                                ))}
                            </select>
                        </div>

                        <div>
                            <label className={labelClass}>{'Categor\u00eda'}</label>
                            <Input value={category} onChange={(e) => setCategory(e.target.value)} placeholder={'Ej: env\u00edos, facturaci\u00f3n, frontend'} />
                        </div>

                        <div>
                            <label className={labelClass}>Fecha objetivo</label>
                            <Input type="date" value={dueDate} onChange={(e) => setDueDate(e.target.value)} />
                        </div>

                        <div>
                            <label className={labelClass}>Asignado a</label>
                            <select value={assignedToId} onChange={(e) => setAssignedToId(e.target.value)} className={controlClass}>
                                <option value="">Sin asignar</option>
                                {users.map((u) => (
                                    <option key={u.id} value={u.id}>{u.name}</option>
                                ))}
                            </select>
                        </div>

                        <div>
                            <label className={labelClass}>Sprint</label>
                            <select value={sprintId} onChange={(e) => setSprintId(e.target.value)} className={controlClass}>
                                <option value="">Sin sprint</option>
                                {sprints.map((s) => (
                                    <option key={s.id} value={s.id}>{s.name}</option>
                                ))}
                            </select>
                        </div>

                        {isSuperAdmin && (
                            <div>
                                <label className={labelClass}>Negocio</label>
                                <select value={businessId} onChange={(e) => setBusinessId(e.target.value)} className={controlClass}>
                                    <option value="">Interno (sin negocio)</option>
                                    {businesses.map((b) => (
                                        <option key={b.id} value={b.id}>{b.name}</option>
                                    ))}
                                </select>
                            </div>
                        )}
                    </div>

                    <div>
                        <label className={labelClass}>Adjuntos</label>
                        <label
                            htmlFor="ticket-create-file-input"
                            onDragOver={(e) => { e.preventDefault(); setDragOver(true); }}
                            onDragLeave={() => setDragOver(false)}
                            onDrop={handleDrop}
                            className={`flex flex-col items-center justify-center gap-1.5 rounded-xl border-[1.5px] border-dashed p-4 cursor-pointer transition ${dragOver ? 'border-sky-400 bg-sky-50 dark:bg-sky-900/20' : 'border-gray-300 dark:border-gray-600 bg-gray-50/60 dark:bg-gray-900/40 hover:border-gray-400 dark:hover:border-gray-500'}`}
                        >
                            <svg className={`w-6 h-6 ${dragOver ? 'text-sky-400' : 'text-gray-400 dark:text-gray-500'}`} fill="none" stroke="currentColor" strokeWidth={1.5} viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M12 16.5V9.75m0 0l3 3m-3-3l-3 3M6.75 19.5a4.5 4.5 0 01-1.41-8.775 5.25 5.25 0 0110.233-2.33 3 3 0 013.758 3.848A3.752 3.752 0 0118 19.5H6.75z" />
                            </svg>
                            <span className="text-sm text-gray-600 dark:text-gray-300">
                                {dragOver ? 'Suelta para adjuntar' : 'Arrastra archivos aqu\u00ed'}
                            </span>
                            {!dragOver && (
                                <span className="text-[11px] text-gray-400 dark:text-gray-500">
                                    {'o '}<span className="text-sky-600 dark:text-sky-400 font-medium">{'selecci\u00f3nalos'}</span>{' desde tu equipo'}
                                </span>
                            )}
                            <input
                                ref={fileRef}
                                id="ticket-create-file-input"
                                type="file"
                                multiple
                                accept={ACCEPTED_TYPES}
                                onChange={handlePick}
                                className="hidden"
                            />
                        </label>

                        {fileError && <div className="mt-2 text-xs text-red-600 dark:text-red-400">{fileError}</div>}

                        {files.length > 0 && (
                            <ul className="mt-2 space-y-1">
                                {files.map((f, i) => (
                                    <li key={f.name + f.size + i} className="flex items-center justify-between gap-2 rounded-lg border border-gray-200 dark:border-gray-700 px-3 py-1.5">
                                        <span className="min-w-0 truncate text-xs text-gray-700 dark:text-gray-200" title={f.name}>{f.name}</span>
                                        <span className="flex items-center gap-2 shrink-0">
                                            <span className="text-[11px] text-gray-400 dark:text-gray-500">{formatSize(f.size)}</span>
                                            <button
                                                type="button"
                                                onClick={() => removeFile(i)}
                                                disabled={submitting}
                                                className="text-[11px] font-semibold text-red-500 hover:text-red-600 disabled:opacity-50"
                                            >
                                                Quitar
                                            </button>
                                        </span>
                                    </li>
                                ))}
                            </ul>
                        )}
                    </div>
                </div>
            </div>

            {error && <div className="text-sm text-red-600 dark:text-red-400">{error}</div>}

            <div className="flex justify-end gap-2 pt-2">
                <Button type="button" variant="outline" onClick={onCancel} disabled={submitting}>Cancelar</Button>
                <Button type="submit" variant="primary" disabled={submitting}>
                    {submitting ? (submitLabel || 'Creando...') : 'Crear ticket'}
                </Button>
            </div>
        </form>
    );
}

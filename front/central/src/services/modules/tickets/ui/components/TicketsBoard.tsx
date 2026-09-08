'use client';

import { useMemo, useState } from 'react';
import { Ticket, TicketStatus, PRIORITY_ACCENT, PRIORITY_META } from '../../domain/types';
import { StatusChip, TypeChip } from './TicketBadges';
import { META_TONE, bySla, cardMeta, channelLabel, initials, isDueSoon } from './ticket-meta';

export interface BoardColumn {
    key: string;
    title: string;
    statuses: TicketStatus[];
    dropStatus: TicketStatus;
    accent: string;
    trackStale: boolean;
}

export const BOARD_COLUMNS: BoardColumn[] = [
    { key: 'todo', title: 'To Do', statuses: ['open', 'in_review'], dropStatus: 'open', accent: 'bg-gray-400 dark:bg-gray-500', trackStale: false },
    { key: 'in_progress', title: 'In Progress', statuses: ['in_development'], dropStatus: 'in_development', accent: 'bg-violet-400', trackStale: true },
    { key: 'in_testing', title: 'In Testing', statuses: ['testing'], dropStatus: 'testing', accent: 'bg-cyan-400', trackStale: true },
    { key: 'blocked', title: 'Blocked', statuses: ['blocked'], dropStatus: 'blocked', accent: 'bg-red-400', trackStale: true },
    { key: 'done', title: 'Done', statuses: ['resolved', 'closed', 'wont_fix'], dropStatus: 'resolved', accent: 'bg-emerald-400', trackStale: false },
];

const CARDS_PER_COLUMN = 8;

interface TicketsBoardProps {
    tickets: Ticket[];
    loading: boolean;
    canDrag: boolean;
    updatingId: number | null;
    onOpen: (ticket: Ticket) => void;
    onMove: (id: number, status: TicketStatus) => void;
    getAvatarUrl?: (ticket: Ticket) => string;
    sortBySla?: boolean;
}

export default function TicketsBoard({
    tickets,
    loading,
    canDrag,
    updatingId,
    onOpen,
    onMove,
    getAvatarUrl,
    sortBySla = true,
}: TicketsBoardProps) {
    const [dragId, setDragId] = useState<number | null>(null);
    const [overKey, setOverKey] = useState<string | null>(null);
    const [expanded, setExpanded] = useState<Record<string, boolean>>({});
    const [openedRails, setOpenedRails] = useState<Record<string, boolean>>({});

    const now = Date.now();

    const grouped = useMemo(() => {
        const sorter = bySla(now);
        return BOARD_COLUMNS.map((column) => {
            const items = tickets.filter((t) => column.statuses.includes(t.status));
            return { column, items: sortBySla ? items.sort(sorter) : items };
        });
    }, [tickets, now, sortBySla]);

    const handleDragStart = (e: React.DragEvent<HTMLDivElement>, ticket: Ticket) => {
        if (!canDrag) return;
        setDragId(ticket.id);
        e.dataTransfer.effectAllowed = 'move';
        e.dataTransfer.setData('text/plain', String(ticket.id));
    };

    const handleDragEnd = () => {
        setDragId(null);
        setOverKey(null);
    };

    const handleDragOver = (e: React.DragEvent<HTMLDivElement>, column: BoardColumn) => {
        if (!canDrag) return;
        e.preventDefault();
        e.dataTransfer.dropEffect = 'move';
        if (overKey !== column.key) setOverKey(column.key);
    };

    const handleDragLeave = (e: React.DragEvent<HTMLDivElement>, column: BoardColumn) => {
        if (e.currentTarget.contains(e.relatedTarget as Node)) return;
        if (overKey === column.key) setOverKey(null);
    };

    const handleDrop = (e: React.DragEvent<HTMLDivElement>, column: BoardColumn) => {
        e.preventDefault();
        setOverKey(null);
        if (!canDrag) return;
        const raw = e.dataTransfer.getData('text/plain');
        const id = Number(raw) || dragId;
        setDragId(null);
        if (!id) return;
        const ticket = tickets.find((t) => t.id === id);
        if (!ticket) return;
        if (column.statuses.includes(ticket.status)) return;
        onMove(id, column.dropStatus);
    };

    return (
        <div className="relative">
            {loading && (
                <div className="absolute inset-0 z-10 flex items-center justify-center rounded-xl bg-white/80 dark:bg-gray-950/70 backdrop-blur-sm transition-opacity duration-200">
                    <div className="flex flex-col items-center gap-2">
                        <div className="w-8 h-8 border-4 border-purple-500 border-t-transparent rounded-full animate-spin"></div>
                        <p className="text-sm text-gray-600 dark:text-gray-300">Actualizando...</p>
                    </div>
                </div>
            )}

            <div className={`flex w-full items-stretch gap-2.5 overflow-x-auto pb-3 transition-opacity duration-200 ${loading ? 'opacity-50' : 'opacity-100'}`}>
                {grouped.map(({ column, items }) => {
                    const isOver = overKey === column.key;
                    const collapsed = items.length === 0 && !isOver && !openedRails[column.key];
                    const dueSoon = items.filter((t) => isDueSoon(t, now)).length;
                    const showAll = !!expanded[column.key];
                    const visible = showAll ? items : items.slice(0, CARDS_PER_COLUMN);
                    const hidden = items.length - visible.length;

                    if (collapsed) {
                        return (
                            <div
                                key={column.key}
                                role="button"
                                tabIndex={0}
                                title={'Columna vac\u00eda: ' + column.title}
                                aria-label={'Columna vac\u00eda: ' + column.title}
                                onClick={() => setOpenedRails((prev) => ({ ...prev, [column.key]: true }))}
                                onKeyDown={(e) => { if (e.key === 'Enter' || e.key === ' ') setOpenedRails((prev) => ({ ...prev, [column.key]: true })); }}
                                onDragOver={(e) => handleDragOver(e, column)}
                                onDragLeave={(e) => handleDragLeave(e, column)}
                                onDrop={(e) => handleDrop(e, column)}
                                className={`w-9 flex-none cursor-pointer rounded-lg border transition-colors ${
                                    column.key === 'blocked'
                                        ? 'border-red-300/60 bg-red-50/40 dark:border-red-500/25 dark:bg-gray-900/60'
                                        : 'border-gray-200 bg-gray-50 dark:border-gray-700/70 dark:bg-gray-900/60'
                                }`}
                            >
                                <div className="flex min-h-[340px] flex-col items-center gap-2 py-2.5">
                                    <span className={`rounded px-1.5 font-mono text-[10.5px] font-bold ${
                                        column.key === 'blocked'
                                            ? 'bg-red-100 text-red-600 dark:bg-red-500/15 dark:text-red-400'
                                            : 'bg-gray-200 text-gray-500 dark:bg-white/10 dark:text-gray-400'
                                    }`}>0</span>
                                    <h3
                                        className={`text-[11px] font-semibold uppercase tracking-wider ${
                                            column.key === 'blocked' ? 'text-red-500 dark:text-red-400' : 'text-gray-500 dark:text-gray-400'
                                        }`}
                                        style={{ writingMode: 'vertical-rl' }}
                                    >
                                        {column.title}
                                    </h3>
                                </div>
                            </div>
                        );
                    }

                    return (
                        <div
                            key={column.key}
                            onDragOver={(e) => handleDragOver(e, column)}
                            onDragLeave={(e) => handleDragLeave(e, column)}
                            onDrop={(e) => handleDrop(e, column)}
                            className={`flex w-[296px] flex-none flex-col gap-2 rounded-lg border p-1 transition-colors ${
                                isOver
                                    ? 'border-purple-400 bg-purple-50/60 dark:border-purple-500 dark:bg-purple-500/10'
                                    : 'border-transparent'
                            }`}
                        >
                            <div className="flex items-center gap-2 px-1 py-0.5">
                                <span className={`h-2 w-2 rounded-sm ${column.accent}`}></span>
                                <h3 className="text-[11.5px] font-semibold uppercase tracking-wider text-gray-700 dark:text-gray-200">
                                    {column.title}
                                </h3>
                                <span className="rounded px-1.5 font-mono text-[10.5px] font-bold text-gray-500 bg-gray-200 dark:bg-white/10 dark:text-gray-400">
                                    {items.length}
                                </span>
                                <span className="flex-1"></span>
                                {dueSoon > 0 && (
                                    <span className="font-mono text-[10px] font-bold text-amber-600 dark:text-amber-400">
                                        {dueSoon + ' por vencer'}
                                    </span>
                                )}
                            </div>

                            <div className="flex flex-1 flex-col gap-2">
                                {items.length === 0 && (
                                    <div className={`flex min-h-[340px] flex-1 items-center justify-center rounded-lg border border-dashed p-4 text-center text-[11.5px] text-gray-400 dark:text-gray-500 ${
                                        column.key === 'blocked'
                                            ? 'border-red-300/60 dark:border-red-500/25'
                                            : 'border-gray-300 dark:border-gray-700'
                                    }`}>
                                        Sin tickets
                                    </div>
                                )}

                                {visible.map((t) => {
                                    const avatarUrl = getAvatarUrl ? getAvatarUrl(t) : '';
                                    const isUpdating = updatingId === t.id;
                                    const isDone = column.key === 'done';
                                    const meta = cardMeta(t, column.trackStale, now);
                                    const assignee = t.assigned_to_name || '';
                                    return (
                                        <div
                                            key={t.id}
                                            draggable={canDrag}
                                            onDragStart={(e) => handleDragStart(e, t)}
                                            onDragEnd={handleDragEnd}
                                            onClick={() => onOpen(t)}
                                            className={`flex flex-col gap-1.5 rounded-lg border border-l-[3px] px-2.5 py-2.5 transition-all ${PRIORITY_ACCENT[t.priority]} ${
                                                isDone
                                                    ? 'border-gray-200 bg-gray-50 opacity-75 hover:opacity-100 dark:border-gray-700/60 dark:bg-gray-900/70'
                                                    : 'border-gray-200 bg-white hover:border-purple-400 hover:shadow-sm dark:border-gray-700/70 dark:bg-gray-800/90 dark:hover:border-purple-500 dark:hover:bg-gray-800'
                                            } ${canDrag ? 'cursor-grab active:cursor-grabbing' : 'cursor-pointer'} ${
                                                dragId === t.id ? 'opacity-40' : ''
                                            } ${isUpdating ? 'ring-2 ring-purple-400 animate-pulse' : ''}`}
                                            title={PRIORITY_META[t.priority].label + ' | ' + t.title}
                                        >
                                            <div className="flex items-center gap-1.5">
                                                <span className="font-mono text-[10.5px] font-medium text-gray-500 dark:text-gray-400">{t.code}</span>
                                                <TypeChip type={t.type} />
                                                {column.statuses.length > 1 && <StatusChip status={t.status} />}
                                                <span className="flex-1"></span>
                                                <span className={`shrink-0 font-mono text-[10.5px] ${META_TONE[meta.tone]}`}>{meta.label}</span>
                                            </div>

                                            <p className={`line-clamp-2 text-[13px] font-medium leading-[1.35] ${
                                                isDone ? 'text-gray-500 dark:text-gray-400' : 'text-gray-900 dark:text-gray-100'
                                            }`}>
                                                {t.title}
                                            </p>

                                            <div className="flex items-center gap-2.5 text-[11px] text-gray-500 dark:text-gray-400">
                                                <span className="inline-flex items-center gap-1" title="Comentarios">
                                                    <svg className="h-3 w-3" fill="none" stroke="currentColor" strokeWidth={1.8} viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" d="M8.625 12a.375.375 0 11-.75 0 .375.375 0 01.75 0zm3.75 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm3.75 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zM21 12c0 4.556-4.03 8.25-9 8.25a9.76 9.76 0 01-2.555-.337A5.972 5.972 0 015.41 20.97a5.969 5.969 0 01-.474-.065 4.48 4.48 0 00.978-2.025c.09-.457-.133-.901-.467-1.226C3.93 16.178 3 14.189 3 12c0-4.556 4.03-8.25 9-8.25s9 3.694 9 8.25z" />
                                                    </svg>
                                                    {t.comments_count}
                                                </span>
                                                {t.attachments_count > 0 && (
                                                    <span className="inline-flex items-center gap-1" title="Adjuntos">
                                                        <svg className="h-3 w-3" fill="none" stroke="currentColor" strokeWidth={1.8} viewBox="0 0 24 24">
                                                            <path strokeLinecap="round" strokeLinejoin="round" d="M18.375 12.739l-7.693 7.693a4.5 4.5 0 01-6.364-6.364l10.94-10.94A3 3 0 1119.5 7.372L8.552 18.32m.009-.01l-.01.01m5.699-9.941l-7.81 7.81a1.5 1.5 0 002.112 2.13" />
                                                        </svg>
                                                        {t.attachments_count}
                                                    </span>
                                                )}
                                                <span className="flex-1"></span>
                                                <span className="max-w-[45%] truncate font-mono text-[10px]" title={channelLabel(t)}>
                                                    {channelLabel(t)}
                                                </span>
                                                {assignee ? (
                                                    avatarUrl ? (
                                                        <img
                                                            src={avatarUrl}
                                                            alt={assignee}
                                                            title={assignee}
                                                            className="h-5 w-5 shrink-0 rounded-full object-cover ring-1 ring-gray-200 dark:ring-gray-600"
                                                        />
                                                    ) : (
                                                        <span
                                                            title={assignee}
                                                            className="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-violet-500 text-[9px] font-bold text-white"
                                                        >
                                                            {initials(assignee)}
                                                        </span>
                                                    )
                                                ) : (
                                                    <span className="shrink-0 font-mono text-[10px] font-semibold text-red-500 dark:text-red-400">Sin asignar</span>
                                                )}
                                            </div>
                                        </div>
                                    );
                                })}

                                {hidden > 0 && (
                                    <button
                                        type="button"
                                        onClick={() => setExpanded((prev) => ({ ...prev, [column.key]: true }))}
                                        className="rounded-lg border border-dashed border-gray-300 py-1.5 text-[11.5px] font-medium text-gray-500 transition-colors hover:border-purple-400 hover:text-purple-600 dark:border-gray-700 dark:text-gray-400 dark:hover:border-purple-500 dark:hover:text-purple-300"
                                    >
                                        {'Ver ' + hidden + ' m\u00e1s'}
                                    </button>
                                )}
                            </div>
                        </div>
                    );
                })}
            </div>
        </div>
    );
}

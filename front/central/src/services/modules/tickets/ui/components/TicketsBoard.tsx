'use client';

import { useState } from 'react';
import { Ticket, TicketStatus } from '../../domain/types';
import { StatusBadge, PriorityBadge, TypeBadge } from './TicketBadges';

export interface BoardColumn {
    key: string;
    title: string;
    statuses: TicketStatus[];
    dropStatus: TicketStatus;
    accent: string;
}

export const BOARD_COLUMNS: BoardColumn[] = [
    { key: 'todo', title: 'To Do', statuses: ['open', 'in_review'], dropStatus: 'open', accent: 'bg-blue-500' },
    { key: 'in_progress', title: 'In Progress', statuses: ['in_development'], dropStatus: 'in_development', accent: 'bg-purple-500' },
    { key: 'in_testing', title: 'In Testing', statuses: ['testing'], dropStatus: 'testing', accent: 'bg-cyan-500' },
    { key: 'blocked', title: 'Blocked', statuses: ['blocked'], dropStatus: 'blocked', accent: 'bg-red-500' },
    { key: 'done', title: 'Done', statuses: ['resolved', 'closed', 'wont_fix'], dropStatus: 'resolved', accent: 'bg-emerald-500' },
];

interface TicketsBoardProps {
    tickets: Ticket[];
    loading: boolean;
    canDrag: boolean;
    updatingId: number | null;
    onOpen: (ticket: Ticket) => void;
    onMove: (id: number, status: TicketStatus) => void;
    getAvatarUrl?: (ticket: Ticket) => string;
}

export default function TicketsBoard({
    tickets,
    loading,
    canDrag,
    updatingId,
    onOpen,
    onMove,
    getAvatarUrl,
}: TicketsBoardProps) {
    const [dragId, setDragId] = useState<number | null>(null);
    const [overKey, setOverKey] = useState<string | null>(null);

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
                <div className="absolute inset-0 z-10 flex items-center justify-center rounded-xl bg-white/80 dark:bg-gray-900/70 backdrop-blur-sm transition-opacity duration-200">
                    <div className="flex flex-col items-center gap-2">
                        <div className="w-8 h-8 border-4 border-purple-500 border-t-transparent rounded-full animate-spin"></div>
                        <p className="text-sm text-gray-600 dark:text-gray-300">Actualizando...</p>
                    </div>
                </div>
            )}

            <div className={`flex gap-3 overflow-x-auto pb-3 transition-opacity duration-200 ${loading ? 'opacity-50' : 'opacity-100'}`}>
                {BOARD_COLUMNS.map((column) => {
                    const items = tickets.filter((t) => column.statuses.includes(t.status));
                    const isOver = overKey === column.key;
                    return (
                        <div
                            key={column.key}
                            onDragOver={(e) => handleDragOver(e, column)}
                            onDragLeave={(e) => handleDragLeave(e, column)}
                            onDrop={(e) => handleDrop(e, column)}
                            className={`flex w-[280px] min-w-[280px] flex-col rounded-xl border transition-colors ${
                                isOver
                                    ? 'border-purple-400 bg-purple-50 dark:border-purple-500 dark:bg-purple-900/20'
                                    : 'border-gray-200 bg-gray-50 dark:border-gray-700 dark:bg-gray-900/40'
                            }`}
                        >
                            <div className="flex items-center gap-2 border-b border-gray-200 dark:border-gray-700 px-3 py-2.5">
                                <span className={`h-2.5 w-2.5 rounded-full ${column.accent}`}></span>
                                <h3 className="text-xs font-bold uppercase tracking-wider text-gray-700 dark:text-gray-200">
                                    {column.title}
                                </h3>
                                <span className="ml-auto inline-flex min-w-6 items-center justify-center rounded-full bg-gray-200 px-2 py-0.5 text-[11px] font-bold text-gray-700 dark:bg-gray-700 dark:text-gray-200">
                                    {items.length}
                                </span>
                            </div>

                            <div className="flex flex-1 flex-col gap-2 p-2 min-h-[140px]">
                                {items.length === 0 && (
                                    <div className="flex flex-1 items-center justify-center rounded-lg border border-dashed border-gray-300 dark:border-gray-700 p-4 text-center text-xs text-gray-400 dark:text-gray-500">
                                        Sin tickets
                                    </div>
                                )}

                                {items.map((t) => {
                                    const avatarUrl = getAvatarUrl ? getAvatarUrl(t) : '';
                                    const isUpdating = updatingId === t.id;
                                    return (
                                        <div
                                            key={t.id}
                                            draggable={canDrag}
                                            onDragStart={(e) => handleDragStart(e, t)}
                                            onDragEnd={handleDragEnd}
                                            onClick={() => onOpen(t)}
                                            className={`rounded-lg border border-gray-200 bg-white p-3 shadow-sm transition-all hover:shadow-md dark:border-gray-700 dark:bg-gray-800 ${
                                                canDrag ? 'cursor-grab active:cursor-grabbing' : 'cursor-pointer'
                                            } ${dragId === t.id ? 'opacity-40' : ''} ${
                                                isUpdating ? 'ring-2 ring-purple-400 animate-pulse' : ''
                                            }`}
                                        >
                                            <div className="flex items-center justify-between gap-2">
                                                <span className="font-mono text-[11px] text-gray-500 dark:text-gray-400">{t.code}</span>
                                                {column.statuses.length > 1 && <StatusBadge status={t.status} />}
                                            </div>

                                            <p className="mt-1.5 line-clamp-2 text-sm font-medium text-gray-900 dark:text-gray-100">
                                                {t.title}
                                            </p>

                                            <div className="mt-2 flex flex-wrap items-center gap-1.5">
                                                <TypeBadge type={t.type} />
                                                <PriorityBadge priority={t.priority} />
                                            </div>

                                            <div className="mt-2.5 flex items-center justify-between gap-2 border-t border-gray-100 dark:border-gray-700 pt-2">
                                                <div className="flex min-w-0 items-center gap-1.5">
                                                    {avatarUrl ? (
                                                        <img
                                                            src={avatarUrl}
                                                            alt=""
                                                            className="h-5 w-5 rounded-full object-cover ring-1 ring-gray-200 dark:ring-gray-600"
                                                        />
                                                    ) : (
                                                        <div className="flex h-5 w-5 items-center justify-center rounded-full bg-gray-200 text-[9px] text-gray-600 dark:bg-gray-600 dark:text-gray-300">
                                                            {t.assigned_to_name ? t.assigned_to_name[0].toUpperCase() : '-'}
                                                        </div>
                                                    )}
                                                    <span className="truncate text-[11px] text-gray-700 dark:text-gray-300">
                                                        {t.assigned_to_name || 'Sin asignar'}
                                                    </span>
                                                </div>
                                                <span className="max-w-[45%] truncate text-[11px] text-gray-500 dark:text-gray-400">
                                                    {t.business_name || (t.business_id ? `#${t.business_id}` : 'Interno')}
                                                </span>
                                            </div>
                                        </div>
                                    );
                                })}
                            </div>
                        </div>
                    );
                })}
            </div>
        </div>
    );
}

'use client';

import { STATUS_META, PRIORITY_META, TYPE_META, TicketStatus, TicketPriority, TicketType } from '../../domain/types';

export function StatusBadge({ status }: { status: TicketStatus }) {
    const m = STATUS_META[status];
    return (
        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold ring-1 ring-inset ${m.bg} ${m.color} ${m.ring}`}>
            {m.label}
        </span>
    );
}

export function PriorityBadge({ priority }: { priority: TicketPriority }) {
    const m = PRIORITY_META[priority];
    return (
        <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${m.bg} ${m.color}`}>
            {m.label}
        </span>
    );
}

export function TypeBadge({ type }: { type: TicketType }) {
    const m = TYPE_META[type];
    return (
        <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs font-medium bg-slate-100 dark:bg-slate-700 text-slate-700 dark:text-slate-200">
            <span className="font-mono text-[10px] opacity-70">{m.icon}</span>
            {m.label}
        </span>
    );
}
